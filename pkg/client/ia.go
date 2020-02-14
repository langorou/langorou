package client

type PotentialState struct {
	s    state
	prob float64
}

// applyMov to a state
func evaluateMoveOutcome(s state, mov Move) []PotentialState {
	// we assume the mov is "safe"/correct

	startRace := s[mov.Start.X][mov.Start.Y].race
	endRace := s[mov.End.X][mov.End.Y].race
	endCount := s[mov.End.X][mov.End.Y].count

	if endCount == 0 {
		// nobody on there
		newState := s.deepCopy()
		newState[mov.Start.X][mov.Start.Y].count -= float64(mov.N)
		newState[mov.End.X][mov.End.Y] = cell{count: float64(mov.N), race: startRace}

		return []PotentialState{PotentialState{s: newState, prob: 1}}

	} else if startRace == endRace {
		// same race, no battle
		newState := s.deepCopy()
		newState[mov.Start.X][mov.Start.Y].count -= float64(mov.N)
		newState[mov.End.X][mov.End.Y].count += float64(mov.N)

		return []PotentialState{PotentialState{s: newState, prob: 1}}
	}
	// Fight with the enemy or neutral

	// We use a float here for later computation
	var isNeutral float64 = 0
	if endRace == Neutral {
		isNeutral = 1
	}

	P := getProba(float64(mov.N), endCount, isNeutral == 1)

	if P == 1 {
		// We surely win, same as "nobody there"
		newState := s.deepCopy()
		newState[mov.Start.X][mov.Start.Y].count -= float64(mov.N)
		newState[mov.End.X][mov.End.Y] = cell{
			count: float64(mov.N) + float64(isNeutral)*endCount, // if we totally win against Neutral, we convert all of them
			race:  startRace,
		}
		return []PotentialState{PotentialState{s: newState, prob: 1}}
	}

	winState := s.deepCopy()
	winState[mov.Start.X][mov.Start.Y].count -= float64(mov.N)
	winState[mov.End.X][mov.End.Y] = cell{count: float64(mov.N)*P + float64(isNeutral)*endCount, // For each ally, he has P chance to survive. Against neutral, we have P chance to convert them
		race: startRace,
	}

	lossState := s.deepCopy()
	lossState[mov.Start.X][mov.Start.Y].count -= float64(mov.N)
	lossState[mov.End.X][mov.End.Y] = cell{count: float64(mov.N) * (1 - P), race: endRace} // For each enemy, he has 1-P chance to survive

	return []PotentialState{
		PotentialState{s: winState, prob: P},
		PotentialState{s: lossState, prob: 1 - P},
	}

	// A neutral or enemy

}

func evaluateNextState(s state, movs []Move) state {
	// All moves are necesseraly independent, we'll have less 2**len(moves) possible states (max 2 outputs per move)

	// Make a deep copy of the state -> is it necessary ?
	newState := s.deepCopy()

	for _, mov := range movs {
		_ = mov
	}

	return newState
}

func scoreState(potSta PotentialState) float64 {

	// Apply the change on the state

	h := 0.
	for _, row := range potSta.s {
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

	return h * potSta.prob
}
