package range_

import (
	"testing"

	"github.com/m4gshm/gollections/loop"
	lrange "github.com/m4gshm/gollections/loop/range_"
	"github.com/m4gshm/gollections/slice"
	srange "github.com/m4gshm/gollections/slice/range_"
)

var (
	max    = 100000
	values = srange.Closed(1, max)
)

func Benchmark_Slice_RangeClosed_Generate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = srange.Closed(1, max)
	}
}

func Benchmark_Slice_RangeClosed_Iterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i, v := range values {
			_, _ = i, v
		}
	}
}

func Benchmark_Loop_RangeClosed_GenerateIterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		next := lrange.Closed(1, max)
		for {
			n, ok := next()
			if !ok {
				break
			}
			_ = n
		}
	}
}

func Benchmark_Slice_Sequence_Generate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = slice.Sequence(1, func(prev int) (int, bool) { return prev + 1, prev <= max })
	}
}

func Benchmark_Loop_Sequence_Generate_Iterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		next := loop.Sequence(1, func(prev int) (int, bool) { return prev + 1, prev <= max })
		for {
			n, ok := next()
			if !ok {
				break
			}
			_ = n
		}
	}
}
