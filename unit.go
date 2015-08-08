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

func (u *Unit) Translate(d Direction) Unit {
	r := Unit{
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

// Deep copy copies the Unit and its cells.
func (u *Unit) DeepCopy() Unit {
	r := Unit{
		Pivot:   u.Pivot,
		Members: make([]Cell, len(u.Members)),
	}

	for i, c := range u.Members {
		r.Members[i] = c
	}

	return r
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
