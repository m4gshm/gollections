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

func BenchmarkHasNextIteratorImmutableVector(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := c.Begin()
				for v, ok := it.GetNext(); ok; v, ok = it.GetNext() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkHeadHasNextNextIteratorImmutableVectorImpl(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := c.Head(); it.HasNext(); {
					casee.load(it.Next())
				}
			}
		})
	}
}

func BenchmarkHeadGetNextIteratorImmutableVectorImpl(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := c.Head()
				for v, ok := it.GetNext(); ok; v, ok = it.GetNext() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkGoForwardGetNextIteratorImmutableVectorImpl(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it, v, ok := c.GoForward(); ok; v, ok = it.GetNext() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkTailGetPrevIteratorImmutableVectorImpl(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := c.Tail()
				for v, ok := it.GetPrev(); ok; v, ok = it.GetPrev() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkGoBackGetPrevIteratorImmutableVectorImpl(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it, v, ok := c.GoBack(); ok; v, ok = it.GetPrev() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkHeadGetNextIteratorMutableVectorImpl(b *testing.B) {
	c := mvector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				it := c.Head()
				for v, ok := it.GetNext(); ok; v, ok = it.GetNext() {
					casee.load(v)
				}
			}
		})
	}
}

func BenchmarkTailHasPrevImmutableVectorImpl(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := c.Tail(); it.HasPrev(); {
					casee.load(it.Prev())
				}
			}
		})
	}
}

func BenchmarkHasNextIteratorWrapSlice(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := it.Wrap(values); it.HasNext(); {
					casee.load(it.Next())
				}
			}
		})
	}
}

func BenchmarkHasNextIteratorImpl(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := impliter.NewHead(values); it.HasNext(); {
					casee.load(it.Next())
				}
			}
		})
	}
}

func BenchmarkForRangeEmbeddedSlice(b *testing.B) {
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

func BenchmarkForByIndexEmbeddedSlice(b *testing.B) {
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

func BenchmarkWrapMapHasNext(b *testing.B) {
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

func BenchmarkNewKVHasNext(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				iter := impliter.NewKV(values)
				for k, _, ok := iter.GetNext(); ok; k, _, ok = iter.GetNext() {
					casee.load(k)
				}
			}
		})
	}
}

func BenchmarkEmbeddedMapFor(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
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

func BenchmarkForEachImmutableVector(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkForRangeofCollectImmutableVectorImplValues(b *testing.B) {
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

func BenchmarkForEachImmutableSet(b *testing.B) {
	c := set.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkForEachImmutableSetImpl(b *testing.B) {
	c := set.New(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkForEachImmutableOrdererSet(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkForEachImmutableOrdererSetImpl(b *testing.B) {
	c := oset.New(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func BenchmarkForRangeOfCollectImmutableOrdererSetValues(b *testing.B) {
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

func BenchmarkForRangeOfCollectImmutableOrdererSetImplValues(b *testing.B) {
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

func BenchmarkForEachMutableOrdererSetImpl(b *testing.B) {
	c := moset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}
