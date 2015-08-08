package main

import (
	"errors"
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
		possible: []Direction{E, W, SE, SW},
	}
}

// Next returns the next move to attempt, or an error if none remain.
func (m *MovePermuter) Next() (Direction, error) {
	if len(m.possible) == 0 {
		return 0, errors.New("No more possible moves")
	}

	i := rand.Intn(len(m.possible))
	d := m.possible[i]

	m.possible = append(m.possible[:i], m.possible[i+1:]...)

	return d, nil
}

// AI wins the game!
type AI struct {
	Game *Game
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

		if done, err := a.Game.Update(m); err == nil {
			// Successful move/game done.
			return done, nil
		}
	}
}
