package main

import (
	"fmt"
)

type InputProblem struct {
	Id           int
	Units        []Unit
	Width        int
	Height       int
	Filled       []Cell
	SourceLength int
	SourceSeeds  []uint64
}

const (
	lcgM   uint64 = 1103515245
	lcgI   uint64 = 12345
	lcgMod uint64 = 1 << 31
)

type GameLCG struct {
	current uint64
}

func NewLCG(seed uint64) GameLCG {
	return GameLCG{seed}
}

func (l *GameLCG) Next() uint64 {
	ret := (l.current >> 16) & 0x7fff
	l.current = (l.current*lcgM + lcgI) % lcgMod
	return ret
}

// A game holds the context of a running game, for simulation.
// the Units slice is initialized in reverse order from what
// was given in the input problem, to make popping them off easy and fast.
// So that means the last unit in the slice is the next one to be added
// to the game.
type Game struct {
	b         *Board
	units     []Unit
	lcg       GameLCG
	numUnits  int
	unitsSent int

	// Keep track of moves for current unit.
	currUnit      *Unit
	previousMoves []Unit
}

func GamesFromProblem(p *InputProblem) []*Game {
	nseeds := len(p.SourceSeeds)
	games := make([]*Game, nseeds)

	for i, s := range p.SourceSeeds {
		g := &Game{
			b:        NewBoard(p.Width, p.Height, p.Filled),
			lcg:      NewLCG(s),
			units:    p.Units,
			numUnits: p.SourceLength,
		}

		next, ok := g.NextUnit()
		if !ok {
			panic("no first move?")
		}

		g.currUnit = next
		if ok := g.placeUnit(next); !ok {
			panic("not ok?")
		}
		games[i] = g
	}

	return games
}

func (g *Game) String() string {
	return fmt.Sprintf(`Game{
	b:             %s,
	units:         %+v,
	lcg:           %+v,
	numUnits:      %d,
	unitsSent:     %d,
	currUnit:      %+v,
	previousMoves: %v,
}`, g.b.StringLevel(2), g.units, g.lcg, g.numUnits, g.unitsSent, g.currUnit, g.previousMoves)
}

func (g *Game) LockUnit(u *Unit) {
	for _, c := range u.Members {
		g.b.MarkFilled(c)
	}
}

func (g *Game) NextUnit() (*Unit, bool) {
	if g.unitsSent >= g.numUnits {
		return nil, false
	}

	rand := g.lcg.Next()
	idx := int(rand) % len(g.units)
	templUnit := &g.units[idx]

	// Do a deep copy of the chosen Unit.
	r := templUnit.DeepCopy()

	g.unitsSent++

	return &r, true
}

func (g *Game) placeUnit(u *Unit) bool {
	l, r := u.widthBounds()
	ucenter := (r - l) / 2
	bcenter := g.b.width / 2
	if ucenter == bcenter {
		return g.b.IsValid(u)
	}

	// We need to center it up.
	rightShift := bcenter - ucenter
	for i := range u.Members {
		u.Members[i].X += rightShift
		u.Pivot.X += rightShift
	}

	return g.b.IsValid(u)
}

// Update returns a bool indicating whether the game is done, and err to indicate and error (backwards move).
func (g *Game) Update(d Direction) (bool, error) {
	moved := g.currUnit.Translate(d)

	if moved.OverlapsAny(g.previousMoves) {
		return true, fmt.Errorf("moved unit from %+v to %+v and it overlaps with a previous move!", g.currUnit, moved)
	}

	if g.b.IsValid(&moved) {
		g.previousMoves = append(g.previousMoves, g.currUnit.DeepCopy())
		g.currUnit = &moved
		return false, nil
	}

	g.LockUnit(g.currUnit)
	g.b.ClearRows()
	nextUnit, ok := g.NextUnit()
	if !ok {
		// Game is done.
		return true, nil
	}

	if ok := g.placeUnit(nextUnit); !ok {
		// Game is done.
		return true, nil
	}

	g.previousMoves = append(g.previousMoves, g.currUnit.DeepCopy())
	g.currUnit = nextUnit
	return false, nil
}
