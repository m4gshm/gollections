//go:build goexperiment.rangefunc

package over_vs_loop

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/range_"
	"github.com/m4gshm/gollections/seq"
)

var max = 100000

var resultStr = ""

func Benchmark_loop_Converted(b *testing.B) {
	integers := range_.Of(0, max)
	for i := 0; i < b.N; i++ {
		loop.Convert(integers, strconv.Itoa).ForEach(func(element string) {
			resultStr = element
		})
	}
}

func Benchmark_loop_Converted_All(b *testing.B) {
	integers := range_.Of(0, max)
	for i := 0; i < b.N; i++ {
		for element := range loop.Convert(integers, strconv.Itoa).All {
			resultStr = element
		}
	}
}

func Benchmark_over_Converted(b *testing.B) {
	integers := range_.Of(0, max)
	for i := 0; i < b.N; i++ {
		for element := range seq.Convert(integers.All, strconv.Itoa) {
			resultStr = element
		}
	}
}

func Benchmark_over_Converted_direct(b *testing.B) {
	integers := range_.Of(0, max)
	for i := 0; i < b.N; i++ {
		seq.Convert(integers.All, strconv.Itoa)(func(element string) bool {
			resultStr = element
			return true
		})
	}
}

func even(i int) bool {
	return i%2 == 0
}

func Benchmark_loop_Convert_Filtered(b *testing.B) {
	integers := range_.Of(0, max)
	for i := 0; i < b.N; i++ {
		loop.Convert(loop.Filter(integers, even), strconv.Itoa).ForEach(func(element string) {
			resultStr = element
		})
	}
}

func Benchmark_loop_Convert_Filtered_rangefunc(b *testing.B) {
	integers := range_.Of(0, max)
	for i := 0; i < b.N; i++ {
		for element := range loop.Convert(loop.Filter(integers, even), strconv.Itoa).All {
			resultStr = element
		}
	}
}

func Benchmark_over_Convert_Filtered(b *testing.B) {
	integers := range_.Of(0, max)
	for i := 0; i < b.N; i++ {
		for element := range seq.Convert(seq.Filter(integers.All, even), strconv.Itoa) {
			resultStr = element
		}
	}
}
