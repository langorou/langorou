package client

// negamaxAlpha returns a score and a turn
func negamaxAlpha(state state, alpha float64, race race, depth uint8) (Coup, float64) {
	bestTurn := Coup{}
	maxScore := -1000000.0

	if depth <= 0 {
		potentialState := PotentialState{s: state, probability: 1}
		return bestTurn, scoreState(potentialState)
	}

	for _, coup := range generateCoups(state, race) {
		potentialStates := applyCoup(state, race, coup)

		score := 0.0
		for _, potentialState := range potentialStates {
			_, tmpScore := negamaxAlpha(potentialState.s, maxScore, race.opponent(), depth-1)
			score += tmpScore * potentialState.probability
		}
		score = -score
		if score > maxScore {
			maxScore = score
			bestTurn = coup
		}

		// TODO: if memory is an issue, we could try to implement an undoCoup function, that would allow to reduce the number of copies made

		if maxScore > -alpha {
			break
		}

	}
	return bestTurn, maxScore
}
