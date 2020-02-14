package client

// func evaluateNextState(s state, movs []Move) state {
// 	nextGrid = state{}
// }

func scoreState(s state) float64 {

	// Apply the change on the state
	newState := s

	h := 0.
	for _, row := range newState.grid {
		for _, c := range row {
			switch c.race {
			case Empty:
				// nothing
			case Neutral:
				// nothing
			case Ally:
				h += float64(c.count)
			case Enemy:
				h -= float64(c.count)
			}
		}
	}

	return h
}
