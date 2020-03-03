package client

// negamaxAlpha returns a score and a turn
func negamaxAlpha(state state, alpha float64, race race, depth uint8) (Coup, float64) {
	bestCoup := Coup{}
	maxScore := -1e64

	// TODO: maybe move it after if we know that we already hit the depth limit
	coups := generateCoups(state, race)

	if depth <= 0 || len(coups) == 0 {
		potentialState := potentialState{s: state, probability: 1}
		score := scoreState(potentialState, race)
		return bestCoup, score
	}

	for _, coup := range coups {
		potentialStates := applyCoup(state, race, coup)

		score := 0.0
		for _, potentialState := range potentialStates {
			_, tmpScore := negamaxAlpha(potentialState.s, maxScore, race.opponent(), depth-1)
			score += tmpScore * potentialState.probability
		}

		score = -score

		if score > maxScore {
			maxScore = score
			bestCoup = coup
		}

		// TODO: if memory is an issue, we could try to implement an undoCoup function, that would allow to reduce the number of copies made

		// TODO: uncomment me
		// if maxScore > - alpha {
		// 	break
		// }

	}
	return bestCoup, maxScore
}
