package benchmark

import (
	"testing"

	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/it"
	impliter "github.com/m4gshm/gollections/it/impl/it"
	mset "github.com/m4gshm/gollections/mutable/set"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	max    = 100000
	values = range_.Of(1, max)
	result = 0
)

func Benchmark_ForEach_Immutable_Vector(b *testing.B) {
	c := vector.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Immutable_Vector_Impl(b *testing.B) {
	c := vector.New(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_of_Collect_Immutable_Vector_Values(b *testing.B) {
	c := vector.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range c.Collect() {
			result = v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_of_Collect_Immutable_Vector_Impl_Values(b *testing.B) {
	c := vector.New(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range c.Collect() {
			result = v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Immutable_OrdererSet(b *testing.B) {
	c := set.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Immutable_OrdererSet_Impl(b *testing.B) {
	c := set.ToOrderedSet(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_of_Collect_Immutable_OrdererSet_Values(b *testing.B) {
	c := set.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range c.Collect() {
			result = v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_of_Collect_Immutable_OrdererSet_Impl_Values(b *testing.B) {
	c := set.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range c.Collect() {
			result = v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Mutable_OrdererSet(b *testing.B) {
	set := mset.ToOrderedSet(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = set.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Mutable_OrdererSet_Impl(b *testing.B) {
	set := mset.ToOrderedSet(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = set.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_Immutable_Vector(b *testing.B) {
	c := vector.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it := c.Begin(); it.HasNext(); {
			result, _ = it.Get()
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_Immutable_Vector_Impl(b *testing.B) {
	c := vector.Convert(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it := c.Iter(); it.HasNext(); {
			result, _ = it.Get()
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_WrapSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it := it.Wrap(values); it.HasNext(); {
			result, _ = it.Get()
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_Impl(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it := impliter.New(values); it.HasNext(); {
			result, _ = it.Get()
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_Impl_Reseteable(b *testing.B) {
	it := impliter.NewReseteable(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it.HasNext() {
			result, _ = it.Get()
		}
		it.Reset()
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_EmbeddedSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			result = v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForByIndex_EmbeddedSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(values); j++ {
			v := values[j]
			result = v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_WrapMap_HasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := it.WrapMap(values)
		for iter.HasNext() {
			kv, _ := iter.Get()
			_ = kv.Key()
			_ = kv.Value()
		}
	}
	b.StopTimer()
}

func Benchmark_NewKVHasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for iter := impliter.NewKV(values); iter.HasNext(); {
			k, v, _ := iter.GetKV()
			_ = k
			_ = v
		}
	}
	b.StopTimer()
}

func Benchmark_NewReflectKVHasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for iter := impliter.NewReflectKV(values); iter.HasNext(); {
			k, v, _ := iter.GetKV()
			_ = k
			_ = v
		}
	}
	b.StopTimer()
}

func Benchmark_EmbeddedMap_For(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k, v := range values {
			_ = k
			_ = v
		}
	}
	b.StopTimer()
}
