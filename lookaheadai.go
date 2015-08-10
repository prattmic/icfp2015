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
func (m *MovePermuter) Next() (Command, error) {
	if len(m.possible) == 0 {
		return 0, errors.New("No more possible moves")
	}

	i := rand.Intn(len(m.possible))
	move := m.possible[i]
	m.possible = append(m.possible[:i], m.possible[i+1:]...)

	// Pick a random command for this direction
	c := directionToCommands[move][rand.Intn(len(directionToCommands[move]))]

	return c, nil
}

// LookaheadAI is really bad (and implements AI)
type LookaheadAI struct {
	game *Game
}

// NewLookaheadAI builds a new LookaheadAI.
func NewLookaheadAI(g *Game) AI {
	return &LookaheadAI{game: g}
}

func (a *LookaheadAI) Game() *Game {
	return a.game
}

type aiResult struct {
	command Command
	game    *Game
	score   float64
	done    bool
	err     error
}

func (a *LookaheadAI) Next() (bool, error) {
	var coms Commands
	for _, d := range []Direction{E, W, SE, SW, CCW, CW} {
		coms = append(coms, directionToCommands[d][0])
	}

	// Try each command, use the highest score.
	ch := make(chan aiResult, len(coms))

	runner := func(c Command) {
		g := a.game.Fork()
		_, done, err := g.Update(c)

		ch <- aiResult{
			command: c,
			game:    g,
			score:   g.Score(),
			done:    done,
			err:     err,
		}
	}

	// Compute moves
	for _, c := range coms {
		go runner(c)
	}

	// Collect results
	var ret []aiResult
	for _ = range coms {
		ret = append(ret, <-ch)
	}

	bestScore := -1.0
	var best aiResult
	for _, r := range ret {
		if r.err != nil {
			log.Printf("skipping: %+v", r)
			continue
		}

		if r.score > bestScore {
			if r.done && !best.done {
				// Don't override non-done with done.
				continue
			}
			log.Printf("best by score: %+v", r)
			bestScore = r.score
			best = r
		}

		if best.done && !r.done {
			log.Printf("best by not done: %+v", r)
			bestScore = r.score
			best = r
		}
	}

	log.Printf("results: %+v", ret)

	if best.game == nil {
		// No best move! Just use first.
		best = ret[0]
		log.Printf("best by best is nil: %+v", best)
	}

	log.Printf("best: %+v", best)

	a.game = best.game
	return best.done, best.err
}
