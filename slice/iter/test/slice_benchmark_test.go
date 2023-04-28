package test

import (
	"testing"

	"github.com/m4gshm/gollections/slice"
)

func Benchmark_IsValidIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := slice.IsValidIndex(5, 0)
		r = slice.IsValidIndex(5, 5)
		r = slice.IsValidIndex(5, -1)
		_ = r
	}
}


func Benchmark_CanIterateByRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := slice.CanIterateByRange(slice.IterNoStarted, 5, 4)
		r = slice.CanIterateByRange(slice.IterNoStarted, 5, 6)
		r = slice.CanIterateByRange(slice.IterNoStarted, 5, slice.IterNoStarted)
		_ = r
	}
}
