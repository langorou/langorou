package client

import (
	"math/rand"
	"time"
)

type DumbIA struct{}

var _ IA = &DumbIA{}

func NewDumbIA() *DumbIA {
	return &DumbIA{}
}

func (dia *DumbIA) Play(state state) Coup {
	// Simulate computation
	time.Sleep(time.Second)

	coups := generateCoups(state, Ally)

	return coups[rand.Intn(len(coups))]
}
