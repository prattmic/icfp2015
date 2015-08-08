package main

import (
	"testing"
)

func TestMovePermuter(t *testing.T) {
	m := NewMovePermuter()
	possible := []Direction{E, W, SE, SW}

	var got []Direction
	for {
		d, err := m.Next()
		if err != nil {
			break
		}
		got = append(got, d)
	}

	if len(got) != 4 {
		t.Errorf("All m.Next() got %v, want len(got) == 4", got)
	}

	found := make([]bool, 4)
	for _, d := range got {
		for i, p := range possible {
			if d != p {
				continue
			}

			if found[i] {
				// Uh-oh, we found this before.
				t.Errorf("Found %v twice, want once", d)
			}
			found[i] = true
		}
	}

	for i, f := range found {
		if !f {
			t.Errorf("Didn't find %v", possible[i])
		}
	}
}
