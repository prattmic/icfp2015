package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
)

var (
	defaultPhrases multiStringValue = []string{
		"ei!",
		"ia! ia!",
		"yuggoth",
		"r'lyeh",
		"planet 10",
		"tsathoggua",
		"monkeyboy",
	}

	// These are registered in init(), below.
	inputFiles   multiStringValue
	powerPhrases multiStringValue
	aiFlag       string

	timeLimit = flag.Int("t", 1000, "Time limit, in seconds, to produce output.")

	// TODO(myenik) Lol we should think about how to deal with this one...
	memLimit = flag.Int("m", 1000, "Memory limit, in megabytes, to produce output")
	cpus     = flag.Int("c", 1, "Number of processor cores available")

	serve = flag.Bool("serve", false, "Launch a web server")

	render   = flag.Bool("render", false, "Render the game")
	border   = flag.Int("border", 40, "Pixels of border to render")
	hexsize  = flag.Int("hexsize", 3, "Which size to render (0=tiny,1=small,2=med,3=large)")
	display  = flag.Bool("display", false, "Open the GIF after rendering")
	gifdelay = flag.Int("gif_delay", 10, "Time in 1/100ths of a second to wait between render frames.")

	graph = flag.String("graph", "", "Dump treeai graph to file with this prefix")

	profile = flag.String("profile", "", "Output CPU profile to file")

	repeat    = flag.String("repeat", "", "String for RepeaterAI to run")
	seed      = flag.Uint64("seed", 0xFFFFFFFFFFFFFFFF, "Use specific seed for single game")
	customtag = flag.String("customtag", "", "Custom tag for solution")
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
		log.Printf("invalid arguments: %v", err)
		flag.Usage()
		return
	}

	if len(powerPhrases) == 0 {
		powerPhrases = defaultPhrases
	}

	if *profile != "" {
		f, err := os.Create(*profile)
		if err != nil {
			log.Fatalf("Could not create profile file %s: %v", *profile, err)
		}

		pprof.StartCPUProfile(f)
	}

	if *serve {
		log.Printf("Running server...")
		runServer()
		return
	}

	var output []OutputEntry
	for _, name := range inputFiles {
		log.Printf("Processing %s", name)

		f, err := os.Open(name)
		if err != nil {
			log.Fatalf("Could not open input file %s: %v", name, err)
		}

		problem, err := ParseInputProblem(f)
		if err != nil {
			log.Fatalf("Could not parse JSON in input file %s: %v", name, err)
		}

		if *seed != uint64(0xFFFFFFFFFFFFFFFF) { // Hijack seed
			for i, s := range problem.SourceSeeds {
				if *seed == s {
					problem.SourceSeeds = problem.SourceSeeds[i : i+1]
				}
			}
		}

		// Take steps with random AI.
		// TODO(myenik) make rendering less gross/if'd out everywhere.
		for gi, g := range GamesFromProblem(problem) {
			var renderer *GameRenderer
			if *render {
				renderer = NewGameRenderer(g, *border, *hexsize)
			}

			log.Printf("Playing %+v", g)
			a := NewAI(g, aiFlag)

			i := 1
			for {
				log.Printf("Step %d", i)
				if *render {
					renderer.AddFrame(a.Game())
				}

				done, err := a.Next()
				if done {
					log.Println("Game done!")
					break
				} else if err != nil {
					log.Printf("a.Next error: %v", err)
					break
				}
				i++
			}

			log.Printf("Commands: %s", a.Game().Commands)
			log.Printf("Final Score: %f", a.Game().Score())

			if *render {
				gifname := fmt.Sprintf("%s_game%d.gif", name, gi)
				gif, err := os.Create(gifname)
				if err != nil {
					log.Fatalf("Failed to open output file %s: %v\n", gifname, err)
				}

				renderer.OutputGIF(gif, *gifdelay)
				if *display {
					c := exec.Command("sensible-browser", gifname)
					c.Start()
				}
			}

			if *customtag == "" {
				*customtag = fmt.Sprintf("Score: %v", a.Game().Score())
			}
			output = append(output, OutputEntry{
				ProblemId: problem.Id,
				Seed:      problem.SourceSeeds[gi],
				Tag:       *customtag,
				Solution:  a.Game().Commands.String(),
			})
		}
	}

	// Dump output
	if err := json.NewEncoder(os.Stdout).Encode(&output); err != nil {
		log.Fatalf("Failed to encode output %+v: %v", output, err)
	}

	if *profile != "" {
		pprof.StopCPUProfile()
	}
}

func init() {
	flag.Var(&inputFiles, "f", "File containing JSON encoded input.")
	flag.Var(&powerPhrases, "p", "Phrase of power")

	var keys string
	comma := ""
	for k := range ais {
		keys += fmt.Sprintf("%s %s", comma, k)
		if comma == "" {
			comma = ","
		}
	}

	flag.StringVar(&aiFlag, "ai", "treeai", fmt.Sprintf("AI to use:%s", keys))

}
