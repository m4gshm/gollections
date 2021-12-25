package benchmark

import (
	"testing"

	"github.com/m4gshm/container/iter"
)

var amount = 100_000

func Benchmark_ForEach_Iterator_Interface(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		iter.ForEach(iter.Wrap(values), func(v int) { result = append(result, v) })
	}
	b.StopTimer()
}

func Benchmark_ForEach_Iterator_Impl(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		iter.ForEach(iter.Wrap(values), func(v int) { result = append(result, v) })
	}
	b.StopTimer()
}

func Benchmark_IterateHasNextGet_Iterator_Interface(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := iter.Wrap(values)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_IterateHasNextGet_Iterator_Impl(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := iter.Wrap(values)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_WrapMapHasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[int]int, len(values))
		iter := iter.WrapMap(values)
		for iter.HasNext() {
			kv := iter.Get()
			result[kv.Key()] = kv.Value()
		}
	}
	b.StopTimer()
}

func Benchmark_NewMapHasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[int]int, len(values))
		iter := iter.NewMap(values)
		for iter.HasNext() {
			kv := iter.Get()
			result[kv.Key()] = kv.Value()
		}
	}
	b.StopTimer()
}

func Benchmark_EmbeddedMapFor(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[int]int, len(values))

		for k, v := range values {
			result[k] = v
		}
	}
	b.StopTimer()
}
