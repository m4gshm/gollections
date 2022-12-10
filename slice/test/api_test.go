package test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/first"
	"github.com/m4gshm/gollections/slice/last"
	"github.com/m4gshm/gollections/slice/range_"
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
	r := slice.Reduce(s, op.Sum[int])
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_Sum(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Sum(s)
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_First(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r, ok := first.Of(s, func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := slice.First(s, func(i int) bool { return i > 12 })
	assert.False(t, nook)
}

func Test_Last(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r, ok := last.Of(s, func(i int) bool { return i < 9 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := slice.Last(s, func(i int) bool { return i < 1 })
	assert.False(t, nook)
}

func Test_Convert(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Map(s, strconv.Itoa)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, r)
}

func Test_ConvertWithIndex(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.MapIndex(s, func(index int, elem int) int { return index + elem })
	assert.Equal(t, slice.Of(1, 1+3, 2+5, 3+7, 4+9, 5+11), r)
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.MapFit(s, even, strconv.Itoa)
	assert.Equal(t, []string{"4", "8"}, r)
}

func Test_ConvertFilteredWithIndex(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.MapFitIndex(s, func(_ int, elem int) bool { return even(elem) }, func(index int, elem int) string { return strconv.Itoa(index + elem) })
	assert.Equal(t, []string{"6", "13"}, r)
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.MapCheck(s, func(i int) (string, bool) { return strconv.Itoa(i), even(i) })
	assert.Equal(t, []string{"4", "8"}, r)
}

func Test_ConvertFilteredWithIndexInPlace(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.MapCheckIndex(s, func(index int, elem int) (string, bool) { return strconv.Itoa(index + elem), even(elem) })
	assert.Equal(t, []string{"6", "13"}, r)
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

func Test_BehaveAsStrings(t *testing.T) {
	type TypeBasedOnString string
	type ArrayTypeBasedOnString []TypeBasedOnString

	vals := ArrayTypeBasedOnString{"1", "2", "3"}
	strs := slice.BehaveAsStrings(vals)

	assert.Equal(t, []string{"1", "2", "3"}, strs)
	pvals := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vals)).Data)
	pstrs := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&strs)).Data)
	assert.Equal(t, pvals, pstrs)
}

func Test_StringsBehaveAs(t *testing.T) {
	type TypeBasedOnString string
	type ArrayTypeBasedOnString []TypeBasedOnString

	vals := []string{"1", "2", "3"}
	strs := slice.StringsBehaveAs[ArrayTypeBasedOnString](vals)

	assert.Equal(t, ArrayTypeBasedOnString{"1", "2", "3"}, strs)
	pvals := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vals)).Data)
	pstrs := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&strs)).Data)
	assert.Equal(t, pvals, pstrs)
}

func Test_StringsBehaveAs2(t *testing.T) {
	type TypeBasedOnString string

	vals := []string{"1", "2", "3"}
	strs := slice.StringsBehaveAs[[]TypeBasedOnString](vals)

	assert.Equal(t, []TypeBasedOnString{"1", "2", "3"}, strs)
	pvals := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vals)).Data)
	pstrs := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&strs)).Data)
	assert.Equal(t, pvals, pstrs)
}

type rows[T any] struct {
	in     []T
	cursor int
}

func (r *rows[T]) hasNext() bool {
	return r.cursor < len(r.in)
}

func (r *rows[T]) next() (T, error) {
	e := r.in[r.cursor]
	r.cursor++
	return e, nil
}

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := slice.OfLoop(stream, (*rows[int]).hasNext, (*rows[int]).next)

	assert.Equal(t, slice.Of(1, 2, 3), result)
}

func Test_Generate(t *testing.T) {
	counter := 0
	result, _ := slice.Generate(func() (int, bool, error) {
		counter++
		return counter, counter < 4, nil

	})

	assert.Equal(t, slice.Of(1, 2, 3), result)
}
