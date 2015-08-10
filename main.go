package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"
)

var (
	defaultPhrases multiStringValue = []string{
		"case nightmare green",
		"john bigboote",
		"necronomicon",
		"the laundry",
		"tsathoggua",
		"blue hades",
		"planet 10",
		"monkeyboy",
		"yuggoth",
		"ia! ia!",
		"r'lyeh",
		"ei!",
	}

	// These are registered in init(), below.
	inputFiles   multiStringValue
	powerPhrases multiStringValue
	//aiFlags      multiStringValue = []string{"mcai", "cmc", "chanterai", "treeai"}
	aiFlags multiStringValue = []string{}

	timeLimit = flag.Int("t", 0, "Time limit, in seconds, to produce output.")

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

	debug = flag.Bool("debug", false, "enable logging")
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

type devNull struct{}

// Write consumes all.
func (d devNull) Write([]byte) (int, error) {
	return 0, nil
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

type AISolution struct {
	name     string
	commands string
	score    float64
}

func main() {
	flag.Parse()

	if err := ArgsOk(); err != nil {
		log.Printf("invalid arguments: %v", err)
		flag.Usage()
		return
	}

	if !*debug {
		log.SetOutput(devNull{})
	}

	if len(powerPhrases) == 0 {
		powerPhrases = defaultPhrases
	}
	normalizePhrases()

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

	var timeout <-chan time.Time
	timedOut := false
	if *timeLimit > 0 {
		var t time.Duration
		if *timeLimit > 10 {
			// Usually provided a 5 second buffer
			t = time.Duration(*timeLimit - 5)
		} else {
			// ... but with hardly any time at all, use
			// 90%.
			t = time.Duration(0.9 * float64(*timeLimit))
		}

		timeout = time.After(t * time.Second)
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
			if timedOut {
				break
			}

			var aiSolutions []AISolution
			for _, ai := range aiFlags {
				if timedOut {
					break
				}

				var renderer *GameRenderer
				if *render {
					renderer = NewGameRenderer(g, *border, *hexsize)
				}

				aiGame := g.Fork()

				log.Printf("Playing %+v", aiGame)
				log.Printf("Using AI: %s", ai)
				a := NewAI(aiGame, ai, *repeat)

				i := 1

			runGame:
				for {
					select {
					case <-timeout:
						timedOut = true
						break runGame
					default:
					}

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
				a.Game().WriteFinalCommands()
				log.Printf("Final Commands: %s", a.Game().FinalCommands)
				log.Printf("Final Score: %f", a.Game().FinalScore())

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

				aiSolutions = append(aiSolutions, AISolution{
					name:     ai,
					commands: a.Game().FinalCommands.String(),
					score:    a.Game().FinalScore(),
				})
			}

			best := AISolution{score: -1.0}
			for _, a := range aiSolutions {
				if a.score > best.score {
					best = a
				}
			}

			log.Printf("All solutions: %+v", aiSolutions)
			log.Printf("Best solution: %+v", best)

			mytag := *customtag
			if mytag == "" {
				mytag = fmt.Sprintf("Final Score: %v", best.score)
			}

			output = append(output, OutputEntry{
				ProblemId: problem.Id,
				Seed:      problem.SourceSeeds[gi],
				Tag:       mytag,
				Solution:  best.commands,
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

	flag.Var(&aiFlags, "ai", fmt.Sprintf("AI to use (can be multiple):%s", keys))

}
