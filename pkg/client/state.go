package client

type race uint

const (
	// Empty race
	Empty race = iota
	// Villager race
	Villager
	// Ally race
	Ally
	// Enemy race
	Enemy
)

type cell struct {
	count uint
	race  race
}

type state struct {
	grid [][]cell
}
