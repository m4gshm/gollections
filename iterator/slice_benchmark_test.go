package iterator

import (
	"testing"
)

func Benchmark_SliceIterator(b *testing.B) {
	values := make([]int, 100000)
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 100000)
		it := Wrap(values)
		for it.Next() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_EmbeddedSliceFor(b *testing.B) {
	values := make([]int, 100000)
	for i := 0; i < 100000; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 100000)
		for _, v := range values {
			result = append(result, v)
		}
	}
	b.StopTimer()
}
