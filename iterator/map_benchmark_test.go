package iterator

import (
	"runtime"
	"testing"
)

func Benchmark_MapIterator(b *testing.B) {
	runtime.GOMAXPROCS(1)

	values := map[int]int{}
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[int]int, len(values))
		it := WrapMap(values)
		for it.Next() {
			kv := it.Get()
			result[kv.Key()] = kv.Value()
		}
	}
	b.StopTimer()
}

func Benchmark_EmbeddedMapFor(b *testing.B) {
	runtime.GOMAXPROCS(1)

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
