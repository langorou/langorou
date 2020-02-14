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
	count uint8
	race  race
}

type state struct {
	grid [][]cell
}
