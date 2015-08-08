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
	off := offsets[d]
	rowIsEven := c.Y%2 == 0
	if rowIsEven {
		return Cell{c.X + off.ex, c.Y + off.ey}
	}

	return Cell{c.X + off.ox, c.Y + off.oy}
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
