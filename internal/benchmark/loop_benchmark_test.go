package benchmark

import (
	"testing"

	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/map_"
	moset "github.com/m4gshm/gollections/mutable/oset"
	mvector "github.com/m4gshm/gollections/mutable/vector"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	max        = 100000
	values     = range_.Of(1, max)
	ResultInt  = 0
	threshhold = max / 2
)

func HighLoad(v int) {
	ResultInt = v *
		v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v *
		v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v *
		v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v *
		v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v *
		v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v *
		v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v
}

func LowLoad(v int) {
	ResultInt = v * v
}

type benchCase struct {
	name string
	load func(int)
}

var cases = []benchCase{ /*{"high", HighLoad}, */ {"low", LowLoad}}

func Benchmark_Loop_ImmutableOrderSet_FirstNext(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it, v, ok := c.First(); ok; v, ok = it.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableOrderSet_LastPrev(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it, v, ok := c.Last(); ok; v, ok = it.Prev() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_IterNextNext(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				iter := c.Iter()
				for v, ok := iter.Next(); ok; v, ok = iter.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_HeadHasNextGetNext(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := c.Head(); it.HasNext(); {
					casee.load(it.GetNext())
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_HeadNextNext(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				h := c.Head()
				for v, ok := h.Next(); ok; v, ok = h.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_FirstNext(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it, v, ok := c.First(); ok; v, ok = it.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_TailPrevPrev(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := c.Tail()
				for v, ok := it.Prev(); ok; v, ok = it.Prev() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_LastPrev(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it, v, ok := c.Last(); ok; v, ok = it.Prev() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_MutableVector_FirstNext(b *testing.B) {
	c := mvector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it, v, ok := c.First(); ok; v, ok = it.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_MutableVector_HeadNext(b *testing.B) {
	c := mvector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := ptr.Of(c.Head()).Next
				for v, ok := next(); ok; v, ok = next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_Loop_ImmutableVector_TailHasPrevGetPrev(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := c.Tail(); it.HasPrev(); {
					casee.load(it.GetPrev())
				}
			}
		})
	}
}

func Benchmark_Loop_Slice_Wrap_NextNext(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := slice.NewIter(values).Next
				for v, ok := next(); ok; v, ok = next() {
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

func Benchmark_Loop_NewKV_NextNextNext(b *testing.B) {
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

func Benchmark_Loop_MutableOrdererSet_FirstNext(b *testing.B) {
	c := moset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i, e, ok := c.First(); ok; e, ok = i.Next() {
					casee.load(e)
				}
			}
		})
	}
}

func Benchmark_Loop_MutableOrdererSet_Head(b *testing.B) {
	c := moset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				i := c.Head()
				for e, ok := i.Next(); ok; e, ok = i.Next() {
					casee.load(e)
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
