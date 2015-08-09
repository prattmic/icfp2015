package main

import (
	"log"
)

type AI interface {
	// Next steps the AI one step, returning done if the game is
	// complete, or an error if the game cannot continue.
	Next() (done bool, err error)

	// Game returns the Game used by the AI.
	// It may change after calls to Next().
	Game() *Game
}

var ais = map[string]func(*Game) AI{
	"treeai":      NewTreeAI,
	"lookaheadai": NewLookaheadAI,
	"repeaterai":  NewRepeaterAI,
}

func NewAI(g *Game) AI {
	fn, ok := ais[aiFlag]
	if !ok {
		log.Fatalf("Invalid AI %q", aiFlag)
	}

	return fn(g)
}
