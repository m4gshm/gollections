package examples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/last"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/ordered"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/slice/reverse"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/m4gshm/gollections/sum"
)

func Test_SortInt(t *testing.T) {
	c := ordered.Sort([]int{-1, 0, 1, 2, 3})
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, c)
}

func Test_SortStructs(t *testing.T) {
	type User struct {
		name string
		age  int
	}
	var (
		users  = []User{{"Bob", 26}, {"Alice", 35}, {"Tom", 18}}
		byName = sort.ByOrdered(clone.Of(users), func(u User) string { return u.name })
		byAge  = sort.ByOrdered(clone.Of(users), func(u User) int { return u.age })
	)
	assert.Equal(t, []User{{"Alice", 35}, {"Bob", 26}, {"Tom", 18}}, byName)
	assert.Equal(t, []User{{"Tom", 18}, {"Bob", 26}, {"Alice", 35}}, byAge)
}

func Test_Reverse(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, slice.Reverse([]int{3, 2, 1, 0, -1}))
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, reverse.Of([]int{3, 2, 1, 0, -1}))
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

func Test_Slice_Filter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	f := slice.Filter(s, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, f)
}

func Test_Flatt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.Flatt(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_Slice_Group(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	g := slice.Group(s, even)
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
	f := slice.Flatt(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_Range(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, range_.Of(-1, 3))
	assert.Equal(t, []int{3, 2, 1, 0, -1}, range_.Of(3, -1))
	assert.Equal(t, []int{1}, range_.Of(1, 1))
}

func Test_First(t *testing.T) {
	r, ok := first.Of(1, 3, 5, 7, 9, 11).By(func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)
}

func Test_Last(t *testing.T) {
	r, ok := last.Of(1, 3, 5, 7, 9, 11).By(func(i int) bool { return i < 9 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)
}
