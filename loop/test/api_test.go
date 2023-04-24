package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/convert"
	"github.com/m4gshm/gollections/loop/first"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/slice"
)

func Test_ReduceSum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r := loop.Reduce(s, op.Sum[int])
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_Sum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r := loop.Sum(s)
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_First(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, ok := first.Of(s, func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := loop.First(s, func(i int) bool { return i > 12 })
	assert.False(t, nook)
}

func Test_Convert(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r := loop.Convert(s, strconv.Itoa)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, slice.Generate(r.Next))
}

func Test_ConvertNotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = convert.NotNil(source, func(e *entity) string { return e.val })
		expected = []string{"first", "third", "fifth"}
	)
	assert.Equal(t, expected, slice.Generate(result.Next))
}

func Test_ConvertToNotNil(t *testing.T) {
	type entity struct{ val *string }
	var (
		first    = "first"
		third    = "third"
		fifth    = "fifth"
		source   = loop.Of([]entity{{&first}, {}, {&third}, {}, {&fifth}}...)
		result   =convert.ToNotNil(source, func(e entity) *string { return e.val })
		expected = []*string{&first, &third, &fifth}
	)
	assert.Equal(t, expected, slice.Generate(result.Next))
}

func Test_ConvertNilSafe(t *testing.T) {
	type entity struct{ val *string }
	var (
		first    = "first"
		third    = "third"
		fifth    = "fifth"
		source   = loop.Of([]*entity{{&first}, {}, {&third}, nil, {&fifth}}...)
		result   = convert.NilSafe(source, func(e *entity) *string { return e.val })
		expected = []*string{&first, &third, &fifth}
	)
	assert.Equal(t, expected, slice.Generate(result.Next))
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.FilterAndConvert(s, even, strconv.Itoa)
	assert.Equal(t, []string{"4", "8"}, slice.Generate(r.Next))
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.ConvertCheck(s, func(i int) (string, bool) { return strconv.Itoa(i), even(i) })
	assert.Equal(t, []string{"4", "8"}, slice.Generate(r.Next))
}

func Test_Flatt(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.Flatt(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, slice.Generate(f.Next))
}

func Test_FlattFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FilterAndFlatt(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i })
	e := []int{1, 2, 3, 5, 6}
	assert.Equal(t, e, slice.Generate(f.Next))
}

func Test_FlattElemFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FlattAndFilter(md, func(i []int) []int { return i }, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, slice.Generate(f.Next))
}

func Test_FilterAndFlattFit(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FilterFlattFilter(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i }, even)
	e := []int{2, 6}
	assert.Equal(t, e, slice.Generate(f.Next))
}

func Test_Filter(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.Filter(s, even)
	assert.Equal(t, slice.Of(4, 8), slice.Generate(r.Next))
}

func Test_Filtering(t *testing.T) {
	r := loop.Filter(loop.Of(1, 2, 3, 4, 5, 6), func(i int) bool { return i%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, slice.Generate(r.Next))
}

type rows[T any] struct {
	row    []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.row) }
func (r *rows[T]) next() (T, error) { e := r.row[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	stream := loop.Of(1, 2, 3)
	result := slice.Generate(stream)

	assert.Equal(t, slice.Of(1, 2, 3), result)
}

func Test_MatchAny(t *testing.T) {
	elements := loop.Of(1, 2, 3, 4)

	ok := loop.HasAny(elements, eq.To(4))
	assert.True(t, ok)

	noOk := loop.HasAny(elements, more.Than(5))
	assert.False(t, noOk)
}
