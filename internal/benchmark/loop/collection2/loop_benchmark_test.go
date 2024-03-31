package collection2

import (
	"testing"

	oset "github.com/m4gshm/gollections/collection/immutable/ordered/set"
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

var cases = []benchCase{{"high", HighLoad}, {"low", LowLoad}}

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

func Benchmark_Loop_ImmutableOrderSet_FirstNext(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for it, v, ok := c.First(); ok; v, ok = it.Next() {
					casee.load(v)
				}
			}
			b.StopTimer()
		})
	}
}

func Benchmark_Loop_ImmutableOrderSet_FirstNext2(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it, v, ok := c.First()
				for ok {
					casee.load(v)
					v, ok = it.Next()
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderSet_Head_HasNext_GetNext(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := c.Head()
				for it.HasNext() {
					casee.load(it.GetNext())
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderSet_Loop_Crank_Next(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for next, v, ok := c.Loop().Crank(); ok; v, ok = next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderSet_Loop_Next(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := c.Loop()
				for {
					v, ok := next()
					if !ok {
						break
					}
					casee.load(v)
				}
			}
		})
	}
}
