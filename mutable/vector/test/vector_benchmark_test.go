package vector

import (
	"testing"

	"github.com/m4gshm/gollections/mutable/vector"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	max    = 100000
	values = range_.Of(1, max)
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

func Benchmark_Vector_Add_ByOne(b *testing.B) {
	v := vector.New[int](max)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, i := range values {
			_, _ = v.Add(i)
		}
	}
	b.StopTimer()
	_ = v
}

func Benchmark_Vector_Add_All(b *testing.B) {
	v := vector.New[int](max)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v.AddAll(values)
	}
	b.StopTimer()
	_ = v
}
