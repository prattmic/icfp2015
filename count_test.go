package main

import (
	"testing"
)

func TestCountOverlap(t *testing.T) {
	var cases = []struct {
		s    string
		sep  string
		want int
	}{
		{s: "Hello World", sep: "blah", want: 0},
		{s: "Hello World", sep: "Hello", want: 1},
		{s: "cthulhu cthulhu cthulhu", sep: "cthulhu", want: 3},
		{s: "cthulhu cthulhu cthulhu", sep: "cthulhu cthulhu", want: 2},
	}

	for _, c := range cases {
		if got := CountOverlap(c.s, c.sep); got != c.want {
			t.Errorf("CountOverlap(%q, %q) got %d want %d", c.s, c.sep, got, c.want)
		}
	}
}
