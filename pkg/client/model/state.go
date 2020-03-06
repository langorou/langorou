package model

import (
	"fmt"
)

type Race string

const (
	// Neutral Race
	Neutral = "Neutral"
	// Ally Race
	Ally = "Ally"
	// Enemy Race
	Enemy = "Enemy"
)

func (r Race) Opponent() Race {
	if r == Ally {
		return Enemy
	} else if r == Enemy {
		return Ally
	} else {
		panic(fmt.Sprintf("Opponent asked for Race: '%s' this should not happen ", r))
	}
}

type Cell struct {
	Count uint8 // float used for potential State, used in heuristic computation
	Race  Race
}

// A Cell with coordinates
type CCell struct {
	Coords Coordinates
	Cell
}

func (c *Cell) IsEmpty() bool {
	return c.Count == 0
}

// State represents a game state, disclaimer we should NOT modify Grid directly, use SetCell, IncreaseCell and DecreaseCell
// methods instead, Grid is only available to ease it's reading process
type State struct {
	Grid            []CCell
	Height          uint8
	Width           uint8
	time            uint8
	CumulativeScore float64
}

func NewState(height uint8, width uint8) State {
	return State{
		Grid:            []CCell{},
		Height:          height,
		Width:           width,
		time:            0,
		CumulativeScore: 0,
	}
}

// Copy copies a state, incrementing the cumulative score and the time
func (s State) Copy(advanceTime bool) State {
	newGrid := make([]CCell, len(s.Grid))
	copy(newGrid, s.Grid)
	score := s.CumulativeScore
	time := s.time
	if advanceTime {
		score += (1 - float64(s.time)/1000) * s.allies()
		time += 1
	}
	return State{Grid: newGrid, Height: s.Height, Width: s.Width, CumulativeScore: score, time: time}
}

func (s *State) cellIndex(pos Coordinates) int {
	for i, cc := range s.Grid {
		if cc.Coords == pos {
			return i
		}
	}

	return -1
}

func (s *State) GetCell(pos Coordinates) Cell {
	idx := s.cellIndex(pos)
	if idx < 0 {
		return Cell{}
	}

	return s.Grid[idx].Cell
}

func (s *State) SetCell(pos Coordinates, race Race, count uint8) {
	cc := CCell{pos, Cell{Race: race, Count: count}}

	idx := s.cellIndex(pos)
	if idx < 0 && count == 0 {
		return
	} else if idx < 0 {
		s.Grid = append(s.Grid, cc)
	} else {
		s.Grid[idx] = cc
	}
}

func (s *State) IncreaseCell(pos Coordinates, count uint8) {
	idx := s.cellIndex(pos)
	if idx < 0 {
		panic(fmt.Sprintf("Tried to increase population at non existing cell: %+v", pos))
	}
	s.Grid[idx].Count += count
}

func (s *State) DecreaseCell(pos Coordinates, count uint8) {
	idx := s.cellIndex(pos)
	if idx < 0 {
		panic(fmt.Sprintf("Tried to decrease population at non existing cell: %+v", pos))
	}
	c := s.Grid[idx]

	if c.Count < count {
		panic(fmt.Sprintf("Invalid move ! From pos: %+v, race: %v, current count: %d, move count: %d", pos, c.Race, c.Count, count))
	}

	s.Grid[idx].Count -= count
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
