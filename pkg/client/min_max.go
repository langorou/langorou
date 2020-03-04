package client

import (
	"math"
)

// minimax computes the best coup going at most at depth depth
func minimax(state state, race race, depth uint8) (Coup, float64) {
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
	value := -1e64
	f := math.Max
	if race == Enemy {
		value = 1e64
		f = math.Min
	}

	// for each generated coup, we compute the list of potential outcomes and compute an average score
	// weighted by the probabilities of these potential outcomes
	for _, coup := range coups {
		outcomes := applyCoup(state, race, coup)

		score := 0.
		// log.Printf("depth: %d", depth)
		for _, outcome := range outcomes {
			_, tmpScore := minimax(outcome.s, race.opponent(), depth-1)
			score += tmpScore * outcome.probability
		}

		// log.Printf("minimax score: %f at depth: %d for race: %v and coup: %+v, grid: %+v, potential: %+v", score, depth, race, coup, state.grid, outcomes)

		newValue := f(value, score)
		if newValue == score {
			value = newValue
			bestCoup = coup
		}
	}

	return bestCoup, value
}
