package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/sum"
	"github.com/stretchr/testify/assert"
)

func Test_Range(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), range_.Of(-1, 3))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), range_.Of(3, -1))
	assert.Equal(t, slice.Of(1), range_.Of(1, 1))
}

func Test_Reverse(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), slice.Reverse(range_.Of(3, -1)))
}

func Test_ReduceSum(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Reduce(s, sum.Of[int])
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_Convert(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Map(s, strconv.Itoa)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, r)
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.MapFit(s, even, strconv.Itoa)
	assert.Equal(t, []string{"4", "8"}, r)
}

func Test_Flatt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.Flatt(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_FlattFilter(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.FlattFit(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i })
	e := []int{1, 2, 3, 5, 6}
	assert.Equal(t, e, f)
}

func Test_FlattElemFilter(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.FlattElemFit(md, func(i []int) []int { return i }, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, f)
}

func Test_FlattFitFit(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.FlattFitFit(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i }, even)
	e := []int{2, 6}
	assert.Equal(t, e, f)
}

func Test_Filter(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.Filter(s, even)
	assert.Equal(t, slice.Of(4, 8), r)
}

func Test_StringRepresentation(t *testing.T) {
	var (
		expected = fmt.Sprint(slice.Of(1, 2, 3, 4))
		actual   = slice.ToString(slice.Of(1, 2, 3, 4))
	)
	assert.Equal(t, expected, actual)
}

func Test_StringReferencesRepresentation(t *testing.T) {
	var (
		expected       = fmt.Sprint(slice.Of(1, 2, 3, 4))
		i1, i2, i3, i4 = 1, 2, 3, 4
		actual         = slice.ToStringRefs(slice.Of(&i1, &i2, &i3, &i4))
	)
	assert.Equal(t, expected, actual)
}

func Test_Filtering(t *testing.T) {
	r := slice.Filter([]int{1, 2, 3, 4, 5, 6}, func(i int) bool { return i%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, r)
}
