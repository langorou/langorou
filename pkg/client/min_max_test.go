package client

import (
	"fmt"
	"github.com/langorou/langorou/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testDepth = 5

var testHeuristic = NewHeuristic(NewDefaultHeuristicParameters())

func TestMinMax(t *testing.T) {

	t.Run("case1", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 10N | XXX | 06N
		// XXX | 08A | XXX
		// 12E | XXX | XXX

		startState := model.NewState(3, 3)
		startState.SetCell(model.Coordinates{}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 2}, model.Neutral, 6)
		startState.SetCell(model.Coordinates{Y: 2}, model.Enemy, 12)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 8)

		coup, _ := testHeuristic.findBestCoup(startState, testDepth)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 1, Y: 1},
			N:     8,
			End:   model.Coordinates{X: 2, Y: 0},
		}}, coup)
	})

	t.Run("case2", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// 08A | 06N | XXX
		// XXX | 10N | XXX
		// XXX | XXX | XXX
		// Far away: 8Enemy

		// Here our only chance to win is to try to steal the group of 10 neutrals then go for the group of 6 neutrals
		startState := model.NewState(10, 10)
		startState.SetCell(model.Coordinates{}, model.Ally, 8)
		startState.SetCell(model.Coordinates{X: 1}, model.Neutral, 6)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Neutral, 10)
		startState.SetCell(model.Coordinates{X: 8, Y: 8}, model.Enemy, 8)

		coup, _ := testHeuristic.findBestCoup(startState, testDepth)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 0, Y: 0},
			N:     8,
			End:   model.Coordinates{X: 1, Y: 0},
		}}, coup)
	})

	t.Run("case3", func(t *testing.T) {
		// N neutral, A ally, E enemy
		// XXX | XXX | XXX | XXX
		// XXX | 68A | XXX | XXX
		// XXX | XXX | 07N | XXX
		// XXX | XXX | XXX | XXX
		// Far away: 75Enemy

		startState := model.NewState(10, 10)
		startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 68)
		startState.SetCell(model.Coordinates{X: 2, Y: 2}, model.Neutral, 7)
		startState.SetCell(model.Coordinates{X: 7, Y: 4}, model.Enemy, 75)

		coup, _ := testHeuristic.findBestCoup(startState, testDepth)
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 1, Y: 1},
			N:     68,
			End:   model.Coordinates{X: 2, Y: 2},
		}}, coup)
	})
}

func TestSimulationAllyNeutral(t *testing.T) {
	startState := model.NewState(2, 2)
	startState.SetCell(model.Coordinates{}, model.Neutral, 10)
	startState.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 20)

	t.Run("sure win", func(t *testing.T) {
		s := startState.Copy(false)

		coup := model.Coup{
			{Start: model.Coordinates{X: 1, Y: 1}, End: model.Coordinates{}, N: 11},
		}

		potentialStates := applyCoup(s, model.Ally, coup, 1)

		assert.Len(t, potentialStates, 1)
		expected := model.NewState(2, 2)
		expected.SetCell(model.Coordinates{}, model.Ally, 21)
		expected.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 9)
		assert.EqualValues(t, 1, potentialStates[0].probability)
		assert.Equal(t, expected.Grid, potentialStates[0].s.Grid)
	})

	t.Run("unsure win", func(t *testing.T) {
		s := startState.Copy(false)

		coup := model.Coup{
			{Start: model.Coordinates{X: 1, Y: 1}, End: model.Coordinates{}, N: 8},
		}

		potentialStates := applyCoup(s, model.Ally, coup, 1)

		assert.Len(t, potentialStates, 2)
		expected1 := model.NewState(2, 2)
		expected1.SetCell(model.Coordinates{}, model.Ally, 7)
		expected1.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 12)
		assert.Equal(t, 0.4, potentialStates[0].probability)
		assert.Equal(t, expected1.Grid, potentialStates[0].s.Grid)

		expected2 := model.NewState(2, 2)
		expected2.SetCell(model.Coordinates{}, model.Neutral, 6)
		expected2.SetCell(model.Coordinates{X: 1, Y: 1}, model.Ally, 12)
		assert.EqualValues(t, 0.6, potentialStates[1].probability)
		assert.Equal(t, expected2.Grid, potentialStates[1].s.Grid)
	})

	t.Run("minmax decision", func(t *testing.T) {
		s := startState.Copy(false)
		s.SetCell(model.Coordinates{X: 1}, model.Enemy, 15)
		coup, _ := testHeuristic.findBestCoup(s, testDepth)

		// Probability 5/6 of winning if we attack the enemy directly
		// only 3/4 if we get the villagers but the enemy attacks us after
		assert.Equal(t, model.Coup{model.Move{
			Start: model.Coordinates{X: 1, Y: 1},
			N:     20,
			End:   model.Coordinates{X: 1, Y: 0},
		}}, coup)
	})
}

func BenchmarkAB(b *testing.B) {
	cplxState := model.GenerateComplicatedState()

	for _, depth := range []uint8{2, 3, 4, 5, 6} {
		b.Run(fmt.Sprintf("depth%d_complex", depth), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				testHeuristic.findBestCoup(cplxState.Copy(false), depth)
			}
		})
	}

	smplState := model.GenerateSimpleState()

	for _, depth := range []uint8{2, 3, 4, 5, 6} {
		b.Run(fmt.Sprintf("depth%d_simple", depth), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				testHeuristic.findBestCoup(smplState.Copy(false), depth)
			}
		})
	}
}

func BenchmarkHeuristic(b *testing.B) {
	startState := model.GenerateComplicatedState()
	for n := 0; n < b.N; n++ {
		testHeuristic.scoreState(startState)
	}
}

func BenchmarkHash(b *testing.B) {
	startState := model.GenerateComplicatedState()

	for n := 0; n < b.N; n++ {
		startState.Hash(model.Ally)
	}
}

func TestGenerateCoups(t *testing.T) {
	startState := model.GenerateComplicatedState()
	generateCoups(startState, model.Ally)
}
