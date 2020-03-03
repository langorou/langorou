package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const alphaTest = 10

func TestNegamax(t *testing.T) {

	t.Run("case1", func(t *testing.T) {
		// N neutral, A ally, E ennemy
		// 10N | XXX | 06N
		// XXX | 08A | XXX
		// 12E | XXX | XXX

		// Here our only chance to win is to try to steal the group of 10 neutrals then go for the group of 6 neutrals
		startState := state{
			grid: map[Coordinates]cell{
				{X: 0, Y: 0}: {
					race:  Neutral,
					count: 10,
				},
				{X: 2, Y: 0}: {
					race:  Neutral,
					count: 6,
				},
				{X: 0, Y: 2}: {
					race:  Enemy,
					count: 12,
				},
				{X: 1, Y: 1}: {
					race:  Ally,
					count: 8,
				},
			},
			height: 2,
			width:  2,
		}

		coup, _ := negamaxAlpha(startState, alphaTest, Ally, 5)
		assert.Equal(t, Coup{Move{
			Start: Coordinates{X: 1, Y: 1},
			N:     8,
			End:   Coordinates{X: 0, Y: 0},
		}}, coup)
	})
}

func TestSimulationAllyNeutral(t *testing.T) {
	startState := state{
		grid: map[Coordinates]cell{
			{X: 0, Y: 0}: {
				race:  Neutral,
				count: 10,
			},
			{X: 1, Y: 1}: {
				race:  Ally,
				count: 20,
			},
		},
		height: 2,
		width:  2,
	}

	t.Run("sure win", func(t *testing.T) {
		s := startState.deepCopy()

		coup := Coup{
			{Start: Coordinates{X: 1, Y: 1}, End: Coordinates{0, 0}, N: 10},
		}

		potentialStates := applyCoup(s, Ally, coup)

		assert.Len(t, potentialStates, 1)
		assert.Equal(t, potentialState{
			probability: 1,
			s: state{
				grid: map[Coordinates]cell{
					{X: 0, Y: 0}: {
						race:  Ally,
						count: 20,
					},
					{X: 1, Y: 1}: {
						race:  Ally,
						count: 10,
					},
				},
				height: 2,
				width:  2,
			},
		}, potentialStates[0])
	})

	t.Run("unsure win", func(t *testing.T) {
		s := startState.deepCopy()

		coup := Coup{
			{Start: Coordinates{X: 1, Y: 1}, End: Coordinates{0, 0}, N: 8},
		}

		potentialStates := applyCoup(s, Ally, coup)

		assert.Len(t, potentialStates, 2)
		assert.Equal(t, potentialState{
			probability: 0.4,
			s: state{
				grid: map[Coordinates]cell{
					{X: 0, Y: 0}: {
						race:  Ally,
						count: 7,
					},
					{X: 1, Y: 1}: {
						race:  Ally,
						count: 12,
					},
				},
				height: 2,
				width:  2,
			},
		}, potentialStates[0])

		assert.Equal(t, potentialState{
			probability: 0.6,
			s: state{
				grid: map[Coordinates]cell{
					{X: 0, Y: 0}: {
						race:  Neutral,
						count: 6,
					},
					{X: 1, Y: 1}: {
						race:  Ally,
						count: 12,
					},
				},
				height: 2,
				width:  2,
			},
		}, potentialStates[1])
	})

	t.Run("negamax decision", func(t *testing.T) {
		s := startState.deepCopy()
		s.grid[Coordinates{X: 1, Y: 0}] = cell{
			race:  Enemy,
			count: 15,
		}
		coup, _ := negamaxAlpha(s, alphaTest, Ally, 3)

		assert.Equal(t, Coup{Move{
			Start: Coordinates{X: 1, Y: 1},
			N:     20,
			End:   Coordinates{X: 0, Y: 0},
		}}, coup)
	})
}
