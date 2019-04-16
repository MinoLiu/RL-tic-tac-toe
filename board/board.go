package board

import (
	"errors"
	"fmt"
	"strings"
)

// Game struct
type Game struct {
	State  string
	Winner string
	Player string
}

// New a boardGame
func New(sign string) *Game {
	return &Game{
		State:  "         ",
		Winner: " ",
		Player: sign,
	}
}

// IsNotFull check board not Full
func (b *Game) IsNotFull() bool {
	return strings.Count(b.State, " ") != 0
}

// PlayAble avaliable to play
func (b *Game) PlayAble() bool {
	return b.IsNotFull() && b.Winner == " "
}

// PredictWinner get the Winner string
// Return the Winner
// If no Winner, return " "
func (b *Game) PredictWinner() string {
	// Check all win result
	lines := [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}
	Winner := " "
	for _, line := range lines {
		lineState := string(b.State[line[0]]) + string(b.State[line[1]]) + string(b.State[line[2]])
		if lineState == "XXX" {
			Winner = "X"
		}
		if lineState == "XXX" {
			Winner = "X"
		} else if lineState == "OOO" {
			Winner = "O"
		}
	}
	return Winner
}

// Overide default String method
func (b *Game) String() string {
	s := strings.Split(b.State, "")
	return fmt.Sprintf(`
	 %s | %s | %s
	-----------
	 %s | %s | %s
	-----------
	 %s | %s | %s
	`, s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8])
}

// AllowMoves allowed moves
func (b *Game) AllowMoves() ([]string, []int) {
	allowed := []string{}
	index := []int{}
	for i, s := range b.State {
		if string(s) == " " {
			allowed = append(allowed, b.State[:i]+b.Player+b.State[i+1:])
			index = append(index, i)
		}
	}
	return allowed, index
}

// MakeMove make move
func (b *Game) MakeMove(nextState string) error {
	if b.Winner != " " {
		return errors.New("game already completed, cannot make another move")
	}
	if b.validMove(nextState) != true {
		return fmt.Errorf("cannot make move %s to %s for Player %s", b.State, nextState, b.Player)
	}

	b.State = nextState
	b.Winner = b.PredictWinner()
	if b.Winner != " " {
		b.Player = " "
	} else if b.Player == "X" {
		b.Player = "O"
	} else if b.Player == "O" {
		b.Player = "X"
	}

	return nil
}

func (b *Game) validMove(nextState string) bool {
	allowed, _ := b.AllowMoves()
	for _, s := range allowed {
		if nextState == s {
			return true
		}
	}
	return false
}
