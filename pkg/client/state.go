package client

type race uint

const (
	// Empty race
	Empty race = iota
	// Villager race
	Villager
	// Werewolf race
	Werewolf
	// Vampire race
	Vampire
)

type cell struct {
	count uint
	race  race
}

type state struct {
	grid [][]cell
}
