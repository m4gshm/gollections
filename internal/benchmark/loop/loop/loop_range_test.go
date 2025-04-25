package loop

import (
	"testing"

	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/convert/ptr"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
)

var max = 10000

func HighLoad(v int) {
	resultInt = v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v
}

var resultInt = 0

func LowLoad(v int) {
	resultInt = v * v * v
}

func MidLoad(v int) {
	resultInt = v * v * v * v * v * v * v * v * v
}

type benchCase struct {
	name string
	load func(int)
}

var cases = []benchCase{ /*{"high", HighLoad}, {"mid", MidLoad},*/ {"low", LowLoad}}

func Benchmark_SliceRange_Iterating(b *testing.B) {
	integers := slice.Range(0, max)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range integers {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_SeqRange_Iterating(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for v := range seq.Range(0, max) {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_LoopRange_Iterating(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := loop.Range(0, max)
				for v, ok := next(); ok; v, ok = next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_LoopRange_Iterating2(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := loop.Range(0, max)
				v, ok := next()
				for ok {
					casee.load(v)
					v, ok = next()
				}
			}
		})
	}
}

func Benchmark_LoopRange_Iterating3(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := loop.Range(0, max)
				for v := range next.All {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Slice_Iter_Iterating(b *testing.B) {
	integers := slice.Range(0, max)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it, v, ok := ptr.Of(slice.NewHead(integers)).Crank(); ok; v, ok = it.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Slice_Iter_Iterating2(b *testing.B) {
	integers := slice.Range(0, max)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := slice.NewHead(integers)
				v, ok := it.Next()
				for ok {
					casee.load(v)
					v, ok = it.Next()
				}
			}
		})
	}
}

func Benchmark_Slice_Mutable_Iter_Iterating(b *testing.B) {
	integers := slice.Range(0, max)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := mutable.NewHead(&integers, nil)
				v, ok := it.Next()
				for ok {
					casee.load(v)
					v, ok = it.Next()
				}
			}
		})
	}
}
