package main

import (
	"fmt"
	"strings"
)

const (
	lcgM   uint64 = 1103515245
	lcgI   uint64 = 12345
	lcgMod uint64 = 1 << 31
)

type GameLCG struct {
	current uint64
}

func NewLCG(seed uint64) GameLCG {
	return GameLCG{seed}
}

func (l *GameLCG) Next() uint64 {
	ret := (l.current >> 16) & 0x7fff
	l.current = (l.current*lcgM + lcgI) % lcgMod
	return ret
}

// A game holds the context of a running game, for simulation.
// the Units slice is initialized in reverse order from what
// was given in the input problem, to make popping them off easy and fast.
// So that means the last unit in the slice is the next one to be added
// to the game.
type Game struct {
	// Accumulated move score so far.
	moveScore float64
	powerWordCount map[string]int

	// All previous commands sent to the game.
	Commands Commands
	// Recomputed commands for submission, using actual phrases
	FinalCommands Commands

	B         *Board
	units     []Unit
	lcg       GameLCG
	numUnits  int
	unitsSent int

	// Keep track of moves for current unit.
	currUnit             *Unit
	previousMoves        []*Unit
	previousLinesCleared int
}

func (g *Game) Fork() *Game {
	n := &Game{
		moveScore:            g.moveScore,
		B:                    g.B.Fork(),
		units:                g.units,
		lcg:                  g.lcg,
		numUnits:             g.numUnits,
		unitsSent:            g.unitsSent,
		currUnit:             g.currUnit.DeepCopy(),
		previousMoves:        CopyUnits(g.previousMoves),
		previousLinesCleared: g.previousLinesCleared,
	}

	n.Commands = make(Commands, len(g.Commands))
	copy(n.Commands, g.Commands)

	n.powerWordCount = make(map[string]int)
	for k, v := range g.powerWordCount {
		n.powerWordCount[k] = v
	}

	return n
}

func GamesFromProblem(p *InputProblem) []*Game {
	nseeds := len(p.SourceSeeds)
	games := make([]*Game, nseeds)

	for i, s := range p.SourceSeeds {
		g := &Game{
			B:        NewBoard(p.Width, p.Height, p.Filled),
			lcg:      NewLCG(s),
			units:    p.Units,
			numUnits: p.SourceLength,
			powerWordCount: make(map[string]int),
		}

		next, ok := g.NextUnit()
		if !ok {
			panic("no first move?")
		}

		g.currUnit = next
		if ok := g.placeUnit(next); !ok {
			panic("not ok?")
		}
		games[i] = g
	}

	return games
}

func (g *Game) String() string {
	return fmt.Sprintf(`Game{
	Score:         %f,
	B:             %s,
	units:         %+v,
	lcg:           %+v,
	numUnits:      %d,
	unitsSent:     %d,
	currUnit:      %+v,
	previousMoves: %v,
}`, g.Score(), g.B.StringLevel(2), g.units, g.lcg, g.numUnits, g.unitsSent, g.currUnit, g.previousMoves)
}

func (g *Game) LockUnit(u *Unit) {
	for _, c := range u.Members {
		g.B.MarkFilled(c)
	}
}

func (g *Game) NextUnit() (*Unit, bool) {
	if g.unitsSent >= g.numUnits {
		return nil, false
	}

	rand := g.lcg.Next()
	idx := int(rand) % len(g.units)
	templUnit := &g.units[idx]

	// Do a deep copy of the chosen Unit.
	r := templUnit.DeepCopy()

	g.unitsSent++

	return r, true
}

func (g *Game) placeUnit(u *Unit) bool {
	l, r := u.Bounds()

	uwidth := r.X - l.X + 1
	extraSpace := g.B.Width - uwidth

	// Center up, rounding down, leaving less space on the left.
	rightShift := extraSpace / 2
	// If the leftmost cell is not at zero, we don't need to shift as much.
	rightShift -= l.X
	for i := range u.Members {
		u.Members[i].X += rightShift
	}
	u.Pivot.X += rightShift

	return g.B.IsValid(u)
}

// updateScore computes the new Game moves score, and remembers linesCleared as
// previous lines cleared. The power score is computed on-demand with Score()
// or PowerScore().
func (g *Game) updateScore(linesCleared int) {
	ls := float64(linesCleared)
	lsOld := float64(g.previousLinesCleared)
	size := float64(g.currUnit.Size())

	points := size + 100.0*(1.0+ls)*ls/2.0

	var lineBonus int
	if lsOld > 1 {
		lineBonus = int((lsOld - 1.0) * points / 10.0)
	}

	moveScore := points + float64(lineBonus)

	g.moveScore += moveScore
	g.previousLinesCleared = linesCleared
}

// Count occurrences of sep within s, allowing for overlap, which the spec
// allows for phrases of power.
func CountOverlap(s, sep string) (count int) {
	for i := strings.Index(s, sep); i >= 0; i = strings.Index(s, sep) {
		count++
		s = s[i+1:]
	}
	return
}

func (g *Game) updatePowerCount() (score int) {
	s := g.Commands.String()

	for _, p := range normalizedPhrases {
		if len(s) < len(p) {
			continue
		}

		end := s[len(s) - len(p):]
		if end == p {
			c, ok := g.powerWordCount[p]
			if !ok {
				c = 0
			}
			g.powerWordCount[p] = c + 1
		}
	}

	return
}

// PowerScore computes the phrase of power score from the currently completed
// moves.
func (g *Game) PowerScore() (score int) {
	for p, n := range g.powerWordCount {
		score += 2 * len(p) * n
		if n > 0 {
			score += 300
		}
	}

	return
}

// PowerScore() but with final phrases
func (g *Game) PowerFinalScore() (score int) {
	s := g.FinalCommands.String()

	for _, p := range powerPhrases {
		n := CountOverlap(s, p)
		score += 2 * len(p) * n
		if n > 0 {
			score += 300
		}
	}

	return
}

// Score returns the total game score so far.
func (g *Game) Score() float64 {
	return g.moveScore + float64(g.PowerScore())
}

// Score() but with final phrases
func (g *Game) FinalScore() float64 {
	return g.moveScore + float64(g.PowerFinalScore())
}

// Rewrite commands with "final" power Phrases instead of normalized ones
func (g *Game) WriteFinalCommands() {
	s := g.Commands.String() // Copy starting commands
	again := true
	for again {
		again = false
		// Loop over all power phrases, try to replace them each once
		for i, np := range normalizedPhrases {
			pp := powerPhrases[i]
			if strings.Index(s, np) >= 0 {
				again = true // go through the list again
				s = strings.Replace(s, np, pp, 1)
			}
		}
	}
	g.FinalCommands = make(Commands, len(s)) // Done
	for i, c := range s {
		g.FinalCommands[i] = Command(c)
	}
}

// Update returns a bool indicating whether a piece was locked, the game is done, and err to indicate and error (backwards move).
func (g *Game) Update(c Command) (bool, bool, error) {
	d, ok := commandToDirection[c]
	if !ok {
		return false, true, fmt.Errorf("unknown command %c", c)
	}

	if d == NOP {
		return false, false, nil
	}

	var moved *Unit
	isRot := d == CCW || d == CW
	if isRot {
		moved = g.currUnit.Rotate(d == CCW)
	} else {
		moved = g.currUnit.Translate(d)
	}

	// We cannot move into the same position, so we must add the current
	// position (before the translate which we are about to check) to the
	// list of moves to check against. We do this now, rather than the end
	// of the previous Update call to prevent a bad rotation on the first
	// move.
	previousMoves := append(g.previousMoves, g.currUnit.DeepCopy())
	if moved.OverlapsAny(previousMoves) {
		return false, true, fmt.Errorf("moved unit from %+v to %+v and it overlaps with a previous move!", g.currUnit, moved)
	}

	// No more error beyond this point, record the command and previous
	// moves.
	g.Commands = append(g.Commands, c)
	g.previousMoves = previousMoves

	g.updatePowerCount()

	if g.B.IsValid(moved) {
		g.currUnit = moved
		return false, false, nil
	}

	g.LockUnit(g.currUnit)

	linesCleared := g.B.ClearRows()
	g.updateScore(linesCleared)

	nextUnit, ok := g.NextUnit()
	if !ok {
		// Game is done.
		return true, true, nil
	}

	if ok := g.placeUnit(nextUnit); !ok {
		// Game is done.
		return true, true, nil
	}

	g.previousMoves = g.previousMoves[:0]
	g.currUnit = nextUnit
	return true, false, nil
}
