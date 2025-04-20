package test

import (
	"testing"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	maxValue = 100000
	values   = range_.Closed(1, maxValue)
)

func Benchmark_Loop_Slice_Filter_plainOld(b *testing.B) {
	c := values

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range c {
			if even(e) {
				_ = e
			}
		}
	}
	b.StopTimer()
}

func Benchmark_Loop_Seq_Filter_Seq(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := range seq.Filter(seq.Of(values...), even) {
			_ = e
		}
	}
	b.StopTimer()
}

func Benchmark_Loop_Loop_Filter_Seq(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		next := loop.Filter(loop.Of(values...), even)
		for {
			e, ok := next()
			if !ok {
				break
			}
			_ = e
		}
	}
	b.StopTimer()
}
