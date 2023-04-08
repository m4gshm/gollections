package examples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/conv"
	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/last"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/exclude"
	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/predicate/one"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/m4gshm/gollections/slice/group"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/slice/reverse"
	"github.com/m4gshm/gollections/sum"
)

func Test_SortInt(t *testing.T) {
	source := []int{1, 3, -1, 2, 0}
	sorted := sort.Of(source)
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, sorted)
}

func Test_SortStructs(t *testing.T) {
	type User struct {
		name string
		age  int
	}
	var users = []User{{"Bob", 26}, {"Alice", 35}, {"Tom", 18}}
	var (
		//sorted
		byName = sort.By(users, func(u User) string { return u.name })
		byAge  = sort.By(users, func(u User) int { return u.age })
	)
	assert.Equal(t, []User{{"Alice", 35}, {"Bob", 26}, {"Tom", 18}}, byName)
	assert.Equal(t, []User{{"Tom", 18}, {"Bob", 26}, {"Alice", 35}}, byAge)
}

func Test_SortStructsByLess(t *testing.T) {
	type User struct {
		name string
		age  int
	}
	var users = []User{{"Bob", 26}, {"Alice", 35}, {"Tom", 18}}
	var (
		//sorted
		byName       = sort.ByLess(users, func(u1, u2 User) bool { return u1.name < u2.name })
		byAgeReverse = sort.ByLess(users, func(u1, u2 User) bool { return u1.age > u2.age })
	)
	assert.Equal(t, []User{{"Alice", 35}, {"Bob", 26}, {"Tom", 18}}, byName)
	assert.Equal(t, []User{{"Alice", 35}, {"Bob", 26}, {"Tom", 18}}, byAgeReverse)
}

func Test_Reverse(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, slice.Reverse([]int{3, 2, 1, 0, -1}))
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, reverse.Of([]int{3, 2, 1, 0, -1}))
}

func Test_Clone(t *testing.T) {
	type entity struct{ val string }
	var (
		entities = []*entity{{"first"}, {"second"}, {"third"}}
		copy     = clone.Of(entities)
	)

	assert.Equal(t, entities, copy)
	assert.NotSame(t, entities, copy)

	for i := range entities {
		assert.Same(t, entities[i], copy[i])
	}
}

func Test_DeepClone(t *testing.T) {
	type entity struct{ val string }
	var (
		entities = []*entity{{"first"}, {"second"}, {"third"}}
		copy     = clone.Deep(entities, clone.Ptr[entity])
	)

	assert.Equal(t, entities, copy)
	assert.NotSame(t, entities, copy)

	for i := range entities {
		assert.Equal(t, entities[i], copy[i])
		assert.NotSame(t, entities[i], copy[i])
	}
}

func Test_Convert(t *testing.T) {
	var (
		source   = slice.Of(1, 3, 5, 7, 9, 11)
		result   = slice.Convert(source, strconv.Itoa)
		expected = slice.Of("1", "3", "5", "7", "9", "11")
	)
	assert.Equal(t, expected, result)
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	var (
		source   = []int{1, 3, 4, 5, 7, 8, 9, 11}
		result   = slice.ConvertFit(source, even, strconv.Itoa)
		expected = []string{"4", "8"}
	)
	assert.Equal(t, expected, result)
}

func Test_ConvertFilteredWithIndexInPlace(t *testing.T) {
	var (
		source   = slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
		result   = slice.ConvertCheckIndexed(source, func(index int, elem int) (string, bool) { return strconv.Itoa(index + elem), even(elem) })
		expected = []string{"6", "13"}
	)
	assert.Equal(t, expected, result)
}

func Test_Slice_Filter(t *testing.T) {
	var (
		source   = []int{1, 2, 3, 4, 5, 6}
		result   = slice.Filter(source, even)
		expected = []int{2, 4, 6}
	)
	assert.Equal(t, expected, result)
}

func Test_Flatt(t *testing.T) {
	var (
		source   = [][]int{{1, 2, 3}, {4}, {5, 6}}
		result   = slice.Flatt(source, conv.AsIs[[]int])
		expected = []int{1, 2, 3, 4, 5, 6}
	)
	assert.Equal(t, expected, result)
}

func Test_Slice_Group(t *testing.T) {
	var (
		source   = []int{1, 2, 3, 4, 5, 6}
		result   = group.Of(source, even)
		expected = map[bool][]int{false: {1, 3, 5}, true: {2, 4, 6}}
	)
	assert.Equal(t, expected, result)
}

func Test_Slice_ReduceSum(t *testing.T) {
	var (
		source   = []int{1, 2, 3, 4, 5, 6}
		sum      = slice.Reduce(source, op.Sum[int])
		expected = 1 + 2 + 3 + 4 + 5 + 6
	)
	assert.Equal(t, expected, sum)
}

func Test_Slice_Sum(t *testing.T) {
	var (
		sum      = sum.Of(1, 2, 3, 4, 5, 6)
		expected = 1 + 2 + 3 + 4 + 5 + 6
	)
	assert.Equal(t, expected, sum)
}

func Test_Slice_Flatt(t *testing.T) {
	var (
		source   = [][]int{{1, 2, 3}, {4}, {5, 6}}
		result   = slice.Flatt(source, conv.AsIs[[]int])
		expected = []int{1, 2, 3, 4, 5, 6}
	)
	assert.Equal(t, expected, result)
}

func Test_Range(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, range_.Of(-1, 3))
	assert.Equal(t, []int{3, 2, 1, 0, -1}, range_.Of(3, -1))
	assert.Equal(t, []int{1}, range_.Of(1, 1))
}

func Test_First(t *testing.T) {
	result, ok := first.Of(1, 3, 5, 7, 9, 11).By(more.Than(5))
	assert.True(t, ok)
	assert.Equal(t, 7, result)
}

func Test_Last(t *testing.T) {
	result, ok := last.Of(1, 3, 5, 7, 9, 11).By(less.Than(9))
	assert.True(t, ok)
	assert.Equal(t, 7, result)
}

func Test_OneOf(t *testing.T) {
	result := slice.Filter(slice.Of(1, 3, 5, 7, 9, 11), one.Of(1, 7).Or(one.Of(11)))
	assert.Equal(t, slice.Of(1, 7, 11), result)
}

func Test_ExcludeAll(t *testing.T) {
	result := slice.Filter(slice.Of(1, 3, 5, 7, 9, 11), exclude.All(1, 7, 11))
	assert.Equal(t, slice.Of(3, 5, 9), result)
}

func Test_Xor(t *testing.T) {
	result := slice.Filter(slice.Of(1, 3, 5, 7, 9, 11), one.Of(1, 3, 5, 7).Xor(one.Of(7, 9, 11)))
	assert.Equal(t, slice.Of(1, 3, 5, 9, 11), result)
}

func Test_BehaveAsStrings(t *testing.T) {
	type (
		TypeBasedOnString      string
		ArrayTypeBasedOnString []TypeBasedOnString
	)

	var (
		source   = ArrayTypeBasedOnString{"1", "2", "3"}
		result   = slice.BehaveAsStrings(source)
		expected = []string{"1", "2", "3"}
	)

	assert.Equal(t, expected, result)
}

type rows[T any] struct {
	row    []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.row) }
func (r *rows[T]) next() (T, error) { e := r.row[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	var (
		stream    = &rows[int]{slice.Of(1, 2, 3), 0}
		result, _ = slice.OfLoop(stream, (*rows[int]).hasNext, (*rows[int]).next)
		expected  = slice.Of(1, 2, 3)
	)
	assert.Equal(t, expected, result)
}

func Test_Generate(t *testing.T) {
	var (
		counter   = 0
		result, _ = slice.Generate(func() (int, bool, error) { counter++; return counter, counter < 4, nil })
		expected  = slice.Of(1, 2, 3)
	)
	assert.Equal(t, expected, result)
}
