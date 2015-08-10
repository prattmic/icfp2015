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
	NOP
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
	case NOP:
		return "NOP"
	default:
		return fmt.Sprintf("Unknown (%d)", d)
	}
}

type Command byte
type Commands []Command

var commandToDirection = map[Command]Direction{
	'p':  W,
	'\'': W,
	'!':  W,
	'.':  W,
	'0':  W,
	'3':  W,

	'b': E,
	'c': E,
	'e': E,
	'f': E,
	'y': E,
	'2': E,

	'a': SW,
	'g': SW,
	'h': SW,
	'i': SW,
	'j': SW,
	'4': SW,

	'l': SE,
	'm': SE,
	'n': SE,
	'o': SE,
	' ': SE,
	'5': SE,

	'd': CW,
	'q': CW,
	'r': CW,
	'v': CW,
	'z': CW,
	'1': CW,

	'k': CCW,
	's': CCW,
	't': CCW,
	'u': CCW,
	'w': CCW,
	'x': CCW,

	'\t': NOP,
	'\n': NOP,
	'\r': NOP,
}

var directionToCommands = map[Direction]Commands{
	W:   Commands{'!', '\'', 'p', '.', '0', '3'},
	E:   Commands{'e', 'c', 'b', 'f', 'y', '2'},
	SW:  Commands{'i', 'g', 'h', 'a', 'j', '4'},
	SE:  Commands{'l', 'm', 'n', 'o', ' ', '5'},
	CW:  Commands{'d', 'q', 'r', 'v', 'z', '1'},
	CCW: Commands{'k', 's', 't', 'u', 'w', 'x'},
}

func (c Command) String() string {
	d, ok := commandToDirection[c]
	if !ok {
		return fmt.Sprintf("%c (Unknown)", byte(c))
	}

	return fmt.Sprintf("%c (%s)", byte(c), d)
}

func (c *Commands) String() string {
	b := make([]byte, len(*c))
	for i, v := range *c {
		b[i] = byte(v)
	}

	return string(b)
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
	X int
	Y int
	Z int
}

// All praise the great and merciful http://www.redblobgames.com/grids/hexagons
func (c Cell) ToCube() CubeCell {
	q := c.X - int((c.Y-(c.Y%2))/2)
	r := c.Y
	s := -q - r
	return CubeCell{q, r, s}
}

func (cc CubeCell) ToCell() Cell {
	col := cc.X + int((cc.Y-(cc.Y%2))/2)
	row := cc.Y
	return Cell{col, row}
}

// We have types, let's use them.
type CubeVector CubeCell

func (c CubeVector) Rotate(counterClockwise bool) CubeVector {
	if counterClockwise {
		q := -c.Z
		r := -c.X
		s := -c.Y
		return CubeVector{q, r, s}
	}

	q := -c.Y
	r := -c.Z
	s := -c.X
	return CubeVector{q, r, s}
}

func (c CubeCell) VectorFrom(other CubeCell) CubeVector {
	x := c.X - other.X
	y := c.Y - other.Y
	z := c.Z - other.Z
	return CubeVector{x, y, z}
}

func (c CubeCell) Add(v CubeVector) CubeCell {
	x := c.X + v.X
	y := c.Y + v.Y
	z := c.Y + v.Z
	return CubeCell{x, y, z}
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
