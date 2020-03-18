package client

import (
	"github.com/langorou/langorou/pkg/client/model"
	"math"
)

const (
	posInfinity = math.MaxFloat64
	negInfinity = -math.MaxFloat64
)

type resultType uint8

type result struct {
	depth uint8
	score float64
	typ   resultType
	coup  model.Coup
}

const (
	lower resultType = iota
	upper
	exact
)

type transpositionTable struct {
	t            map[uint64]result
	hits, misses uint
}

func (t *transpositionTable) get(hash uint64, depth uint8) (result, bool) {
	rec, ok := t.t[hash]
	if ok && rec.depth >= depth {
		t.hits += 1
		return rec, true
	}

	t.misses += 1

	return rec, false
}

func (t *transpositionTable) save(hash uint64, coup model.Coup, value float64, depth uint8, alpha float64, beta float64) {
	s := result{coup: coup, score: value, depth: depth, typ: exact}
	if alpha > value {
		s.typ = lower
	} else if value > beta {
		s.typ = upper
	}

	t.t[hash] = s
}

func (h *Heuristic) findBestCoup(state *model.State, maxDepth uint8) (coup model.Coup, score float64) {
	for depth := uint8(1); depth <= maxDepth; depth++ {
		tt := &transpositionTable{map[uint64]result{}, 0, 0}
		coup, score = h.alphabeta(tt, state, model.Ally, negInfinity, posInfinity, 0, depth)
		// log.Printf("misses: %d, hits: %d, hit ratio: %f, entries: %d", tt.misses, tt.hits, float64(tt.hits)/(float64(tt.hits+tt.misses)), len(tt.t))
	}

	// TODO: use transposition table to use the best move found at previous depth

	return coup, score
}

// alphabeta computes the best coup going at most at depth depth
func (h *Heuristic) alphabeta(tt *transpositionTable, state *model.State, race model.Race, alpha float64, beta float64, depth uint8, maxDepth uint8) (model.Coup, float64) {
	bestCoup := model.Coup{}

	hash := state.Hash(race)

	rec, cached := tt.get(hash, depth)
	if cached {
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

	if depth >= maxDepth || state.GameOver() { // Max depth reached or game is over
		value := h.scoreState(state)
		tt.save(hash, bestCoup, value, depth, alpha, beta)
		return bestCoup, value
	}

	coups := generateCoups(state, race)
	if len(coups) == 0 { // or no more moves found
		value := h.scoreState(state)
		tt.save(hash, bestCoup, value, depth, alpha, beta)
		return bestCoup, value
	}

	// Chose if we want to maximize (us) or minimize (enemy) our score
	value := negInfinity
	f := math.Max
	if race == model.Enemy {
		value = posInfinity
		f = math.Min
	}

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

		if f(value, score) == score { // score >= value if max playing or value >= score if min playing
			value = score
			bestCoup = coup
			// log.Printf("better value found %f: depth: %d, race: %v", value, depth, race)
		}

		// Check for possible cuts
		if race == model.Enemy {
			// alpha cut
			if alpha > value {
				break
			}
			beta = f(beta, value)
		} else {
			// beta cut
			if value > beta {
				break
			}
			alpha = f(alpha, value)
		}
	}

	tt.save(hash, bestCoup, value, depth, alpha, beta)
	return bestCoup, value
}
