package iterator

import (
	"testing"
)

func Benchmark_WrapMapHasNextGet(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[int]int, len(values))
		it := WrapMap(values)
		for it.HasNext() {
			kv := it.Get()
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
		it := NewMap(values)
		for it.HasNext() {
			kv := it.Get()
			result[kv.Key()] = kv.Value()
		}
	}
	b.StopTimer()
}

func Benchmark_WrapMapNextNext(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[int]int, len(values))
		it := WrapMap(values)
		for kv, ok := it.Next(); ok; kv, ok = it.Next() {
			result[kv.Key()] = kv.Value()
		}
	}
	b.StopTimer()
}

func Benchmark_NewMapNextNext(b *testing.B) {
	values := map[int]int{}
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[int]int, len(values))
		it := NewMap(values)
		for kv, ok := it.Next(); ok; kv, ok = it.Next() {
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
