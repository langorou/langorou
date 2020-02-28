package client

// negamaxAlpha returns a score and a turn
func negamaxAlpha(state state, alpha float64, race race, depth uint8) (Coup, float64) {
	bestTurn := Coup{}
	maxScore := -1000000.0

	if depth <= 0 {
		potentialState := PotentialState{s: state, prob: 1}
		return bestTurn, scoreState(potentialState)
	}

	for _, coup := range getCoups(state, race) {
		potentialStates := evaluateMoveOutcomes(state, race, coup)

		score := 0.0
		for _, potentialState := range potentialStates {
			_, tmpScore := negamaxAlpha(potentialState.s, maxScore, race.opponent(), depth-1)
			score += tmpScore * potentialState.prob
		}
		score = -score
		if score > maxScore {
			maxScore = score
			bestTurn = coup
		}
		//todo save previously modified cells by the applyMoves function
		//we use a deep copy for now
		//undoMoves(potentialState.s, race, coup) dans le cas
		//ou pas de bataille aleatoire pour eviter les copies

		if maxScore > -alpha {
			break
		}

	}
	return bestTurn, maxScore
}

func getMovesFromCell(coord Coordinates) []Move {
	//TODO returns possible moves from a cell
	//we don't consider the splits for now
	return nil
}

// getMoves returns a list of moves
func getCoups(state state, race race) []Coup {
	// TODO
	return nil
}

func undoCoup(state state, player uint8, coup Coup) {
	//TODO
}
