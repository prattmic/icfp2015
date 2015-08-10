package main

import (
	"math/rand"
)

func tryDirection(g *Game, cc string, c string, scoresofar float64, tries int) (bool, float64, error) {
	//defer log.Printf("leaving! %+v\n", n)
	//log.Printf("tries: %d try dir: %+v node %+v\n", tries, d, n)
	thisUnit := g.currUnit.DeepCopy()
	locked, done, err := g.Update(Command(c[0]))
	if err != nil {
		return true, scoresofar + g.Score() - 1000000.0, err
	}

	if done {
		return true, scoresofar + g.Score() - 1000000.0, nil
	}

	if locked {
		if g.B.GapBelowAny(thisUnit) {
			return false, scoresofar + g.Score(), nil
		}
	}

	var ded bool
	var score float64
	nextcom := c[1:]
	if len(nextcom) > 0 {
		ded, score, err = tryDirection(g, cc, nextcom, scoresofar, tries-1)
		return false, score, nil
	}

	sz := len(defaultPhrases)
	rn := rand.Intn(sz)
	nextcom = defaultPhrases[rn]

	if tries == 0 {
		return false, scoresofar + g.Score(), err
	}

	for i := 0; i < pathEndRetries; i++ {
		ded, score, err = tryDirection(g.Fork(), nextcom, nextcom, scoresofar, tries-1)
		if !ded {
			break
		}
	}

	return false, score, err
}

type CMonteCarloid struct {
	g *Game
}

func NewCMonteCarloid(g *Game) AI {
	return &CMonteCarloid{g: g}
}

func (m *CMonteCarloid) Game() *Game {
	return m.g
}

var (
	chantDepth   = 6
	chantRetries = 100
)

func (m *CMonteCarloid) Next() (bool, error) {
	sz := len(defaultPhrases)
	command := defaultPhrases[rand.Intn(sz)]

	var ded bool
	var err error
	for i := 0; i < chantRetries; i++ {
		ded, _, err = tryDirection(m.g.Fork(), command, command, m.g.Score(), chantDepth)

		if !ded {
			break
		}

		//log.Printf("retry needed cause best node %+v ended game\n", best)
		command = defaultPhrases[rand.Intn(sz)]
	}

	for _, c := range command {
		m.g.Update(Command(c))
	}

	//log.Printf("next done: %+v", m.root)
	return ded, err
}
