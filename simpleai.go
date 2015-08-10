package main

import (
	"log"
)

type SimpleAI struct {
	index int
	game  *Game
}

func NewSimpleAI(g *Game, _ string) AI {
	return &SimpleAI{index: 0, game: g}
}

// Game returns the Game used by the AI.
// It may change after calls to Next().
func (ai *SimpleAI) Game() *Game {
	return ai.game
}

// func CheckPath(g *Game, int xOrigin, int yOrigin, int xTarget, int yTarget) bool {
// 	board := g.B.Cells

// }

func moveDirection(ai *SimpleAI, d Direction) (bool, bool, error) {
	fork := ai.game.Fork()
	locked, done, err := fork.Update(directionToCommands[d][0])

	if !locked {
		ai.game = fork
	}

	return locked, done, err
}

// Next steps the AI one step, returning true if the game is
// complete, or an error if the game cannot continue.
func (ai *SimpleAI) Next() (bool, error) {
	var coms Commands
	for _, d := range []Direction{E, W, SE, SW, CCW, CW} {
		coms = append(coms, directionToCommands[d][0])
	}

	leftMost := 0
	bottomMost := 0

	for y := range ai.game.B.Cells {
		for x := range ai.game.B.Cells[y] {
			cell := ai.game.B.Cells[x][y]
			if cell.Filled == false {
				if bottomMost < y {
					bottomMost = y
					leftMost = x
					continue
				}
			}
		}
	}

	log.Printf("currUnit %v,", ai.game.currUnit)

	firstMember := ai.game.currUnit.Members[0]

	// move left
	if firstMember.X > leftMost {
		locked, done, err := moveDirection(ai, W)

		if !locked {
			return done, err
		}
	}

	// move right
	if firstMember.X < leftMost {
		locked, done, err := moveDirection(ai, E)

		if !locked {
			return done, err
		}
	}

	// move southwest
	locked, done, err := moveDirection(ai, SW)

	if !locked {
		return done, err
	}

	// // move southeast
	locked, done, err = moveDirection(ai, SE)

	if !locked {
		return done, err
	}

	_, done, err = ai.game.Update(directionToCommands[SE][0])

	return done, err
}
