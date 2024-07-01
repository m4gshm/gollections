//go:build goexperiment.rangefunc

package test

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/seq"
)

func Benchmark_OrderedSet_Filter_Convert_go1_22(b *testing.B) {
	c := set.Of(values...)

	var s []string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := range seq.Convert(seq.Filter(c.All, func(e int) bool { return e%2 == 0 }), strconv.Itoa) {
			s = append(s, e)
		}
	}
	b.StopTimer()

	_ = s
}

func Benchmark_Seq_Filter_Convert_go1_22(b *testing.B) {
	var s []string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := range seq.Convert(seq.Filter(seq.Of(values...), func(e int) bool { return e%2 == 0 }), strconv.Itoa) {
			s = append(s, e)
		}
	}
	b.StopTimer()

	_ = s
}
