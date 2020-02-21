package client

type race uint

const (
	// Empty race
	Empty race = iota
	// Neutral race
	Neutral
	// Ally race
	Ally
	// Enemy race
	Enemy
)

type cell struct {
	count uint8 // float used for potential state, used in heuristic computation
	race  race
}

func (c *cell) isEmpty() bool {
	return c.count == 0
}

type state [][]cell

func (s state) deepCopy() state {
	newState := make([][]cell, len(s))
	for i := range s {
		newState[i] = make([]cell, len(s[i]))
		copy(newState[i], s[i])
	}
	return newState
}
