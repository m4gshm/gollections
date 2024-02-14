//go:build goexperiment.rangefunc

package collection

import (
	"testing"

	oset "github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/collection/immutable/vector"
)


func Benchmark_Loop_ImmutableOrderSet_go_1_22_direct(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.All(func(v int) bool {
					casee.load(v)
					return true
				})
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderSet_go_1_22(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for v := range c.All {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderSet_go_1_22_2(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for v := range c.All2 {
					casee.load(v)
				}
			}
		})
	}
}


func Benchmark_Loop_ImmutableVector_go_1_22(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range c.All {
					casee.load(v)
				}
			}
		})
	}
}

