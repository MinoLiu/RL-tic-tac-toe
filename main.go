package main

import (
	"fmt"
	"os"
	"time"

	agent "github.com/sean2525/RL-tic-tac-toe/agent"
	board "github.com/sean2525/RL-tic-tac-toe/board"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func nbWins(wins []entry, c string) int {
	count := 0
	for _, e := range wins {
		if e.Value == c {
			count++
		}
	}
	return count
}

type entry struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

func (e entry) String() string {
	return fmt.Sprintf("%v, %v\n", e.Timestamp, e.Value)
}

func demoGameStats(a1 *agent.Agent, numEpisodes int) {
	wins := make([]entry, 0, numEpisodes)
	draw := make([]entry, 0, numEpisodes)
	start := time.Now().UnixNano()
	var b *board.Game
	for i := 0; i < numEpisodes; i++ {
		b = board.New("X")

		for b.PlayAble() {
			a1.Play(b)
		}

		if b.Winner == a1.Sign {
			wins = append(wins, entry{time.Now().UnixNano() - start, b.Winner})
		} else if b.Winner == " " {
			draw = append(draw, entry{time.Now().UnixNano() - start, "Draw"})
		} else {
			wins = append(wins, entry{time.Now().UnixNano() - start, b.Winner})
		}
	}
	// Display new stats
	fmt.Printf("%v wins %v%% times\n", "X", float64(nbWins(wins, "X"))/float64(numEpisodes)*100)
	fmt.Printf("%v wins %v%% times\n", "O", float64(nbWins(wins, "O"))/float64(numEpisodes)*100)
	fmt.Printf("Draws %v%% times\n", float64(len(draw))/float64(numEpisodes)*100)
	if len(os.Args) > 1 && os.Args[1] == "--plot" {
		generateFigure(wins, draw, numEpisodes, a1)
	}

}

func main() {
	agent.Init()
	a1 := agent.New(0.5, 0.9, 0.1, "X")

	fmt.Println("before leaning")
	demoGameStats(a1, 3000)

	for i := 1; i <= 10; i++ {
		fmt.Printf("after learning %d times\n", i*3000)
		a1.LearnGame(3000)
		demoGameStats(a1, 3000)
	}

	a1.InteractiveGame()
}

// Generate figure from the wins array
func generateFigure(wins []entry, draw []entry, loopNb int, a1 *agent.Agent) {
	// Create plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	// Set plot meta data
	p.Title.Text = "Both learning"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Number of wins"
	// Build plot data
	ptsX := make(plotter.XYs, nbWins(wins, "X")+1)
	ptsO := make(plotter.XYs, nbWins(wins, "O")+1)
	ptsDraw := make(plotter.XYs, len(draw))
	countX := 0
	countO := 0
	for _, w := range wins[:loopNb] {
		if w.Value == "X" {
			countX++
			ptsX[countX].Y = float64(countX)
			ptsX[countX].X = float64(w.Timestamp)
		} else if w.Value == "O" {
			countO++
			ptsO[countO].Y = float64(countO)
			ptsO[countO].X = float64(w.Timestamp)
		}
	}
	for i, w := range draw {
		ptsDraw[i].Y = float64(i)
		ptsDraw[i].X = float64(w.Timestamp)
	}
	// Add data to plot
	err = plotutil.AddLines(p, "X", ptsX, "O", ptsO, "_", ptsDraw)
	if err != nil {
		panic(err)
	}
	// Save the plot to a PNG file.
	err = p.Save(4*vg.Inch, 4*vg.Inch, "points.png")
	if err != nil {
		panic(err)
	}
}
