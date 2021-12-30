package benchmark

import (
	"testing"

	"github.com/m4gshm/container/immutable"
	"github.com/m4gshm/container/iter"
	impliter "github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/mutable"
)

var amount = 100_000

func Benchmark_ForEach_Pipe(b *testing.B) {
	amount = 2
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}

	var result int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter.Pipe(iter.Wrap(values)).ForEach(func(v int) { result = v })
	}

	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Iterator_Interface(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter.ForEach(iter.Wrap(values), func(v int) { result++ })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Iterator_Impl(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		impliter.ForEach(impliter.New(&values), func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Immutable_OrdererSet_Iterator(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	set := immutable.NewOrderedSet(values...)
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		impliter.ForEach(set.Begin(), func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Mutable_OrdererSet_Iterator(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	var set mutable.Set[int] = mutable.NewOrderedSet(values...)
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		impliter.ForEach(set.Begin(), func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGet_Iterator_Interface(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	result := 0
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
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for it := impliter.New(&values); it.HasNext(); {
			result = it.Get()
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_HasNextGetReset_Iterator_Impl(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	it := impliter.NewReseteable(&values)
	result := 0
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

func Benchmark_ForEach_Immutable_OrdererSet_Walker(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	set := immutable.NewOrderedSet(values...)
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Immutable_OrdererSet_Walker_Impl(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	set := immutable.NewSet(values)
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}
func Benchmark_ForEach_Mutable_OrdererSet_Walker(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	set := mutable.NewOrderedSet(values...)
	var result int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForEach_Mutable_OrdererSet_Walker_Impl(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	set := mutable.NewSet(values)
	var result int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.ForEach(func(v int) { result = v })
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForRange_EmbeddedSlice(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			result += v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_ForByIndex_EmbeddedSlice(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	result := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(values); j++ {
			v := values[j]
			result += v
		}
	}
	b.StopTimer()
	_ = result
}

func Benchmark_WrapMap_HasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < 100000; i++ {
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
	for i := 0; i < 100000; i++ {
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
	for i := 0; i < 100000; i++ {
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
	for i := 0; i < 100000; i++ {
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
