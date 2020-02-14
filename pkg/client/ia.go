package client

func evaluateHeuristic(s state, changes []Changes) float64 {

	// Apply the change on the state
	newState := s

	h := 0.
	for _, row := range newState.grid {
		for _, c := range row {
			switch c.race {
			case Empty:
				// nothing
			case Villager:
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
