package main

import (
	"fmt"
	"strings"
)

type BoardCell struct {
	Cell
	Filled bool
}

// Boards facilitate search/simulation.
// Indexing into a grid with (x, y) gives
// you the cell on column x, in row y.
type Board struct {
	Width  int
	Height int
	Cells  [][]BoardCell
}

func NewBoard(w, h int, filled []Cell) *Board {
	b := &Board{Width: w, Height: h}

	// Make columns, according to [w][h]Cell.
	b.Cells = make([][]BoardCell, w)

	// Fill out each column with one spot for each row.
	for i := 0; i < w; i++ {
		b.Cells[i] = make([]BoardCell, h)
	}

	// Mark cell coordinates
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			b.Cells[x][y].X = x
			b.Cells[x][y].Y = y
		}
	}

	// Mark filled cells as filled.
	for _, c := range filled {
		b.MarkFilled(c)
	}

	return b
}

func (b *Board) Fork() *Board {
	w := b.Width
	h := b.Height
	bcopy := &Board{Width: w, Height: h}

	// Make columns, according to [w][h]Cell.
	bcopy.Cells = make([][]BoardCell, w)

	// Fill out each column with one spot for each row.
	for i := 0; i < w; i++ {
		bcopy.Cells[i] = make([]BoardCell, h)
	}

	// Copy cells.
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			bcopy.Cells[x][y] = b.Cells[x][y]
		}
	}

	return bcopy
}

// Pretty-print Board, indenting n levels
func (b *Board) StringLevel(n int) string {
	indent := strings.Repeat("\t", n)
	endindent := strings.Repeat("\t", n-1)

	return fmt.Sprintf(`Board{
%swidth:  %d,
%sheight: %d,
%scells:  %+v,
%s}`, indent, b.Width, indent, b.Height, indent, b.Cells, endindent)
}

func (b *Board) String() string {
	return b.StringLevel(1)
}

func (b *Board) BoardCell(c Cell) *BoardCell {
	return &b.Cells[c.X][c.Y]
}

func (b *Board) MarkFilled(c Cell) {
	b.BoardCell(c).Filled = true
}

func (b *Board) MarkUnfilled(c Cell) {
	b.BoardCell(c).Filled = false
}

func (b *Board) IsFilled(c Cell) bool {
	return b.BoardCell(c).Filled
}

func (b *Board) InBounds(c Cell) bool {
	isNotTooLow := c.X >= 0 && c.Y >= 0
	isNotTooHigh := c.X < b.Width && c.Y < b.Height
	return isNotTooLow && isNotTooHigh
}

func (b *Board) RowIsFilled(row int) bool {
	for i := 0; i < b.Width; i++ {
		if !b.Cells[i][row].Filled {
			return false
		}
	}

	return true
}

func (b *Board) UnfillRow(row int) bool {
	for i := 0; i < b.Width; i++ {
		b.Cells[i][row].Filled = false
	}

	return true
}

func (b *Board) TranslateRowDown(row int) {
	// Odd rows translate SW, even rows translate SE.
	dir := SW
	if row%2 == 0 {
		dir = SE
	}

	for i := 0; i < b.Width; i++ {
		if b.Cells[i][row].Filled {
			cell := Cell{X: i, Y: row}
			b.MarkFilled(cell.Translate(dir))
			b.MarkUnfilled(cell)
		}
	}
}

// ClearRows clears the lowest filled row and moves the tiles down.
// It returns true if a row was cleared.
func (b *Board) ClearRow() bool {
	for i := b.Height - 1; i >= 0; i-- {
		if b.RowIsFilled(i) {
			b.UnfillRow(i)
			for j := i - 1; j >= 0; j-- {
				b.TranslateRowDown(j)
			}
			return true
		}
	}

	return false
}

// ClearRows clears rows until there are no more to clear.
// It is a no-op if there are no rows to clear.
func (b *Board) ClearRows() (cleared int) {
	for b.ClearRow() {
		cleared++
	}
	return
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
		if b.GapBelow(c, 3) {
			return true
		}
	}

	return false
}
