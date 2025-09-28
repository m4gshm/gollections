package collection2

import (
	"testing"

	oset "github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	max    = 10000
	values = range_.Closed(1, max)
)

var resultInt = 0

func LowLoad(v int) {
	resultInt = v * v * v
}

func HighLoad(v int) {
	resultInt = v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v
}

type benchCase struct {
	name string
	load func(int)
}

var cases = []benchCase{ /*{"high", HighLoad},*/ {"low", LowLoad}}

func Benchmark_Loop_ImmutableOrderSet_ForEach(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}

func Benchmark_Loop_Slice_Seq_ForByRange(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for v := range seq.Of(values...) {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_Slice_Embedded_ForByRange(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range values {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderSet_ForRange_All(b *testing.B) {
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

func Benchmark_Loop_Slice_Loop_NextNext(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := loop.Of(values...)
				for v, ok := next(); ok; v, ok = next() {
					casee.load(v)
				}
			}
		})
	}
}
