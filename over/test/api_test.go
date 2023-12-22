package test

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/over"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/stretchr/testify/assert"
)

func Test_AllFiltered(t *testing.T) {
	from := set.Of(1, 2, 3, 5, 7, 8, 9, 11)

	s := []int{}

	for e := range over.Filtered(from.All, func(e int) bool { return e%2 == 0 }) {
		s = append(s, e)
	}

	assert.Equal(t, slice.Of(2, 8), sort.Asc(s))
}

func Test_AllConverted(t *testing.T) {
	from := set.Of(1, 2, 3, 5, 7, 8, 9, 11)
	s := []string{}

	for e := range over.Converted(from.All, strconv.Itoa) {
		s = append(s, e)
	}

	assert.Equal(t, slice.Of("1", "2", "3", "5", "7", "8", "9", "11"), s)
}

var (
	max        = 100000
	values     = range_.Closed(1, max)
)

func Benchmark_OrderedSet_Filter_Convert(b *testing.B) {
	c := set.Of(values...)

	var s []string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = collection.Convert(c.Filter(func(i int) bool { return i%2 == 0 }), strconv.Itoa).Slice()
	}
	b.StopTimer()

	_ = s
}

func Benchmark_OrderedSet_Filter_Convert_go1_22(b *testing.B) {
	c := set.Of(values...)

	var s []string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := range over.Converted(over.Filtered(c.All, func(e int) bool { return e%2 == 0 }), strconv.Itoa) {
			s = append(s, e)
		}
	}
	b.StopTimer()

	_ = s
}

func Benchmark_Slice_Filter_Convert(b *testing.B) {
	c := slice.Of(values...)

	var s []string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = slice.Convert(slice.Filter(c, func(i int) bool { return i%2 == 0 }), strconv.Itoa)
	}
	b.StopTimer()

	_ = s
}

func Benchmark_Slice_Filter_Convert_plainOld(b *testing.B) {
	c := values

	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range c {
			if e%2 == 0 {
				s = append(s, strconv.Itoa(e))
			}
		}
	}
	b.StopTimer()

	_ = s
}

func Benchmark_Slice_Filter_Convert_plainOld2(b *testing.B) {
	c := values

	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := make([]int, 0, len(s)/2)
		for _, e := range c {
			if e%2 == 0 {
				f = append(f, e)
			}
		}
		s = make([]string, len(f))
		for i := range f {
			s[i] = strconv.Itoa(f[i])
		}
	}
	b.StopTimer()

	_ = s
}
