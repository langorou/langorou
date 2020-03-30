package model

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestHashing(t *testing.T) {
	s1 := NewState(10, 10)
	s1.SetCell(Coordinates{2, 2}, Ally, 75)
	s1.SetCell(Coordinates{7, 4}, Enemy, 75)

	s2 := NewState(10, 10)
	s2.SetCell(Coordinates{0, 0}, Ally, 68)
	s2.SetCell(Coordinates{2, 2}, Neutral, 7)
	s2.SetCell(Coordinates{7, 4}, Enemy, 75)

	assert.NotEqual(t, s1.Hash(Ally, nil), s2.Hash(Ally, nil))
}

func TestSort(t *testing.T) {
	length := 64
	for i := 0; i < 100; i++ {
		arr := make(sortableU32, length)
		arr2 := make(sortableU32, length)
		for idx := range arr {
			arr[idx] = uint32(rand.Int())
			arr2[idx] = arr[idx]
		}

		sort.Sort(arr)
		sortQuick(arr2)
		assert.Equal(t, arr, arr2)
	}
}
