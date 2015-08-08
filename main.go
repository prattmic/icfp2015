package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	inputFile = flag.String("f", "", "File containing JSON encoded input.")
	timeLimit = flag.Int("t", 1000, "Time limit, in seconds, to produce output.")
	powerPhrase = flag.String("p", "", "Phrase of power")

	// TODO(myenik) Lol we should think about how to deal with this one...
	memLimit = flag.Int("m", 1000, "Memory limit, in megabytes, to produce output")
)

func ArgsOk() error {
	if *inputFile == "" {
		return fmt.Errorf("no file name specified")
	}

	return nil
}

func main() {
	flag.Parse()

	if err := ArgsOk(); err != nil {
		fmt.Printf("invalid arguments: %v\n", err)
		flag.Usage()
		return
	}

	f, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Could not open input file %s: %v\n", *inputFile, err)
		return
	}

	problem, err := ParseInputProblem(f)
	if err != nil {
		fmt.Printf("Could not parse JSON in input file %s: %v\n", *inputFile, err)
		return
	}

	fmt.Printf("Dump of parsed input:\n%+v\n", problem)
}
