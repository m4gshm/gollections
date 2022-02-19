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

var cases = []benchCase{/*{"high", HighLoad}, */{"low", LowLoad}}

func Benchmark_HasNext_Iterator_Immutable_Vector(b *testing.B) {
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

func Benchmark_HasNext_Iterator_Immutable_Vector_Impl(b *testing.B) {
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

func Benchmark_GetNext_Iterator_Immutable_Vector_Impl(b *testing.B) {
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

func Benchmark_GetPrev_Iterator_Immutable_Vector_Impl(b *testing.B) {
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

func Benchmark_HasNext_Iterator_Mutable_Vector_Impl(b *testing.B) {
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

func Benchmark_HasPrev_Immutable_Vector_Impl(b *testing.B) {
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

func Benchmark_HasNext_Iterator_WrapSlice(b *testing.B) {
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

func Benchmark_HasNext_Iterator_Impl(b *testing.B) {
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

func Benchmark_ForRange_EmbeddedSlice(b *testing.B) {
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

func Benchmark_ForByIndex_EmbeddedSlice(b *testing.B) {
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

func Benchmark_WrapMap_HasNext(b *testing.B) {
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

func Benchmark_NewKVHasNext(b *testing.B) {
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

func Benchmark_EmbeddedMap_For(b *testing.B) {
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

func Benchmark_ForEach_Immutable_Vector(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func Benchmark_ForRange_of_Collect_Immutable_Vector_Impl_Values(b *testing.B) {
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

func Benchmark_ForEach_Immutable_Set(b *testing.B) {
	c := set.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func Benchmark_ForEach_Immutable_Set_Impl(b *testing.B) {
	c := set.New(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func Benchmark_ForEach_Immutable_OrdererSet(b *testing.B) {
	c := oset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func Benchmark_ForEach_Immutable_OrdererSet_Impl(b *testing.B) {
	c := oset.New(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func Benchmark_ForRange_of_Collect_Immutable_OrdererSet_Values(b *testing.B) {
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

func Benchmark_ForRange_of_Collect_Immutable_OrdererSet_Impl_Values(b *testing.B) {
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

func Benchmark_ForEach_Mutable_OrdererSet_Impl(b *testing.B) {
	c := moset.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}
