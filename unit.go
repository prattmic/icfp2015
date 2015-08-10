package main

import (
	"math"
)

type Unit struct {
	Members []Cell
	Pivot   Cell
}

func (u *Unit) Size() int {
	return len(u.Members)
}

func (u *Unit) Translate(d Direction) *Unit {
	r := &Unit{
		Pivot:   u.Pivot.Translate(d),
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

func CopyUnits(units []*Unit) []*Unit {
	copies := make([]*Unit, len(units))

	for i, u := range units {
		copies[i] = u.DeepCopy()
	}

	return copies
}

// Deep copy copies the Unit and its cells.
func (u *Unit) DeepCopy() *Unit {
	r := &Unit{
		Pivot:   u.Pivot,
		Members: make([]Cell, len(u.Members)),
	}

	for i, c := range u.Members {
		r.Members[i] = c
	}

	return r
}

func (u *Unit) OverlapsAny(others []*Unit) bool {
	for _, other := range others {
		if u.Overlaps(other) {
			return true
		}
	}

	return false
}

// Left and rightmost Cells.
func (u *Unit) Bounds() (Cell, Cell) {
	leftmost := Cell{math.MaxInt32, 0}
	rightmost := Cell{0, 0}

	for _, c := range u.Members {
		if c.X < leftmost.X {
			leftmost = c
		}

		if c.X > rightmost.X {
			rightmost = c
		}
	}

	return leftmost, rightmost
}

func (u *Unit) Rotate(counterClockwise bool) *Unit {
	r := &Unit{
		Pivot:   u.Pivot,
		Members: make([]Cell, len(u.Members)),
	}

	p := r.Pivot.ToCube()
	for i, c := range u.Members {
		cc := c.ToCube().VectorFrom(p).Rotate(counterClockwise)
		r.Members[i] = p.Add(cc).ToCell()
	}

	return r
}

func (b *Board) GapBelow(c Cell, layers int) bool {
	if layers == 0 {
		return false
	}

	sw := c.Translate(SW)
	se := c.Translate(SE)
	swEmpty := b.InBounds(sw) && !b.IsFilled(sw)
	seEmpty := b.InBounds(se) && !b.IsFilled(se)
	return swEmpty || seEmpty || b.GapBelow(sw, layers-1) || b.GapBelow(se, layers-1)
}

func (b *Board) GapBelowAny(u *Unit) bool {
	for _, c := range u.Members {
		if b.GapBelow(c, 2) {
			return true
		}
	}

	return false
}
