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

// LearnGame learn tic tac toe game
func (a *Agent) LearnGame(numEpisodes int) {
	for i := 0; i < numEpisodes; i++ {
		a.learnFromEpisode()
	}
}
func (a *Agent) learnFromEpisode() {
	b := board.New("O")
	_, move := a.policy(b)
	for move != "" {
		move = a.learnFromMove(b, move)
	}
}

func (a *Agent) learnFromMove(b *board.Game, move string) string {
	if err := b.MakeMove(move); err != nil {
		fmt.Println(err)
		return ""
	}
	r := a.getReward(b.Winner)
	nextStateValue := 0.0
	bestMove := ""
	nextMove := ""
	if b.PlayAble() {
		bestMove, nextMove = a.policy(b)
		nextStateValue = a.Values[bestMove]
	}
	currStateValue := a.Values[move]
	a.Values[move] = currStateValue + a.alpha*(r+a.gamma*nextStateValue-currStateValue)
	return nextMove
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

	geteway := a.Values[posStates[0]]

	if b.Player == a.Sign {
		for _, state := range posStates {
			stateValue := a.Values[state]
			if stateValue >= geteway {
				geteway = stateValue
			}
		}
	} else {
		for _, state := range posStates {
			stateValue := a.Values[state]
			if stateValue <= geteway {
				geteway = stateValue
			}
		}
	}

	actions := []string{}
	if b.Player == a.Sign {
		for _, state := range posStates {
			if geteway == a.Values[state] {
				actions = append(actions, state)
			}
		}
	} else {
		for _, state := range posStates {
			if geteway == a.Values[state] {
				actions = append(actions, state)
			}
		}
	}

	var bestAction string
	if len(actions) > 1 {
		bestAction = actions[r.Intn(len(actions))]
	} else {
		bestAction = actions[0]
	}

	nextAction = bestAction
	// 10% of the time chose random actions
	if float64(r.Intn(100))/100.0 < a.epsilon {
		nextAction = posStates[r.Intn(len(posStates))]
	}
	return bestAction, nextAction
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
			fmt.Printf("Round %d", round)
			fmt.Println(b)
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

		}
		if b.Winner != " " {
			fmt.Println(b.Winner + " Win this game")
		} else {
			fmt.Println("Draw")
		}
	}
}
