package main

import (
	"fmt"
	"os"
	"testing"
)

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

// TODO(myenik) Moar tests lol
