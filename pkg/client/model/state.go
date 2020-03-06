package model

import (
	"fmt"
)

type Race uint8

const (
	// Neutral Race
	Neutral Race = iota
	// Ally Race
	Ally
	// Enemy Race
	Enemy
)

func (r Race) Opponent() Race {
	if r == Ally {
		return Enemy
	} else if r == Enemy {
		return Ally
	} else {
		panic(fmt.Sprintf("Opponent asked for Race: '%s' this should not happen ", []string{"Neutral", "Ally", "Enemy"}[r]))
	}
}

type Cell struct {
	Count uint8
	Race  Race
}

func (c *Cell) IsEmpty() bool {
	return c.Count == 0
}

// State represents a game state, disclaimer we should NOT modify Grid directly, use SetCell, IncreaseCell and DecreaseCell
// methods instead, Grid is only available to ease it's reading process
type State struct {
	Grid            map[Coordinates]Cell
	Height          uint8
	Width           uint8
	time            uint8
	CumulativeScore float64
	Diffs map[Coordinates]Cell
}

func NewState(height uint8, width uint8) State {
	return State{
		Grid:            map[Coordinates]Cell{},
		Height:          height,
		Width:           width,
		time:            0,
		CumulativeScore: 0,
		Diffs: map[Coordinates]Cell{},
	}
}

// Copy copies a state, incrementing the cumulative score and the time
func (s State) Copy(advanceTime bool) State {
	newDiffs := make(map[Coordinates]Cell, len(s.Diffs))
	for k, v := range s.Diffs {
		newDiffs[k] = v
	}
	score := s.CumulativeScore
	time := s.time
	if advanceTime {
		score += (1 - float64(s.time)/1000) * s.allies()
		time += 1
	}
	return State{Grid: s.Grid, Height: s.Height, Width: s.Width, CumulativeScore: score, time: time, Diffs: newDiffs}
}

func (s *State) SetCell(pos Coordinates, race Race, count uint8) {
	// If we set a cell to 0, remove it except if it's neutral (because the HUM message from the server does this)
	if count == 0 && race != Neutral {
		s.EmptyCell(pos)
	}
	s.Grid[pos] = Cell{Race: race, Count: count}
}

func (s *State) GetCellWithDiff(pos Coordinates) (Cell, bool) {
	c, ok := s.Diffs[pos]
	if ok {
		return c, true
	}

	c, ok = s.Grid[pos]
	return c, ok
}

func (s *State) DiffSetCell(pos Coordinates, race Race, count uint8) {
	s.Diffs[pos] = Cell{Race: race, Count: count}
}


func (s *State) DiffDecreaseCell(pos Coordinates, count uint8) {
	c, ok := s.GetCellWithDiff(pos)
	if !ok {
		panic(fmt.Sprintf("Tried to decrease population at non existing cell: %+v", pos))
	}

	if c.Count < count {
		panic(fmt.Sprintf("Invalid move ! From pos: %+v, race: %v, current count: %d, move count: %d", pos, c.Race, c.Count, count))
	}

	s.Diffs[pos] = Cell{Race:c.Race, Count:c.Count -count}
}

func (s *State) EmptyCell(pos Coordinates) {
	delete(s.Grid, pos)
}

func (s State) allies() float64 {
	// TODO: could be computed iteratively from setcell and friends
	count := uint8(0)

	for _, c := range s.Grid {
		if c.Race == Ally {
			count += c.Count
		}
	}

	return float64(count)
}
