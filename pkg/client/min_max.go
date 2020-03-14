package client

import (
	"github.com/langorou/langorou/pkg/client/model"
	"math"
)

const (
	posInfinity = math.MaxFloat64
	negInfinity = -math.MaxFloat64
)

type save struct {
	depth uint8
	score float64
	typ   uint8
	coup  model.Coup
}

const (
	lower uint8 = iota
	upper
	exact
)

func (h *Heuristic) findBestCoup(state model.State, maxDepth uint8) (coup model.Coup, score float64) {
	tt := map[uint64]save{}

	for depth := uint8(1); depth <= maxDepth; depth++ {
		coup, score = h.alphabeta(tt, state, model.Ally, negInfinity, posInfinity, 0, depth)
	}

	return coup, score
}

// alphabeta computes the best coup going at most at depth depth
func (h *Heuristic) alphabeta(tt map[uint64]save, state model.State, race model.Race, alpha float64, beta float64, depth uint8, maxDepth uint8) (model.Coup, float64) {
	bestCoup := model.Coup{}

	hash := state.Hash()

	if rec, ok := tt[hash]; ok && rec.depth == depth {
		if rec.typ == exact {
			return rec.coup, rec.score
		} else if rec.typ == lower {
			alpha = math.Max(alpha, rec.score)
		} else if rec.typ == upper {
			beta = math.Min(beta, rec.score)
		}

		if alpha >= beta {
			return rec.coup, rec.score
		}
	}

	if depth >= maxDepth { // Max depth reached
		score := h.scoreState(state)
		tt[hash] = save{coup: bestCoup, score: score, depth: depth, typ: exact}
		return bestCoup, score
	}

	coups := generateCoups(state, race)
	if len(coups) == 0 { // or no more moves found
		score := h.scoreState(state)
		tt[hash] = save{coup: bestCoup, score: score, depth: depth, typ: exact}
		return bestCoup, h.scoreState(state)
	}

	if rec, ok := tt[hash]; ok && len(rec.coup) != 0 {
		// Put the current move first
		// TODO: duplicated
		coups = append([]model.Coup{rec.coup}, coups...)
	}

	// Chose if we want to maximize (us) or minimize (enemy) our score
	value := negInfinity
	f := math.Max
	if race == model.Enemy {
		value = posInfinity
		f = math.Min
	}

	// Sort by killer moves
	// model.SortCoupsByQuickScore(coups, state)

	// for each generated coup, we compute the list of potential outcomes and compute an average score
	// weighted by the probabilities of these potential outcomes
	for _, coup := range coups {
		outcomes := applyCoup(state, race, coup, h.WinThreshold)

		score := 0.
		// log.Printf("depth: %d", depth)
		for _, outcome := range outcomes {
			_, tmpScore := h.alphabeta(tt, outcome.s, race.Opponent(), alpha, beta, depth+1, maxDepth)
			score += tmpScore * outcome.probability
		}

		// log.Printf("cumulative: %f, findBestCoup score: %f at depth: %d for race: %v and coup: %+v, grid: %+v, potential: %+v", state.cumScore, score, depth, race, coup, state.grid, outcomes)

		if f(value, score) == score { // score >= value if max playing or value >= score if min playing
			value = score
			bestCoup = coup
			// log.Printf("better value found %f: depth: %d, race: %v", value, depth, race)
		}

		// Check for possible cuts
		if race == model.Enemy {
			// alpha cut
			if alpha >= value {
				break
			}
			beta = f(beta, value)
		} else {
			// beta cut
			if value >= beta {
				break
			}
			alpha = f(alpha, value)
		}
	}

	s := save{depth: depth, coup: bestCoup, score: value, typ: exact}
	if alpha >= value {
		s.typ = lower
	} else if value >= beta {
		s.typ = upper
	}

	tt[hash] = s

	return bestCoup, value
}
