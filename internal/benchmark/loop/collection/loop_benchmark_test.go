package collection

import (
	"testing"

	oset "github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/collection/immutable/set"
	"github.com/m4gshm/gollections/collection/immutable/vector"
	moset "github.com/m4gshm/gollections/collection/mutable/ordered/set"
	"github.com/m4gshm/gollections/convert/ptr"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	max        = 100000
	values     = range_.Closed(1, max)
	resultInt  = 0
	threshhold = max / 2
)

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

func Benchmark_Loop_ImmutableOrderSet_All(b *testing.B) {
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

func Benchmark_Loop_Slice_Loop_All(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := loop.Of(values...)
				for v := range l.All {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_Slice_NewHead_HasNextGetNext(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := slice.NewHead(values); it.HasNext(); {
					casee.load(it.GetNext())
				}
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

func Benchmark_Loop_Slice_Embedded_ForByRangeIndex(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range values {
					v := values[j]
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_Slice_Embedded_ForByIndex(b *testing.B) {
	l := len(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < l; j++ {
					v := values[j]
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_Map_Embedded_ForByKeyValueRange(b *testing.B) {
	values := map[int]struct{}{}
	for i := 0; i < max; i++ {
		values[i] = struct{}{}
	}
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for k, v := range values {
					_ = v
					casee.load(k)
				}
			}
		})
	}
}

func Benchmark_Loop_Map_NewIter_NextNext(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				iterator := map_.NewIter(values)
				for k, _, ok := iterator.Next(); ok; k, _, ok = iterator.Next() {
					casee.load(k)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_ForEach(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_ForRangeSlice(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range c.Slice() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableSet_ForEach(b *testing.B) {
	c := set.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderedSet_ForEach(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderedSet_ForRangeSlice(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range c.Slice() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_MutableOrdererSet_ForEach(b *testing.B) {
	c := moset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}

func Benchmark_Loop_Loop_RangeClosed_ForEach(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				loop.RangeClosed(1, max).ForEach(casee.load)
			}
		})
	}
}

func Benchmark_Loop_Loop_Of_ForEach(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				loop.Of(values...).ForEach(casee.load)
			}
		})
	}
}

func Benchmark_Loop_Slice_Head_Next(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i, v, ok := ptr.Of(slice.NewHead(values)).Crank(); ok; v, ok = i.Next() {
					casee.load(v)
				}
			}
		})
	}
}
