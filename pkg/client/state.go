package client

import "fmt"

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
}

func (s state) deepCopy() state {
	newGrid := make(map[Coordinates]cell, len(s.grid))
	for k, v := range s.grid {
		newGrid[k] = v
	}
	return state{grid: newGrid, height: s.height, width: s.width}
}
