package benchmark

import (
	"testing"

	"github.com/m4gshm/container/immutable/set"
	"github.com/m4gshm/container/immutable/vector"
	"github.com/m4gshm/container/iter"
	impliter "github.com/m4gshm/container/iter/impl/iter"
	mset "github.com/m4gshm/container/mutable/set"
	"github.com/m4gshm/container/slice"
)

var (
	max    = 100000
	values = slice.Range(1, max)
	result = 0
)

func Benchmark_ForEach_Immutable_Vector(b *testing.B) {
	c := vector.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Immutable_Vector_Impl(b *testing.B) {
	c := vector.New(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_Immutable_Vector_Values(b *testing.B) {
	c := vector.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range c.Values() {
			result = v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_Immutable_Vector_Impl_Values(b *testing.B) {
	c := vector.New(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range c.Values() {
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
		c.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Immutable_OrdererSet_Impl(b *testing.B) {
	c := set.NewOrderedSet(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_Immutable_OrdererSet_Values(b *testing.B) {
	c := set.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range c.Values() {
			result = v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Immutable_OrdererSet_Impl_Values(b *testing.B) {
	c := set.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range c.Values() {
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
		set.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Mutable_OrdererSet_Impl(b *testing.B) {
	set := mset.ToOrderedSet(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_Immutable_Vector(b *testing.B) {
	c := vector.Of(values...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it := c.Begin(); it.HasNext(); {
			result = it.Get()
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
			result = it.Get()
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_WrapSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it := iter.Wrap(values); it.HasNext(); {
			result = it.Get()
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_Impl(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it := impliter.New(values); it.HasNext(); {
			result = it.Get()
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
			result = it.Get()
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
		iter := iter.WrapMap(values)
		for iter.HasNext() {
			kv := iter.Get()
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
			kv := iter.Get()
			_ = kv.Key()
			_ = kv.Value()
		}
	}
	b.StopTimer()
}

func Benchmark_NewMap_HasNextGet_Reset(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < max; i++ {
		values[i] = i
	}
	iter := impliter.NewKV(values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for iter.HasNext() {
			kv := iter.Get()
			_ = kv.Key()
			_ = kv.Value()
		}
		iter.Reset()
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
