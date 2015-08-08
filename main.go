package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// These are registered in init(), below.
	inputFiles multiStringValue
	powerPhrases multiStringValue

	timeLimit = flag.Int("t", 1000, "Time limit, in seconds, to produce output.")

	// TODO(myenik) Lol we should think about how to deal with this one...
	memLimit = flag.Int("m", 1000, "Memory limit, in megabytes, to produce output")
	cpus = flag.Int("c", 1, "Number of processor cores available")
)

// multiStringValue is a flag.Value which can be specified multiple times
// on the command line.
type multiStringValue []string

// String implements flag.Value.String.
func (s *multiStringValue) String() string {
	return fmt.Sprintf("%v", *s)
}

// String implements flag.Value.Set, adding each new item to the slice.
func (s *multiStringValue) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func ArgsOk() error {
	if len(inputFiles) == 0 {
		return fmt.Errorf("no file names specified")
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

	for _, name := range inputFiles {
		fmt.Printf("Processing %s\n", name)

		f, err := os.Open(name)
		if err != nil {
			fmt.Printf("Could not open input file %s: %v\n", name, err)
			return
		}

		problem, err := ParseInputProblem(f)
		if err != nil {
			fmt.Printf("Could not parse JSON in input file %s: %v\n", name, err)
			return
		}

		fmt.Printf("Dump of parsed input:\n%+v\n", problem)
	}
}

func init() {
	flag.Var(&inputFiles, "f",  "File containing JSON encoded input.")
	flag.Var(&powerPhrases, "p",  "Phrase of power")
}
