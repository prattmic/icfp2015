package main

import (
	"fmt"
)

type Cell struct {
	X int
	Y int
}

type Direction int

const (
	E Direction = iota
	NE
	NW
	W
	SW
	SE
	CCW
	CW
)

func (d Direction) String() string {
	switch d {
	case E:
		return "E"
	case NE:
		return "NE"
	case W:
		return "W"
	case SW:
		return "SW"
	case SE:
		return "SE"
	case CCW:
		return "CCW"
	case CW:
		return "CW"
	default:
		return fmt.Sprintf("Unknown (%d)", d)
	}
}

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
	off, ok := offsets[d]
	if !ok {
		panic(fmt.Sprintf("Cannot translate in direction: %s\n", d))
	}

	rowIsEven := c.Y%2 == 0
	if rowIsEven {
		return Cell{c.X + off.ex, c.Y + off.ey}
	}

	return Cell{c.X + off.ox, c.Y + off.oy}
}

type CubeCell struct {
	x int
	y int
	z int
}

// All praise the great and merciful http://www.redblobgames.com/grids/hexagons
func (c Cell) ToCube() CubeCell {
	q := c.X - int((c.Y-(c.Y%2))/2)
	r := c.Y
	s := -q - r
	return CubeCell{q, r, s}
}

func (cc CubeCell) ToCell() Cell {
	col := cc.x + int((cc.y-(cc.y%2))/2)
	row := cc.y
	return Cell{col, row}
}

func (cc CubeCell) Rotate(counterClockwise bool) CubeCell {
	if counterClockwise {
		q := -cc.y
		r := -cc.z
		s := -cc.x
		return CubeCell{q, r, s}
	}

	q := -cc.z
	r := -cc.x
	s := -cc.y
	return CubeCell{q, r, s}
}

func (c Cell) OffsetFrom(other Cell) Cell {
	x := c.X - other.X
	y := c.Y - other.Y
	return Cell{x, y}
}

func (c Cell) Add(other Cell) Cell {
	x := c.X + other.X
	y := c.Y + other.Y
	return Cell{x, y}
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
