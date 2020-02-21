package client

// Adapted from github.com/Succo/twilight, but we should use float since we evaluate probability of winning.

// getProba of winning for the attaquant 1 with an effectif E1, agains E2
// E2 might be Neutral
func getProba(E1, E2 uint8, E2isNeutral bool) float64 {

	var cste float64
	if E2isNeutral {
		cste = 1
	} else {
		cste = 1.5
	}

	// True by property
	if float64(E1) >= cste*float64(E2) {
		return 1
	}

	if E1 == E2 {
		return 0.5
	} else if E1 < E2 {
		return float64(E1) / (2 * float64(E2))
	} else {
		return (float64(E1) / float64(E2)) - 0.5
	}
}

/*
// INFO: github.com/Succo/twilight version
// getProba reimplements the getProba logic from Board.cs in the C# implementation
func getProba(E1, E2 int, involveHumans bool) float64 {
	if E1 == E2 {
		return 0.5
	}
	var cste float64
	if involveHumans {
		cste = 1
	} else {
		cste = 1.5
	}

	// True by property
	if float64(E1) >= cste*float64(E2) {
		return 1
	}

	var x0, y0 float64
	x1 := float64(E2)
	y1 := 0.5
	if E1 < E2 {
		x0 = 0
		y0 = 0
		return (y0 - y1) / (x0 - x1) * float64(E1)
	} else {
		x0 = cste * float64(E2)
		y0 = 1
		m := (y0 - y1) / (x0 - x1)
		c := 1 - m*x0
		return m*float64(E2) + c
	}
}
*/
