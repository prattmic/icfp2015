package main

type SimpleAI struct {
	game *Game
}

func NewSimpleAI(g *Game) AI {
	return &SimpleAI{game: g}
}

func (a *SimpleAI) Game() *Game {
	return a.game
}

func (a *SimpleAI) Next() (bool, error) {
	// We want to go SW, pick a command from the
	// command slice.
	c := directionToCommands[SW][0]

	// If you want to test a move before committing it
	// to the 'official' game, Fork it first and update
	// the fork.
	// fork : = a.game.Fork()

	_, done, err := a.game.Update(c)

	return done, err
}
