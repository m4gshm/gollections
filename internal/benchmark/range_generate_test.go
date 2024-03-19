package benchmark

import (
	"testing"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

func Benchmark_LoopRange_Iterating(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := loop.Range(0, 100000)
				for v, ok := next(); ok; v, ok = next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_SliceRange_Iterating(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				integers := slice.Range(0, 100000)
				for i := range integers {
					casee.load(integers[i])
				}
			}
		})
	}
}
