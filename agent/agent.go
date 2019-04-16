package agent

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sean2525/RL-tic-tac-toe/board"
)

var r *rand.Rand

// Init the rand seed
func Init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Agent definition
// An agent that will learn to play tic-tac-toe
type Agent struct {
	Values  map[string]float64 // States encountered during of all the games
	Sign    string             // The sign use to play
	alpha   float64
	gamma   float64
	epsilon float64
}

// New a agent
func New(alpha, gamma, epsilon float64, sign string) *Agent {
	return &Agent{
		alpha:   alpha,
		gamma:   gamma,
		Sign:    sign,
		epsilon: epsilon,
		Values:  make(map[string]float64),
	}
}

// Play Called when we need the agent to play
func (a *Agent) Play(b *board.Game) (string, error) {
	// Get the action from the policy
	bestAction, _ := a.policy(b)
	// Apply the action
	if err := b.MakeMove(bestAction); err != nil {
		fmt.Println(err)
		return "", err
	}
	return b.State, nil
}

// TrainPlay depend on epsilon random choice
func (a *Agent) TrainPlay(b *board.Game) error {
	_, nextAction := a.policy(b)
	if err := b.MakeMove(nextAction); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (a *Agent) getReward(winner string) float64 {
	if winner == a.Sign {
		return 1.0
	} else if winner == " " {
		return 0.0
	} else {
		return -1.0
	}
}

// The policy
// Given a state, it return an action
func (a *Agent) policy(b *board.Game) (string, string) {
	posStates, _ := b.AllowMoves()
	var nextAction string
	if len(posStates) == 0 {
		return "", ""
	}
	maxVal := a.Values[posStates[0]]
	// Get the hightest valued state from posStates
	for _, state := range posStates {
		stateValue := a.Values[state]
		if stateValue >= maxVal {
			maxVal = stateValue
		}
	}

	maxActions := []string{}
	for _, state := range posStates {
		if maxVal == a.Values[state] {
			maxActions = append(maxActions, state)
		}
	}

	var bestAction string
	if len(maxActions) > 1 {
		bestAction = maxActions[r.Intn(len(maxActions))]
	} else {
		bestAction = maxActions[0]
	}

	nextAction = bestAction
	// 10% of the time chose random actions
	if float64(r.Intn(100))/100.0 < a.epsilon {
		nextAction = posStates[r.Intn(len(posStates))]
	}
	return bestAction, nextAction
}

// LearnFromMove reward
// Q learning
func (a *Agent) LearnFromMove(state string, b *board.Game) {
	reward := a.getReward(b.Winner)
	currentStateValue := a.Values[state]
	nextStateValue := a.Values[b.State]
	// if b.PlayAble() {
	// 	best, _ := a.policy(b)
	// 	nextStateValue = a.Values[best]
	// }

	a.Values[state] = currentStateValue + a.alpha*(reward+a.gamma*nextStateValue-currentStateValue)
}

// Reset value
func (a *Agent) Reset() {
	a.Values = make(map[string]float64)
}

// InteractiveGame play
func (a *Agent) InteractiveGame() {
	t := 0
	round := 0
	reader := bufio.NewReader(os.Stdin)
	var b *board.Game
	for {
		round = 0
		if t%2 == 0 {
			b = board.New("X")
		} else {
			b = board.New("O")
		}
		t++
		for b.PlayAble() {
			round++
			if b.Player == a.Sign {
				a.Play(b)
			} else {
				for {
					allowed, index := b.AllowMoves()
					pIndex := []int{}
					for _, i := range index {
						pIndex = append(pIndex, i+1)
					}
					fmt.Printf("Please select one %v", pIndex)
					text, _ := reader.ReadString('\n')
					text = strings.ReplaceAll(text, "\r", "")
					text = strings.ReplaceAll(text, "\n", "")
					textInt, err := strconv.Atoi(text)
					if err != nil {
						fmt.Println(err)
						continue
					}
					for _i, i := range index {
						if i == textInt-1 {
							if _err := b.MakeMove(allowed[_i]); _err != nil {
								fmt.Println(_err)
								continue
							}
						}
					}
					break
				}
			}
			fmt.Printf("Round %d", round)
			fmt.Println(b)
		}
		if b.Winner != " " {
			fmt.Println(b.Winner + " Win this game")
		} else {
			fmt.Println("Draw")
		}
	}
}
