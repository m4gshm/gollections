package range_

import (
	"testing"

	"github.com/m4gshm/gollections/loop"
	lrange "github.com/m4gshm/gollections/loop/range_"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
	srange "github.com/m4gshm/gollections/slice/range_"
)

var (
	maxVal = 100000
	values = srange.Closed(1, maxVal)
)

func Benchmark_Slice_RangeClosed_Generate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = srange.Closed(1, maxVal)
	}
}

func Benchmark_Slice_RangeClosed_Iterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i, v := range values {
			_, _ = i, v
		}
	}
}

func Benchmark_Seq_RangeClosed_Iterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for v := range seq.RangeClosed(1, maxVal) {
			_ = v
		}
	}
}

func Benchmark_Loop_RangeClosed_GenerateIterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		next := lrange.Closed(1, maxVal)
		for {
			n, ok := next()
			if !ok {
				break
			}
			_ = n
		}
	}
}

func Benchmark_Slice_Series_Generate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := range slice.Series(1, func(prev int) (int, bool) { return prev + 1, prev <= maxVal }) {
			_ = n
		}
	}
}

func Benchmark_Seq_Series_Generate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := range seq.Series(1, func(prev int) (int, bool) { return prev + 1, prev <= maxVal }) {
			_ = n
		}
	}
}

func Benchmark_Loop_Series_Generate_Iterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		next := loop.Series(1, func(prev int) (int, bool) { return prev + 1, prev <= maxVal })
		for {
			n, ok := next()
			if !ok {
				break
			}
			_ = n
		}
	}
}
