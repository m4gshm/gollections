package vector

import (
	"testing"

	"github.com/m4gshm/gollections/collection/mutable/vector"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	maxValue = 100000
	values   = range_.Closed(1, maxValue)
)

func Benchmark_Vector_Add(b *testing.B) {
	v := vector.NewCap[int](maxValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Add(values...)
	}
	b.StopTimer()
	_ = v
}

func Benchmark_Vector_Add_ByOne(b *testing.B) {
	v := vector.NewCap[int](maxValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, i := range values {
			v.Add(i)
		}
	}
	b.StopTimer()
	_ = v
}

func Benchmark_Vector_Add_All(b *testing.B) {
	v := vector.NewCap[int](maxValue)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Add(values...)
	}
	b.StopTimer()
	_ = v
}
