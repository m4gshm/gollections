//go:build goexperiment.rangefunc

package collection2

import (
	"testing"

	oset "github.com/m4gshm/gollections/collection/immutable/ordered/set"
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
