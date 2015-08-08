package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// These are registered in init(), below.
	inputFiles   multiStringValue
	powerPhrases multiStringValue

	timeLimit = flag.Int("t", 1000, "Time limit, in seconds, to produce output.")

	// TODO(myenik) Lol we should think about how to deal with this one...
	memLimit = flag.Int("m", 1000, "Memory limit, in megabytes, to produce output")
	cpus     = flag.Int("c", 1, "Number of processor cores available")

	render = flag.Bool("render", false, "Render the board and exit")
	serve  = flag.Bool("serve", false, "Launch a web server")
	gifdelay  = flag.Int("gif_delay", 100, "Time in 1/100ths of a second to wait between render frames.")
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
	if *serve {
		return nil
	}

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

	if *serve {
		fmt.Printf("Running server...\n")
		runServer()
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

		// Take steps with random AI.
		// TODO(myenik) make rendering less gross/if'd out everywhere.
		for gi, g := range GamesFromProblem(problem)[:1] {
			var renderer *GameRenderer
			if *render {
				renderer = NewGameRenderer()
			}

			fmt.Printf("Playing %+v\n", g)
			a := NewAI(g)

			i := 1
			for {
				fmt.Printf("Step %d\n", i)
				if *render {
					renderer.AddFrame(g)
				}

				done, err := a.Next()
				if done {
					fmt.Printf("Game done!\n")
					break
				} else if err != nil {
					fmt.Printf("a.Next error: %v\n", err)
					break
				}
				i++
			}

			if *render {
				gifname := fmt.Sprintf("%s_game%d.gif", name, gi)
				gif, err := os.Create(gifname)
				if err != nil {
					fmt.Printf("Failed to open output file %s: %v\n", gifname, err)
				}

				renderer.OutputGIF(gif, *gifdelay)
				return
			}
		}

	}
}

func init() {
	flag.Var(&inputFiles, "f", "File containing JSON encoded input.")
	flag.Var(&powerPhrases, "p", "Phrase of power")
}
