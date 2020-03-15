package model

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func BenchmarkSortStd(b *testing.B) {
	s := GenerateComplicatedState()

	buf := make(sortableU32, 0, 128)

	for i := 0; i < b.N; i++ {
		buf = s.packedU32(buf)
		sort.Sort(buf)
	}
}

func BenchmarkSortClassic(b *testing.B) {
	s := GenerateComplicatedState()

	buf := make([]uint32, 0, 128)

	for i := 0; i < b.N; i++ {
		buf = s.packedU32(buf)
		classicSort(buf)
	}
}

func BenchmarkSortCustomQSort(b *testing.B) {
	s := GenerateComplicatedState()

	buf := make([]uint32, 0, 128)

	for i := 0; i < b.N; i++ {
		buf = s.packedU32(buf)
		sortQuick(buf)
	}
}

func classicSort(buf []uint32) {
	for j, ej := range buf {
		for k, ek := range buf {
			if ej < ek && j > k {
				buf[j], buf[k] = ek, ej
			}
		}
	}
}

func TestHashing(t *testing.T) {
	s1 := NewState(10, 10)
	s1.SetCell(Coordinates{2, 2}, Ally, 75)
	s1.SetCell(Coordinates{7, 4}, Enemy, 75)


	s2 := NewState(10, 10)
	s2.SetCell(Coordinates{0, 0}, Ally, 68)
	s2.SetCell(Coordinates{2, 2}, Neutral, 7)
	s2.SetCell(Coordinates{7, 4}, Enemy, 75)

	assert.NotEqual(t, s1.Hash(Ally), s2.Hash(Ally))
}
