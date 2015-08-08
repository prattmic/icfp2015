package main

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
	SourceSeeds  []int
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

func (c Cell) Translate(d Direction) Cell {
	switch d {
	case E:
		return Cell{c.X + 1, c.Y}
	case NE:
		return Cell{c.X + 1, c.Y - 1}
	case NW:
		return Cell{c.X - 1, c.Y - 1}
	case W:
		return Cell{c.X - 1, c.Y}
	case SW:
		return Cell{c.X - 1, c.Y + 1}
	case SE:
		return Cell{c.X + 1, c.Y + 1}
	}

	// wat
	return Cell{}
}
