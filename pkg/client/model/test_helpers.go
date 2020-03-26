package model

func GenerateComplicatedState() *State {
	s := NewState(10, 10)
	s.SetCell(Coordinates{X: 1, Y: 1}, Ally, 68)
	s.SetCell(Coordinates{X: 2, Y: 2}, Neutral, 7)
	s.SetCell(Coordinates{X: 2, Y: 7}, Neutral, 18)
	s.SetCell(Coordinates{X: 3, Y: 3}, Ally, 11)
	s.SetCell(Coordinates{X: 5, Y: 7}, Neutral, 3)
	s.SetCell(Coordinates{X: 5, Y: 8}, Neutral, 4)
	s.SetCell(Coordinates{X: 5, Y: 9}, Neutral, 18)
	s.SetCell(Coordinates{X: 8, Y: 1}, Enemy, 2)
	s.SetCell(Coordinates{X: 9, Y: 0}, Enemy, 53)

	return s
}

func GenerateSimpleState() *State {
	s := NewState(5, 5)
	s.SetCell(Coordinates{X: 0, Y: 0}, Ally, 68)
	s.SetCell(Coordinates{X: 4, Y: 2}, Neutral, 4)
	s.SetCell(Coordinates{X: 0, Y: 4}, Enemy, 68)

	return s
}
