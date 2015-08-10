package main

import (
	"log"
)

type RepeaterAI struct {
	index int
	str   string
	game  *Game
}

func NewRepeaterAI(g *Game, repeatStr string) AI {
	return &RepeaterAI{index: 0, str: repeatStr, game: g}
}

// Game returns the Game used by the AI.
// It may change after calls to Next().
func (ai *RepeaterAI) Game() *Game {
	return ai.game
}

// Next steps the AI one step, returning true if the game is
// complete, or an error if the game cannot continue.
func (ai *RepeaterAI) Next() (bool, error) {
	if ai.index < len(ai.str) {
		c := Command(ai.str[ai.index])
		ai.index++
		_, done, err := ai.game.Update(c)
		log.Printf("Update(%s) -> %v, %v", c, done, err)
		return done, err
	}
	return true, nil
}
