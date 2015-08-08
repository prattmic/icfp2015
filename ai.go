package main

import (
	"errors"
	"log"
	"math/rand"
)

// MovePermuter provides each of the possible moves in a random order.
// TODO)prattmic): Better name...
type MovePermuter struct {
	possible []Direction
}

// NewMovePermuter builds a MovePermuter with all of the legal game moves.
func NewMovePermuter() MovePermuter {
	return MovePermuter{
		// All the possible moves.
		possible: []Direction{E, W, SE, SW, CCW, CW},
	}
}

// Next returns the next move to attempt, or an error if none remain.
func (m *MovePermuter) Next() (Direction, error) {
	if len(m.possible) == 0 {
		return 0, errors.New("No more possible moves")
	}

	i := rand.Intn(len(m.possible))
	move := m.possible[i]
	m.possible = append(m.possible[:i], m.possible[i+1:]...)

	return move, nil
}

// AI wins the game!
type AI struct {
	Game *Game

	// moves are the moves we actually took.
	moves []Direction
}

// NewAI builds a new AI.
func NewAI(g *Game) AI {
	return AI{Game: g}
}

// Next steps the game, returning true when the game is done.
func (a *AI) Next() (bool, error) {
	moves := NewMovePermuter()

	// Keep trying moves until one works.
	for {
		m, err := moves.Next()
		if err != nil {
			// No possible moves, we are stuck!
			return false, err
		}

		done, err := a.Game.Update(m)
		log.Printf("Update(%s) -> %v, %v", m, done, err)
		if err == nil {
			// Successful move/game done.
			a.moves = append(a.moves, m)
			return done, nil
		}
	}
}

// Moves dumps out the moves taken as a spec string.
func (a *AI) Moves() string {
	commands := map[Direction]byte{
		W:   '!', // []byte{'p', '\'', '!', '.', '0', '3'}
		E:   'e', // []byte{'b', 'c', 'e', 'f', 'y', '2'}
		SW:  'i', // []byte{'a', 'g', 'h', 'i', 'j', '4'}
		SE:  'l', // []byte{'l', 'm', 'n', 'o', ' ', '5'}
		CW:  'd', // []byte{'d', 'q', 'r', 'v', 'z', '1'}
		CCW: 'k', // []byte{'k', 's', 't', 'u', 'w', 'x'}
	}

	var b []byte
	for _, m := range a.moves {
		b = append(b, commands[m])
	}

	return string(b)
}
