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
}

func NewState(height uint8, width uint8) State {
	return State{
		Grid:            map[Coordinates]Cell{},
		Height:          height,
		Width:           width,
		time:            0,
		CumulativeScore: 0,
	}
}

// Copy copies a state, incrementing the cumulative score and the time
func (s State) Copy(advanceTime bool) State {
	newGrid := make(map[Coordinates]Cell, len(s.Grid))
	for k, v := range s.Grid {
		newGrid[k] = v
	}
	score := s.CumulativeScore
	time := s.time
	if advanceTime {
		score += (1 - float64(s.time)/1000) * s.allies()
		time += 1
	}
	return State{Grid: newGrid, Height: s.Height, Width: s.Width, CumulativeScore: score, time: time}
}

func (s *State) SetCell(pos Coordinates, race Race, count uint8) {
	// If we set a cell to 0, remove it except if it's neutral (because the HUM message from the server does this)
	if count == 0 && race != Neutral {
		s.EmptyCell(pos)
	}
	s.Grid[pos] = Cell{Race: race, Count: count}
}

func (s *State) IncreaseCell(pos Coordinates, count uint8) {
	c, ok := s.Grid[pos]
	if !ok {
		panic(fmt.Sprintf("Tried to increase population at non existing cell: %+v", pos))
	}
	c.Count += count
	s.Grid[pos] = c
}

func (s *State) DecreaseCell(pos Coordinates, count uint8) {
	c, ok := s.Grid[pos]
	if !ok {
		panic(fmt.Sprintf("Tried to decrease population at non existing cell: %+v", pos))
	}

	if c.Count == count {
		// If cell is going to be empty, let's remove it
		s.EmptyCell(pos)
		return
	} else if c.Count < count {
		panic(fmt.Sprintf("Invalid move ! From pos: %+v, race: %v, current count: %d, move count: %d", pos, c.Race, c.Count, count))
	}

	c.Count -= count
	s.Grid[pos] = c
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
