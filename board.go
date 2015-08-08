package main

import (
	"fmt"
	"strings"
)

type BoardCell struct {
	Cell
	filled bool
}

// Boards facilitate search/simulation.
// Indexing into a grid with (x, y) gives
// you the cell on column x, in row y.
type Board struct {
	width  int
	height int
	cells  [][]BoardCell
}

func NewBoard(w, h int, filled []Cell) *Board {
	b := &Board{width: w, height: h}

	// Make columns, according to [w][h]Cell.
	b.cells = make([][]BoardCell, w)

	// Fill out each column with one spot for each row.
	for i := 0; i < w; i++ {
		b.cells[i] = make([]BoardCell, h)
	}

	// Mark cell coordinates
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			b.cells[x][y].X = x
			b.cells[x][y].Y = y
		}
	}

	// Mark filled cells as filled.
	for _, c := range filled {
		b.MarkFilled(c)
	}

	return b
}

func (b *Board) Fork() *Board {
	w := b.width
	h := b.height
	bcopy := &Board{width: w, height: h}

	// Make columns, according to [w][h]Cell.
	bcopy.cells = make([][]BoardCell, w)

	// Fill out each column with one spot for each row.
	for i := 0; i < w; i++ {
		bcopy.cells[i] = make([]BoardCell, h)
	}

	// Copy cells.
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			bcopy.cells[x][y] = b.cells[x][y]
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
%s}`, indent, b.width, indent, b.height, indent, b.cells, endindent)
}

func (b *Board) String() string {
	return b.StringLevel(1)
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
