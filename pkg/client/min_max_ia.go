package client

import "log"

type MinMaxIA struct {
	alpha float64
	depth uint8
}

var _ IA = &MinMaxIA{}

func NewMinMaxIA(alpha float64, depth uint8) *MinMaxIA {
	return &MinMaxIA{
		alpha: alpha,
		depth: depth,
	}
}

func (m *MinMaxIA) Play(state state) Coup {
	coup, score := negamaxAlpha(state, m.alpha, Ally, m.depth)
	log.Printf("MinMaxIA computed a coup with score: %f", score)
	return coup
}
