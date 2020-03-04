package client

import (
	"math"
)

func minimax(state state, race race, depth uint8) (Coup, float64) {
	return alphabeta(state, race, math.SmallestNonzeroFloat64, math.MaxFloat64, depth)
}

// minimax computes the best coup going at most at depth depth
func alphabeta(state state, race race, alpha float64, beta float64, depth uint8) (Coup, float64) {
	bestCoup := Coup{}

	// Max depth reached
	if depth <= 0 {
		return bestCoup, scoreState(state, Ally)
	}

	coups := generateCoups(state, race)
	// No moves found
	if len(coups) == 0 {
		return bestCoup, scoreState(state, Ally)
	}

	// Chose if we want to maximize (us) or minimize (enemy) our score
	value := math.SmallestNonzeroFloat64
	f := math.Max
	if race == Enemy {
		value = math.MaxFloat64
		f = math.Min
	}

	// for each generated coup, we compute the list of potential outcomes and compute an average score
	// weighted by the probabilities of these potential outcomes
	for _, coup := range coups {
		outcomes := applyCoup(state, race, coup)

		score := 0.
		// log.Printf("depth: %d", depth)
		for _, outcome := range outcomes {
			_, tmpScore := alphabeta(outcome.s, race.opponent(), alpha, beta, depth-1)
			score += tmpScore * outcome.probability
		}

		// log.Printf("minimax score: %f at depth: %d for race: %v and coup: %+v, grid: %+v, potential: %+v", score, depth, race, coup, state.grid, outcomes)

		newValue := f(value, score)
		if newValue == score {
			value = newValue
			bestCoup = coup
		}

		// Check for possible cuts
		if race == Enemy {
			// alpha cut
			if alpha >= value {
				return bestCoup, value
			}
			beta = f(beta, value)
		} else {
			// beta cut
			if value >= beta {
				return bestCoup, value
			}
			alpha = f(alpha, value)
		}
	}

	return bestCoup, value
}
