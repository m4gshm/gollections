package iterator

import (
	"testing"
)

var amount = 100_000

func Benchmark_WrapSliceHasNextGet(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := Wrap(values)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_WrapSliceNextNext(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := Wrap(values)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_NewSliceNextNext(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := New(values...)
		for v, ok := it.Next(); ok; v, ok = it.Next() {
			result = append(result, v)
		}
	}
	b.StopTimer()
}

func Benchmark_WrapSliceForEach(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		ForEach(Wrap(values), func(v int) { result = append(result, v) })
	}
	b.StopTimer()
}

func Benchmark_NewSliceHasNextGet(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := New(values...)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}
