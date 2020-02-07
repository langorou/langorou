package client

type race int

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
