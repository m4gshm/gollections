package benchmark

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/it"
	impliter "github.com/m4gshm/gollections/it/impl/it"

	moset "github.com/m4gshm/gollections/mutable/oset"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	max       = 100000
	values    = range_.Of(1, max)
	ResultStr = strconv.Itoa(ResultInt)
	ResultInt = 0
)

func HighLoad(v int) {
	ResultStr = strconv.Itoa(v)
}

func LowLoad(v int) {
	ResultInt = v * v
}

type benchCase struct {
	name string
	load func(int)
}

var cases = []benchCase{/*{"high", HighLoad}, */{"low", LowLoad}}

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

func Benchmark_ForRange_of_Collect_Immutable_Vector_Values(b *testing.B) {
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

func Benchmark_ForRange_of_Collect_Immutable_Vector_Impl_Values(b *testing.B) {
	c := vector.New(values)
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

func Benchmark_ForEach_Mutable_OrdererSet(b *testing.B) {
	c := moset.Convert(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func Benchmark_ForEach_Mutable_OrdererSet_Impl(b *testing.B) {
	c := moset.Convert(values)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.ForEach(func(v int) { casee.load(v) })
			}
		})
	}
}

func Benchmark_HasNextGet_Iterator_Immutable_Vector(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := c.Begin(); it.HasNext(); {
					casee.load(it.Next())
				}
			}
		})
	}
}

func Benchmark_HasNextGet_Iterator_Immutable_Vector_Impl(b *testing.B) {
	c := vector.Of(values...)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := c.Iter(); it.HasNext(); {
					casee.load(it.Next())
				}
			}
		})
	}
}

func Benchmark_HasNextGet_Iterator_WrapSlice(b *testing.B) {
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

func Benchmark_HasNextGet_Iterator_Impl(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := impliter.New(values); it.HasNext(); {
					casee.load(it.Next())
				}
			}
		})
	}
}

func Benchmark_HasNextGet_Iterator_Impl_Point(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for it := impliter.NewP(&values); it.HasNext(); {
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

func Benchmark_WrapMap_HasNextGet(b *testing.B) {
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

func Benchmark_NewKVHasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for iter := impliter.NewKV(values); iter.HasNext(); {
					k, v := iter.Next()
					_ = v
					casee.load(k)
				}
			}
		})
	}
}

func Benchmark_NewReflectKVHasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for iter := impliter.NewReflectKV(values); iter.HasNext(); {
					k, v := iter.Next()
					_ = v
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
