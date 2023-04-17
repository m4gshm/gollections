package test

import (
	"testing"

	"github.com/m4gshm/gollections/iter/impl/iter"
)

func Benchmark_IsValidIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := iter.IsValidIndex(5, 0)
		r = iter.IsValidIndex(5, 5)
		r = iter.IsValidIndex(5, -1)
		_ = r
	}
}

func Benchmark_IsValidIndex2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := iter.IsValidIndex2(5, 0)
		r = iter.IsValidIndex2(5, 5)
		r = iter.IsValidIndex2(5, -1)
		_ = r
	}
}

func Benchmark_CanIterateByRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := iter.CanIterateByRange(iter.NoStarted, 5, 4)
		r = iter.CanIterateByRange(iter.NoStarted, 5, 6)
		r = iter.CanIterateByRange(iter.NoStarted, 5, iter.NoStarted)
		_ = r
	}
}
