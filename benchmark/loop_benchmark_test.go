package benchmark

import (
	"testing"

	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/iter"
	impliter "github.com/m4gshm/gollections/iter/impl/iter"
	moset "github.com/m4gshm/gollections/mutable/oset"
	mvector "github.com/m4gshm/gollections/mutable/vector"
	"github.com/m4gshm/gollections/ptr"
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

func BenchmarkLoopImmutableOrderSetFirstNext(b *testing.B) {
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

func BenchmarkLoopImmutableOrderSetLastPrev(b *testing.B) {
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

func BenchmarkLoopImmutableVectorBeginNextNext(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := c.Begin().Next
				for v, ok := next(); ok; v, ok = next() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkLoopImmutableVectorHeadHasNextGetNext(b *testing.B) {
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

func BenchmarkLoopImmutableVectorHeadNextNext(b *testing.B) {
	c := vector.Of(values...)
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

func BenchmarkLoopImmutableVectorFirstNext(b *testing.B) {
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

func BenchmarkLoopImmutableVectorTailPrevPrev(b *testing.B) {
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

func BenchmarkLoopImmutableVectorLastPrev(b *testing.B) {
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

func BenchmarkLoopMutableVectorFirstNext(b *testing.B) {
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

func BenchmarkLoopMutableVectorHeadNext(b *testing.B) {
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

func BenchmarkLoopImmutableVectorTailHasPrevGetPrev(b *testing.B) {
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

func BenchmarkLoopSliceWrapNextNext(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := iter.Wrap(values).Next
				for v, ok := next(); ok; v, ok = next() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkLoopSliceNewHeadHasNextGetNext(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := impliter.NewHead(values); it.HasNext(); {
					casee.load(it.GetNext())
				}
			}
		})
	}
}

func BenchmarkLoopSliceEmbeddedForByRange(b *testing.B) {
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

func BenchmarkLoopSliceEmbeddedForByRangeIndex(b *testing.B) {
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

func BenchmarkLoopSliceEmbeddedForByIndex(b *testing.B) {
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

func BenchmarkLoopMapEmbeddedForByKeyValueRange(b *testing.B) {
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

func BenchmarkLoopNewKVNextNextNext(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				iterator := impliter.NewEmbedMapKV(values)
				for k, _, ok := iterator.Next(); ok; k, _, ok = iterator.Next() {
					casee.load(k)
				}
			}
		})
	}
}

func BenchmarkImmutableVectorForEachLoop(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}

func BenchmarkImmutableVectorCollecAndLoop(b *testing.B) {
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

func BenchmarkImmutableSetForEachLoop(b *testing.B) {
	c := set.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}

func BenchmarkImmutableOrderedSetForEachLoop(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}

func BenchmarkImmutableOrderedSetCollectAndLoop(b *testing.B) {
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

func BenchmarkMutableOrdererSetForEachLoop(b *testing.B) {
	c := moset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(casee.load)
			}
		})
	}
}
