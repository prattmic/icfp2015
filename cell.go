package main

import (
	"fmt"
	"math"
)

type Cell struct {
	X int
	Y int
}

type Unit struct {
	Members []Cell
	Pivot   Cell
}

type InputProblem struct {
	Id           int
	Units        []Unit
	Width        int
	Height       int
	Filled       []Cell
	SourceLength int
	SourceSeeds  []uint64
}

type Direction int

const (
	E Direction = iota
	NE
	NW
	W
	SW
	SE
)

type offset struct {
	// for even row
	ex int
	ey int

	// for odd row
	ox int
	oy int
}

var (
	offsets = map[Direction]offset{
		E:  offset{1, 0, 1, 0},
		NE: offset{0, -1, 1, -1},
		NW: offset{-1, -1, 0, -1},
		W:  offset{-1, 0, -1, 0},
		SW: offset{-1, 1, 0, 1},
		SE: offset{0, 1, 1, 1},
	}
)

func (c Cell) Equals(other Cell) bool {
	return c.X == other.X && c.Y == other.Y
}

func (c Cell) Translate(d Direction) Cell {
	off := offsets[d]
	rowIsEven := c.Y%2 == 0
	if rowIsEven {
		return Cell{c.X + off.ex, c.Y + off.ey}
	}

	return Cell{c.X + off.ox, c.Y + off.oy}
}

// Rotate will return a new Unit that has undergone the specified
// number of rotations about its pivot cell. Rotations are mod 6
// wrt to the argument, for example a rotation of 2 yields:
// Before:
//   __    __    __
//  /NW\__/ 2\__/ E\
//  \__/ 3\__/ 1\__/
//  /  \__/ P\__/  \
//  \__/ 4\__/ 6\__/
//  / W\__/ 5\__/SE\
//  \__/  \__/  \__/
//
//  After:
//   __    __    __
//  /NW\__/ 6\__/ E\
//  \__/ 1\__/ 5\__/
//  /  \__/ P\__/  \
//  \__/ 2\__/ 4\__/
//  / W\__/ 3\__/SE\
//  \__/  \__/  \__/
//func (u *Unit) Rotate(i int) *Unit {
//	return nil
//}

// Boards facilitate search/simulation.
// Indexing into a grid with (x, y) gives
// you the cell on column x, in row y.
type Board struct {
	width  int
	height int
	cells  [][]BoardCell
}

type BoardCell struct {
	Cell
	filled bool
}

func (b *Board) BoardCell(c Cell) *BoardCell {
	return &b.cells[c.X][c.Y]
}

func (b *Board) MarkFilled(c Cell) {
	b.BoardCell(c).filled = true
}

func (b *Board) IsFilled(c Cell) bool {
	return b.BoardCell(c).filled
}

func (b *Board) InBounds(c Cell) bool {
	isNotTooLow := c.X >= 0 && c.Y >= 0
	isNotTooHigh := c.X < b.width && c.Y < b.height
	return isNotTooLow && isNotTooHigh
}

func (b *Board) RowIsFilled(row int) bool {
	for i := 0; i < b.width; i++ {
		if !b.cells[i][row].filled {
			return false
		}
	}

	return true
}

func (b *Board) UnfillRow(row int) bool {
	for i := 0; i < b.width; i++ {
		b.cells[i][row].filled = false
	}

	return true
}

func (b *Board) TranslateRowDown(row int) {
	// Odd rows translate SW, even rows translate SE.
	dir := SW
	if row%2 == 0 {
		dir = SE
	}

	for i := 0; i < b.width; i++ {
		if b.cells[i][row].filled {
			cell := Cell{X: i, Y: row}
			b.MarkFilled(cell.Translate(dir))
		}
	}
}

// ClearRows clears the lowest filled row and moves the tiles down.
// It returns true if a row was cleared.
func (b *Board) ClearRow() bool {
	for i := b.height - 1; i >= 0; i-- {
		if b.RowIsFilled(i) {
			for j := i; j >= 0; j-- {
				b.UnfillRow(j)
				b.TranslateRowDown(j)
			}
			return true
		}
	}

	return false
}

// ClearRows clears rows until there are no more to clear.
// It is a no-op if there are no rows to clear.
func (b *Board) ClearRows() {
	for b.ClearRow() {
	}
}

func NewBoard(w, h int, filled []Cell) *Board {
	b := &Board{width: w, height: h}

	// Make columns, according to [w][h]Cell.
	b.cells = make([][]BoardCell, w)

	// Fill out each column with one spot for each row.
	for i := 0; i < w; i++ {
		b.cells[i] = make([]BoardCell, h)
	}

	// Mark filled cells as filled.
	for _, c := range filled {
		b.MarkFilled(c)
	}

	return b
}

type GameLCG struct {
	current uint64
}

func NewLCG(seed uint64) GameLCG {
	return GameLCG{seed}
}

const (
	lcgM   uint64 = 1103515245
	lcgI   uint64 = 12345
	lcgMod uint64 = 1 << 31
)

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

func (u *Unit) Translate(d Direction) Unit {
	r := Unit{
		Pivot:   u.Pivot,
		Members: make([]Cell, len(u.Members)),
	}

	for i, c := range u.Members {
		r.Members[i] = c.Translate(d)
	}

	return r
}

func (u *Unit) Occupied() []Cell {
	return u.Members
}

// Returns whether or not any cell in the input slice equals the cell c.
func (c *Cell) EqualsAny(cells []Cell) bool {
	for _, other := range cells {
		if c.Equals(other) {
			return true
		}
	}

	return false
}

// TODO(myenik) This is O(stupid), there has to be a better way...
func (u *Unit) Overlaps(other *Unit) bool {
	// If there is any cell in the input unit that is not present in the
	// other unit, they do not overlap. Otherwise, they do overlap.
	for _, c := range u.Members {
		if !c.EqualsAny(other.Members) {
			return false
		}
	}

	return true
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
		games[i] = g
	}

	return games
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
	r := &Unit{
		Pivot:   templUnit.Pivot,
		Members: make([]Cell, len(templUnit.Members)),
	}

	for i, c := range templUnit.Members {
		r.Members[i] = c
	}

	return r, true
}

func (u *Unit) OverlapsAny(others []Unit) bool {
	for _, other := range others {
		if u.Overlaps(&other) {
			return true
		}
	}

	return false
}

func (u *Unit) widthBounds() (int, int) {
	leftmost := math.MaxInt32
	rightmost := 0
	for _, c := range u.Members {
		if c.X < leftmost {
			leftmost = c.X
		}

		if c.X > rightmost {
			rightmost = c.X
		}
	}

	return leftmost, rightmost
}

// A move is "invalid" if the moved unit has members that overlap filled cells
// or if any of the members leave the board.
func (b *Board) IsValid(u *Unit) bool {
	for _, c := range u.Members {
		if !b.InBounds(c) || b.IsFilled(c) {
			return false
		}
	}

	return true
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
		return false, nil
	}

	g.LockUnit(&moved)
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

	g.currUnit = nextUnit
	return false, nil
}
