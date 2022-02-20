package benchmark

import (
	"testing"

	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/it"
	impliter "github.com/m4gshm/gollections/it/impl/it"
	moset "github.com/m4gshm/gollections/mutable/oset"
	mvector "github.com/m4gshm/gollections/mutable/vector"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	max       = 100000
	values    = range_.Of(1, max)
	ResultInt = 0
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

func BenchmarkLoopImmutableOrderSetImplFirstNext(b *testing.B) {
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

func BenchmarkLoopImmutableOrderSetImplLastPrev(b *testing.B) {
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
				it := c.Begin()
				for v, ok := it.Next(); ok; v, ok = it.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkLoopImmutableVectorImplHeadHasNextGetNext(b *testing.B) {
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

func BenchmarkLoopImmutableVectorImplHeadNextNext(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := c.Head()
				for v, ok := it.Next(); ok; v, ok = it.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkLoopImmutableVectorImplFirstNext(b *testing.B) {
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

func BenchmarkLoopImmutableVectorImplTailPrevPrev(b *testing.B) {
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

func BenchmarkLoopImmutableVectorImplLastPrev(b *testing.B) {
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

func BenchmarkLoopMutableVectorImplHeadHeadNext(b *testing.B) {
	c := mvector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := c.Head()
				for v, ok := it.Next(); ok; v, ok = it.Next() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkLoopImmutableVectorImplTailHasPrevGetPrev(b *testing.B) {
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

func BenchmarkLoopSliceWrapHasNextGetNext(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := it.Wrap(values); it.HasNext(); {
					casee.load(it.GetNext())
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

func BenchmarkLoopSliceEmbeddedForByIndex(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(values); j++ {
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
				iter := impliter.NewKV(values)
				for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
					casee.load(k)
				}
			}
		})
	}
}

func BenchmarkImmutableVectorForEach(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkLoopImmutableVectorCollectEmbeddedForByRange(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range c.Collect() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkLoopImmutableSetByOfForEach(b *testing.B) {
	c := set.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkLoopImmutableSetByNewForEach(b *testing.B) {
	c := set.New(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkLoopImmutableOrdererSetByOfForEach(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkLoopImmutableOrdererSetByNewForEach(b *testing.B) {
	c := oset.New(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkLoopImmutableOrdererSetCollectByEmbeddedForByRange(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range c.Collect() {
					casee.load(v)
				}
			}
		})
	}
}

func LoopMutableOrdererSetByOfForEach(b *testing.B) {
	c := moset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}
