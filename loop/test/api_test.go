package test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/conv"
	"github.com/m4gshm/gollections/loop/convert"
	"github.com/m4gshm/gollections/loop/first"
	"github.com/m4gshm/gollections/loop/range_"
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

func Test_NotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = loop.NotNil(source)
		expected = []*entity{{"first"}, {"third"}, {"fifth"}}
	)
	assert.Equal(t, expected, loop.Slice(result.Next))
}

func Test_ConvertPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = loop.ToValues(source)
		expected = []entity{{"first"}, {}, {"third"}, {}, {"fifth"}}
	)
	assert.Equal(t, expected, loop.Slice(result.Next))
}

func Test_ConvertNotnilPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = loop.GetValues(source)
		expected = []entity{{"first"}, {"third"}, {"fifth"}}
	)
	assert.Equal(t, expected, loop.Slice(result.Next))
}

func Test_Convert(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r := loop.Convert(s, strconv.Itoa)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, loop.Slice(r.Next))
}

func Test_ConvertNotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = convert.NotNil(source, func(e *entity) string { return e.val })
		expected = []string{"first", "third", "fifth"}
	)
	assert.Equal(t, expected, loop.Slice(result.Next))
}

func Test_ConvertToNotNil(t *testing.T) {
	type entity struct{ val *string }
	var (
		first    = "first"
		third    = "third"
		fifth    = "fifth"
		source   = loop.Of([]entity{{&first}, {}, {&third}, {}, {&fifth}}...)
		result   = convert.ToNotNil(source, func(e entity) *string { return e.val })
		expected = []*string{&first, &third, &fifth}
	)
	assert.Equal(t, expected, loop.Slice(result.Next))
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
	assert.Equal(t, expected, loop.Slice(result.Next))
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.FilterAndConvert(s, even, strconv.Itoa)
	assert.Equal(t, []string{"4", "8"}, loop.Slice(r.Next))
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.ConvertCheck(s, func(i int) (string, bool) { return strconv.Itoa(i), even(i) })
	assert.Equal(t, []string{"4", "8"}, loop.Slice(r.Next))
}

func Test_Flatt(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.Flatt(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, loop.Slice(f.Next))
}

func Test_FlattFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FilterAndFlatt(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i })
	e := []int{1, 2, 3, 5, 6}
	assert.Equal(t, e, loop.Slice(f.Next))
}

func Test_FlattElemFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FlattAndFilter(md, func(i []int) []int { return i }, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, loop.Slice(f.Next))
}

func Test_FilterAndFlattFit(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FilterFlattFilter(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i }, even)
	e := []int{2, 6}
	assert.Equal(t, e, loop.Slice(f.Next))
}

func Test_Filter(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.Filter(s, even)
	assert.Equal(t, slice.Of(4, 8), loop.Slice(r.Next))
}

func Test_Filtering(t *testing.T) {
	r := loop.Filter(loop.Of(1, 2, 3, 4, 5, 6), func(i int) bool { return i%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, loop.Slice(r.Next))
}

type rows[T any] struct {
	row    []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.row) }
func (r *rows[T]) next() (T, error) { e := r.row[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	stream := loop.Of(1, 2, 3)
	result := loop.Slice(stream)

	assert.Equal(t, slice.Of(1, 2, 3), result)
}

func Test_MatchAny(t *testing.T) {
	elements := loop.Of(1, 2, 3, 4)

	ok := loop.HasAny(elements, eq.To(4))
	assert.True(t, ok)

	noOk := loop.HasAny(elements, more.Than(5))
	assert.False(t, noOk)
}

func Test_MultipleKeyValuer(t *testing.T) {
	type Role struct {
		name string
	}

	type User struct {
		name  string
		age   int
		roles []Role
	}

	var users = []User{
		{name: "Bob", age: 26, roles: []Role{{"Admin"}, {"manager"}}},
		{name: "Alice", age: 35, roles: []Role{{"Manager"}}},
		{name: "Tom", age: 18}, {},
	}

	m := kvloop.Group(loop.ToKVs(loop.Of(users...),
		func(u User) []string {
			return slice.Convert(u.roles, func(r Role) string { return strings.ToLower(r.name) })
		},
		func(u User) []string { return []string{u.name, strings.ToLower(u.name)} },
	).Next)

	assert.Equal(t, m["admin"], slice.Of("Bob", "bob"))
	assert.Equal(t, m["manager"], slice.Of("Bob", "bob", "Alice", "alice"))
	assert.Equal(t, m[""], slice.Of("Tom", "tom", "", ""))
}

func Test_Range(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), loop.Slice(range_.Of(-1, 4)))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), loop.Slice(range_.Of(3, -2)))
	assert.Nil(t, loop.Slice(range_.Of(1, 1)))
}

func Test_RangeClosed(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), loop.Slice(range_.Closed(-1, 3)))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), loop.Slice(range_.Closed(3, -1)))
	assert.Equal(t, slice.Of(1), loop.Slice(range_.Closed(1, 1)))
}

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")
	result := loop.Slice(loop.OfIndexed(len(indexed), func(i int) string { return indexed[i] }))
	assert.Equal(t, indexed, result)
}

func Test_ConvertIndexed(t *testing.T) {
	indexed := slice.Of(10, 11, 12, 13, 14)
	result := loop.Slice(convert.FromIndexed(len(indexed), func(i int) int { return indexed[i] }, strconv.Itoa).Next)
	assert.Equal(t, slice.Of("10", "11", "12", "13", "14"), result)
}

func Test_ConvIndexed(t *testing.T) {
	indexed := slice.Of("10", "11", "12", "13", "14")
	result, err := breakLoop.Slice(conv.FromIndexed(len(indexed), func(i int) string { return indexed[i] }, strconv.Atoi).Next)
	assert.NoError(t, err)
	assert.Equal(t, slice.Of(10, 11, 12, 13, 14), result)
}
