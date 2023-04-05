package examples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/conv"
	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/last"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/m4gshm/gollections/slice/group"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/slice/reverse"
	"github.com/m4gshm/gollections/sum"
)

func Test_SortInt(t *testing.T) {
	ints := []int{1, 3, -1, 2, 0}
	sorted := sort.Of(ints)
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
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Convert(s, strconv.Itoa)
	assert.Equal(t, slice.Of("1", "3", "5", "7", "9", "11"), r)
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := []int{1, 3, 4, 5, 7, 8, 9, 11}
	r := slice.ConvertFit(s, even, strconv.Itoa)
	assert.Equal(t, []string{"4", "8"}, r)
}

func Test_ConvertFilteredWithIndexInPlace(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.ConvertCheckIndexed(s, func(index int, elem int) (string, bool) { return strconv.Itoa(index + elem), even(elem) })
	assert.Equal(t, []string{"6", "13"}, r)
}

func Test_Slice_Filter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	f := slice.Filter(s, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, f)
}

func Test_Flatt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.Flatt(md, conv.AsIs[[]int])
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_Slice_Group(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	g := group.Of(s, even)
	e := map[bool][]int{false: {1, 3, 5}, true: {2, 4, 6}}
	assert.Equal(t, e, g)
}

func Test_Slice_ReduceSum(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	sum := slice.Reduce(s, op.Sum[int])
	e := 1 + 2 + 3 + 4 + 5 + 6
	assert.Equal(t, e, sum)
}

func Test_Slice_Sum(t *testing.T) {
	sum := sum.Of(1, 2, 3, 4, 5, 6)
	e := 1 + 2 + 3 + 4 + 5 + 6
	assert.Equal(t, e, sum)
}

func Test_Slice_Flatt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.Flatt(md, conv.AsIs[[]int])
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_Range(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, range_.Of(-1, 3))
	assert.Equal(t, []int{3, 2, 1, 0, -1}, range_.Of(3, -1))
	assert.Equal(t, []int{1}, range_.Of(1, 1))
}

func Test_First(t *testing.T) {
	r, ok := first.Of(1, 3, 5, 7, 9, 11).By(more.Than(5))
	assert.True(t, ok)
	assert.Equal(t, 7, r)
}

func Test_Last(t *testing.T) {
	r, ok := last.Of(1, 3, 5, 7, 9, 11).By(less.Than(9))
	assert.True(t, ok)
	assert.Equal(t, 7, r)
}

func Test_BehaveAsStrings(t *testing.T) {
	type TypeBasedOnString string
	type ArrayTypeBasedOnString []TypeBasedOnString

	vals := ArrayTypeBasedOnString{"1", "2", "3"}
	strs := slice.BehaveAsStrings(vals)

	assert.Equal(t, []string{"1", "2", "3"}, strs)
}

type rows[T any] struct {
	row    []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.row) }
func (r *rows[T]) next() (T, error) { e := r.row[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := slice.OfLoop(stream, (*rows[int]).hasNext, (*rows[int]).next)

	assert.Equal(t, slice.Of(1, 2, 3), result)
}

func Test_Generate(t *testing.T) {
	counter := 0
	result, _ := slice.Generate(func() (int, bool, error) { counter++; return counter, counter < 4, nil })

	assert.Equal(t, slice.Of(1, 2, 3), result)
}
