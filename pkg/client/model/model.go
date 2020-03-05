package model

import "math"

// Coordinates represents coordinates on the grid
type Coordinates struct {
	X uint8
	Y uint8
}

func (c1 *Coordinates) Distance(c2 Coordinates) float64 {
	x1, y1 := float64(c1.X), float64(c1.Y)
	x2, y2 := float64(c2.X), float64(c2.Y)

	return math.Max(math.Abs(x1-x2), math.Abs(y1-y2))
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
