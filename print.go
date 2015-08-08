package main

import (
	"fmt"
	"io"
)

func IsFilled(filled []Cell, x, y int) bool {
	for _, cell := range filled {
		if (cell.X == x) && cell.Y == y {
			return true
		}
	}
	return false
}

func PrintBoard(w io.Writer, width, height int, filled []Cell, unit *Unit) {
	fmt.Fprintf(w, "Width %d\n", width)
	fmt.Fprintf(w, "Height %d\n", height)
	fmt.Fprintf(w, "Cells %d\n", len(filled))
	if unit != nil {
		fmt.Fprintf(w, "Unit\n")
	}

	fmt.Fprintf(w, "\n\n")

	for i := 0; i < height; i++ {
		s := make([]string, 4)
		if i%2 == 1 { // indent odd rows
			s[0] = ` \  `
			s[1] = `   \`
			s[2] = `    `
			s[3] = `    `
		}
		for j := 0; j < width; j++ {
			s[0] += fmt.Sprintf(`   / \  `)
			s[1] += fmt.Sprintf(` / %x,%x \`, j%16, i%16)
			if IsFilled(filled, j, i) {
				s[2] += `| xxxxx `
				s[3] += `| xxxxx `
			} else {
				s[2] += `|       `
				s[3] += `|       `
			}
		}
		if (i%2 == 0) && (i > 0) { // trailing edges of rows
			s[0] += `   /`
			s[1] += ` /`
		}
		fmt.Fprintf(w, "%s\n", s[0])
		fmt.Fprintf(w, "%s\n", s[1])
		fmt.Fprintf(w, "%s|\n", s[2])
		fmt.Fprintf(w, "%s|\n", s[3])
	}
	a := ``
	b := ``
	if height%2 == 0 {
		a += `    `
		b += `    `
	}
	for j := 0; j < width; j++ {
		a += ` \     /`
		b += `   \ /  `
	}
	fmt.Fprintf(w, "%s\n", a)
	fmt.Fprintf(w, "%s\n", b)
}
