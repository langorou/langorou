package client

import "log"

type MinMaxIA struct {
	depth uint8
}

var _ IA = &MinMaxIA{}

func NewMinMaxIA(depth uint8) *MinMaxIA {
	return &MinMaxIA{
		depth: depth,
	}
}

func (m *MinMaxIA) Play(state state) Coup {
	coup, score := minimax(state, Ally, m.depth)
	log.Printf("MinMaxIA computed a coup with score: %f", score)
	return coup
}
