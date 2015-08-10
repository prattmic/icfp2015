package main

import (
	"log"
	"math/rand"
)

// This remembers where was probed
type MCNode struct {
	probed []*MCNode
	g      *Game
	done   bool
	err    error
	score  float64
}

type weightedDir struct {
	w int
	d Direction
}

var (
	probeDepth     = 300
	probeWidth     = 6
	gameEndRetries = 1000000
	pathEndRetries = 10
	weightedDirs   = []weightedDir{
		weightedDir{25, SE},
		weightedDir{25, SW},
		weightedDir{75, E},
		weightedDir{75, W},
		weightedDir{25, CW},
		weightedDir{25, CCW},
	}
)

func drawDir(wds []weightedDir) Direction {
	cum := 0
	lim := 0
	for _, wd := range wds {
		lim += wd.w
	}

	//log.Printf("lim: %d\n", lim)
	rand := rand.Intn(lim)
	for _, wd := range wds {
		if cum+wd.w > rand {
			return wd.d
		}
		cum += wd.w
	}

	panic("REALITY IS BROKEN")
}

func dirInSlice(d Direction, ds []Direction) bool {
	for _, di := range ds {
		if d == di {
			return true
		}
	}
	return false
}

func weightedDirsCopy() []weightedDir {
	wds := append([]weightedDir(nil), weightedDirs...)
	return wds
}

func removeDir(d Direction, wds *[]weightedDir) {
	for i, wd := range *wds {
		if wd.d == d {
			(*wds)[i] = (*wds)[len(*wds)-1]
			*wds = (*wds)[:len(*wds)-1]
			return
		}
	}
}

func removeDirs(ds []Direction, wds *[]weightedDir) {
	for _, d := range ds {
		removeDir(d, wds)
	}
}

func drawDirs(n int, wds []weightedDir) []Direction {
	var ds []Direction
	if n > len(wds) {
		n = len(wds)
	}

	for done := 0; done < n; {
		d := drawDir(wds)
		if !dirInSlice(d, ds) {
			ds = append(ds, d)
			done++
		}
	}

	return ds
}

func (n *MCNode) tryDirection(d Direction, scoresofar float64, tries int) (bool, float64) {
	//defer log.Printf("leaving! %+v\n", n)
	//log.Printf("tries: %d try dir: %+v node %+v\n", tries, d, n)
	thisUnit := n.g.currUnit
	locked, done, err := n.g.Update(directionToCommands[d][0])
	n.done, n.err = done, err
	if err != nil {
		return true, scoresofar + n.g.Score() - 1000000.0
	}

	if done {
		return true, scoresofar + n.g.Score() - 1000000.0
	}

	// We must go deeper
	if d == SW || d == SE {
		scoresofar += 10.0
	}

	if locked {
		if n.g.B.GapBelowAny(thisUnit) {
			return false, scoresofar + n.g.Score()
		}
	}

	if tries == 0 {
		return false, scoresofar + n.g.Score()
	}

	tried := n.probed[int(d)]
	var ded bool
	var score float64
	if tried == nil {
		currdirs := weightedDirsCopy()
		for i := 0; i < pathEndRetries; i++ {
			tried = &MCNode{
				g:      n.g.Fork(),
				probed: make([]*MCNode, int(NOP)+1),
			}

			dir := drawDir(currdirs)
			ded, score = tried.tryDirection(dir, scoresofar, tries-1)
			if !ded {
				break
			}

			//log.Printf("dir %v ded at try %d\n", dir, tries)
			removeDir(dir, &currdirs)
			if len(currdirs) == 0 {
				//log.Printf("no dirs left for try %d\n", tries)
				break
			}
		}

		n.probed[int(d)] = tried
		n.score = score
	} else {
		score = tried.score
	}

	return false, score
}

func (root *MCNode) tryDirections(n int, wds *[]weightedDir) (Direction, []Direction) {
	ds := drawDirs(n, *wds)
	//log.Printf("drawn dirs: %+v\n", ds)
	for _, d := range ds {
		chld := root.probed[int(d)]
		if chld == nil {
			chld = &MCNode{
				g:      root.g.Fork(),
				probed: make([]*MCNode, int(NOP)+1),
			}

			_, _ = chld.tryDirection(d, 0.0, probeDepth)
			root.probed[int(d)] = chld
		}
	}

	// Find best direction to go in.
	bestDir := ds[0]
	//log.Printf("best dir: %s %d\n", bestDir, int(bestDir))
	//log.Printf("root: %+v\n", root)
	bestScore := root.probed[int(bestDir)].score
	//log.Printf("best start score: %v\n", bestScore)
	for _, d := range ds {
		chld := root.probed[int(d)]
		if chld.score > bestScore {
			bestDir = d
			bestScore = chld.score
			//log.Printf("new best score: %v\n", bestScore)
		}
	}

	return bestDir, ds
}

type MonteCarloid struct {
	g    *Game
	root *MCNode
}

func NewMonteCarloid(g *Game, st string) AI {
	newroot := &MCNode{
		g:      g,
		probed: make([]*MCNode, int(NOP)+1),
	}
	return &MonteCarloid{g: g, root: newroot}
}

func (m *MonteCarloid) Game() *Game {
	return m.g
}

func (m *MonteCarloid) Next() (bool, error) {
	var best *MCNode
	var d Direction
	//var ds []Direction

	wds := weightedDirsCopy()
	for i := 0; i < gameEndRetries; i++ {
		d, _ = m.root.tryDirections(probeWidth, &wds)
		//d, ds = m.root.tryDirections(probeWidth, &wds)
		//log.Printf("best chosen direction: %s", d)

		// Keep this one if it didn't end the game
		best = m.root.probed[int(d)]
		//log.Printf("tried %+v best dir %v best %+v\n", ds, d, best)
		if !best.done && best.err == nil {
			//log.Printf("using best chosen direction: %s", d)
			break
		}

		//log.Printf("retry needed cause best node %+v ended game\n", best)
		removeDir(d, &wds)
		if len(wds) == 0 {
			log.Printf("NO DIRS LEFT TO TRY!\n")
			break
		}

	}

	if m.root == nil {
		panic("ROOT IS NIL BUT YOU JUST WENT IN THAT DIRECTION")
	}

	m.root = best
	m.g = m.root.g
	//log.Printf("next done: %+v", m.root)
	return m.root.done, m.root.err
}
