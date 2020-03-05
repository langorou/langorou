package client

import (
	"fmt"
)

type race string

const (
	// Neutral race
	Neutral = "Neutral"
	// Ally race
	Ally = "Ally"
	// Enemy race
	Enemy = "Enemy"
)

func (r race) opponent() race {
	if r == Ally {
		return Enemy
	} else if r == Enemy {
		return Ally
	} else {
		panic(fmt.Sprintf("opponent asked for race: %s this should not happen ", r))
	}
}

type cell struct {
	count uint8 // float used for potential state, used in heuristic computation
	race  race
}

func (c *cell) isEmpty() bool {
	return c.count == 0
}

type state struct {
	grid   map[Coordinates]cell
	height uint8
	width  uint8
	time uint8
	cumScore float64
}

func (s state) deepCopy(dt uint8) state {
	newGrid := make(map[Coordinates]cell, len(s.grid))
	for k, v := range s.grid {
		newGrid[k] = v
	}
	sc := s.cumScore
	if dt != 0 {
		sc += ( 1 - float64(s.time)/1000) * s.allies()
	}
	return state{grid: newGrid, height: s.height, width: s.width, cumScore: sc, time: s.time + dt}
}

func (s state) allies() float64 {
	count := uint8(0)

	for _, c := range s.grid {
		if c.race == Ally {
			count += c.count
		}
	}

	return float64(count)
}