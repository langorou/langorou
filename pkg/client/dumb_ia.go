package client

import (
	"log"
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
	log.Printf("coups: %v", coups)

	if len(coups) == 0 {
		return nil
	}

	return coups[rand.Intn(len(coups))]
}
