package main

import (
	"fmt"
	"os"
	"testing"
)

var (
	inputProblems []*InputProblem
)

func init() {
	// Initialize qualifier problems. If there are errors, this will fail to
	// work, but so wil the TestInputParsing test case...
	for _, fname := range QualifierProblemFilenames() {
		f, _ := os.Open(fname)
		p, _ := ParseInputProblem(f)
		inputProblems = append(inputProblems, p)
	}
}

func QualifierProblemFilenames() []string {
	var inputFiles []string

	for i := 0; i < 24; i++ {
		inputFiles = append(inputFiles, fmt.Sprintf("qualifiers/problem_%d.json", i))
	}

	return inputFiles
}

func TestInputParsing(t *testing.T) {
	for _, fname := range QualifierProblemFilenames() {
		f, err := os.Open(fname)
		if err != nil {
			t.Fatalf("Could not open %s: %v", fname, err)
		}

		_, err = ParseInputProblem(f)
		if err != nil {
			t.Errorf("Could not parse json in %s into an InputReader: %v", fname, err)
		}
	}
}

func QualifierProblems() []*InputProblem {
	return inputProblems
}

// TODO(myenik) More extensive testing, this just ensures the game builder doesn't blow up...
func TestGamesFromProblem(t *testing.T) {
	for _, p := range QualifierProblems() {
		GamesFromProblem(p)
	}
}

func TestLCG(t *testing.T) {
	// From the specs.
	testSeed := uint64(17)
	testNums := []uint64{0, 24107, 16552, 12125, 9427, 13152, 21440, 3383, 6873, 16117}

	testLCG := NewLCG(testSeed)
	for i, n := range testNums {
		ln := testLCG.Next()
		if n != ln {
			t.Errorf("wrong value for LCG output %d: got %d want %d", i, ln, n)
		}
	}
}
