package model

import (
	"sort"
)

// Coordinates represents coordinates on the grid
type Coordinates struct {
	X uint8
	Y uint8
}

func uint8Sub(a, b uint8) uint8 {
	if a > b {
		return a - b
	}

	return b - a
}

func uint8Max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func (c1 *Coordinates) Distance(c2 Coordinates) float64 {
	return float64(uint8Max(uint8Sub(c1.X, c2.X), uint8Sub(c1.Y, c2.Y)))
}

// Changes is an update of the board sent by the server
type Changes struct {
	Coords  Coordinates
	Neutral uint8
	Ally    uint8
	Enemy   uint8
}

// Move is an allowed move
type Move struct {
	Start Coordinates
	N     uint8
	End   Coordinates
}

// Coup represents a list of moves/actions, it implements the sort.Interface to sort by target cells
type Coup []Move

func (coup Coup) Len() int {
	return len(coup)
}

func (coup Coup) Less(i, j int) bool {
	return coup[i].End.Y > coup[j].End.Y && coup[i].End.X > coup[j].End.X
}

func (coup Coup) Swap(i, j int) {
	coup[i], coup[j] = coup[j], coup[i]
}

func SortCoupsByQuickScore(coups []Coup, s State) {
	cbqs := coupByQuickScore{coups, s}
	sort.Sort(cbqs)
}

// quickScore scores a move quickly (quicker than the heuristic) just to be able to sort coups by their score
// and have a more effective minmax
// It simply computes the count of the battles taking places:
// No battle will return 0
// Nice battles will be > 0
// Bad battles will be < 0
func (coup Coup) quickScore(state State) float64 {
	battles := map[Coordinates]struct{}{}
	score := 0.

	for _, move := range coup {
		end := state.Grid[move.End]
		if !end.IsEmpty() && end.Race != Ally {
			if _, ok := battles[move.End]; !ok {
				score -= float64(end.Count)
			}
			score += float64(move.N)
		}
	}

	return score
}

// coupByQuickScore is used to sort coups by quick score
type coupByQuickScore struct {
	c []Coup
	s State
}

func (cb coupByQuickScore) Len() int {
	return len(cb.c)
}

func (cb coupByQuickScore) Less(i, j int) bool {
	return cb.c[i].quickScore(cb.s) > cb.c[j].quickScore(cb.s)
}

func (cb coupByQuickScore) Swap(i, j int) {
	cb.c[i], cb.c[j] = cb.c[j], cb.c[i]
}
