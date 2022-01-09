package vector

import (
	"testing"

	"github.com/m4gshm/container/mutable/vector"
	"github.com/m4gshm/container/slice"
)

var (
	max    = 100000
	values = slice.Range(1, max)
)

func Benchmark_Vector_Add(b *testing.B) {
	v := vector.New[int](max)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v.Add(values...)
	}
	b.StopTimer()
	_ = v
}

func Benchmark_VectorImpl_Add(b *testing.B) {
	v := vector.Create[int](max)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v.Add(values...)
	}
	b.StopTimer()
	_ = v
}

func Benchmark_VectorImpl_Add_ByOne(b *testing.B) {
	v := vector.Create[int](max)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, i := range values {
			_, _ = v.Add(i)
		}
	}
	b.StopTimer()
	_ = v
}

func Benchmark_VectorImpl_Add_One(b *testing.B) {
	v := vector.Create[int](max)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, i := range values {
			_, _ = v.AddOne(i)
		}
	}
	b.StopTimer()
	_ = v
}

func Benchmark_VectorImpl_Add_All(b *testing.B) {
	v := vector.Create[int](max)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v.AddAll(values)
	}
	b.StopTimer()
	_ = v
}