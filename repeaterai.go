package main

import (
	"log"
)

var (
	repeatIndex = 0
)

type RepeaterAI struct {
	game *Game
}

func NewRepeaterAI(g *Game) AI {
	return &RepeaterAI{game: g}
}

// Game returns the Game used by the AI.
// It may change after calls to Next().
func (ai *RepeaterAI) Game() *Game {
	return ai.game
}

// Next steps the AI one step, returning true if the game is
// complete, or an error if the game cannot continue.
func (ai *RepeaterAI) Next() (bool, error) {
	if repeatIndex < len(*repeat) {
		c := Command((*repeat)[repeatIndex])
		repeatIndex++
		done, err := ai.game.Update(c)
		log.Printf("Update(%s) -> %v, %v", c, done, err)
		return done, err
	}
	return true, nil
}
