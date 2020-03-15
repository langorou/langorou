package model

import (
	"sort"
	"testing"
)

func BenchmarkSortStd(b *testing.B) {
	s := GenerateComplicatedState()

	buf := make(sortableU32, 0, 128)

	for i := 0; i < b.N; i++ {
		s.packedU32(buf)
		sort.Sort(buf)
	}
}

func BenchmarkSortClassic(b *testing.B) {
	s := GenerateComplicatedState()

	buf := make([]uint32, 0, 128)

	for i := 0; i < b.N; i++ {
		s.packedU32(buf)
		classicSort(buf)
	}
}

func BenchmarkSortCustomQSort(b *testing.B) {
	s := GenerateComplicatedState()

	buf := make([]uint32, 0, 128)

	for i := 0; i < b.N; i++ {
		s.packedU32(buf)
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
