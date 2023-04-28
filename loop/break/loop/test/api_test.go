package test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/as"
	kvBreakLoop "github.com/m4gshm/gollections/kv/loop/break/loop"
	"github.com/m4gshm/gollections/loop"
	breakLoop "github.com/m4gshm/gollections/loop/break/loop"
	"github.com/m4gshm/gollections/op/break/op"
	"github.com/m4gshm/gollections/predicate/break/predicate/eq"
	"github.com/m4gshm/gollections/predicate/break/predicate/more"
	"github.com/m4gshm/gollections/slice"
)

func Test_ReduceSum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, _ := breakLoop.Reduce(breakLoop.From(s), op.Sum[int])
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_Sum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, _ := breakLoop.Sum(breakLoop.From(s))
	assert.Equal(t, 1+3+5+7+9+11, r)
}

// func Test_First(t *testing.T) {
// 	s := loop.Of(1, 3, 5, 7, 9, 11)
// 	r, ok := first.Of(s, func(i int) bool { return i > 5 })
// 	assert.True(t, ok)
// 	assert.Equal(t, 7, r)

// 	_, nook := loop.First(s, func(i int) bool { return i > 12 })
// 	assert.False(t, nook)
// }

func Test_Convert(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r := breakLoop.Convert(breakLoop.From(s), strconv.Itoa)
	o, _ := breakLoop.ToSlice(r.Next)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, o)
}

// func Test_ConvertNotNil(t *testing.T) {
// 	type entity struct{ val string }
// 	var (
// 		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
// 		result   = convert.NotNil(source, func(e *entity) string { return e.val })
// 		expected = []string{"first", "third", "fifth"}
// 	)
// 	assert.Equal(t, expected, slice.Generate(result.Next))
// }

// func Test_ConvertToNotNil(t *testing.T) {
// 	type entity struct{ val *string }
// 	var (
// 		first    = "first"
// 		third    = "third"
// 		fifth    = "fifth"
// 		source   = loop.Of([]entity{{&first}, {}, {&third}, {}, {&fifth}}...)
// 		result   = convert.ToNotNil(source, func(e entity) *string { return e.val })
// 		expected = []*string{&first, &third, &fifth}
// 	)
// 	assert.Equal(t, expected, slice.Generate(result.Next))
// }

// func Test_ConvertNilSafe(t *testing.T) {
// 	type entity struct{ val *string }
// 	var (
// 		first    = "first"
// 		third    = "third"
// 		fifth    = "fifth"
// 		source   = loop.Of([]*entity{{&first}, {}, {&third}, nil, {&fifth}}...)
// 		result   = convert.NilSafe(source, func(e *entity) *string { return e.val })
// 		expected = []*string{&first, &third, &fifth}
// 	)
// 	assert.Equal(t, expected, slice.Generate(result.Next))
// }

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := breakLoop.FilterAndConvert(breakLoop.From(s), even, strconv.Itoa)
	o, _ := breakLoop.ToSlice(r.Next)
	assert.Equal(t, []string{"4", "8"}, o)
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := breakLoop.ConvCheck(breakLoop.From(s), func(i int) (string, bool, error) { return strconv.Itoa(i), even(i), nil })
	o, _ := breakLoop.ToSlice(r.Next)
	assert.Equal(t, []string{"4", "8"}, o)
}

func Test_Flatt(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.Flatt(breakLoop.From(md), as.Is[[]int])
	e := []int{1, 2, 3, 4, 5, 6}
	o, _ := breakLoop.ToSlice(f.Next)
	assert.Equal(t, e, o)
}

func Test_FlattFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FilterAndFlatt(breakLoop.From(md), func(from []int) bool { return len(from) > 1 }, as.Is[[]int])
	e := []int{1, 2, 3, 5, 6}
	o, _ := breakLoop.ToSlice(f.Next)
	assert.Equal(t, e, o)
}

func Test_FlattElemFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FlattAndFilter(breakLoop.From(md), as.Is[[]int], even)
	e := []int{2, 4, 6}
	o, _ := breakLoop.ToSlice(f.Next)
	assert.Equal(t, e, o)
}

func Test_FilterAndFlattFit(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FilterFlattFilter(breakLoop.From(md), func(from []int) bool { return len(from) > 1 }, as.Is[[]int], even)
	e := []int{2, 6}
	o, _ := breakLoop.ToSlice(f.Next)
	assert.Equal(t, e, o)
}

func Test_Filter(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	f := breakLoop.Filter(breakLoop.From(s), even)
	e := []int{4, 8}
	o, _ := breakLoop.ToSlice(f.Next)
	assert.Equal(t, e, o)
}

func Test_Filtering(t *testing.T) {
	r := breakLoop.Filt(breakLoop.From(loop.Of(1, 2, 3, 4, 5, 6)), func(i int) (bool, error) { return i%2 == 0, nil })
	o, _ := breakLoop.ToSlice(r.Next)
	assert.Equal(t, []int{2, 4, 6}, o)
}

func Test_MatchAny(t *testing.T) {
	elements := loop.Of(1, 2, 3, 4)

	ok, _ := breakLoop.HasAny(breakLoop.From(elements), eq.To(4))
	assert.True(t, ok)

	noOk, _ := breakLoop.HasAny(breakLoop.From(elements), more.Than(5))
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

	m, _ := kvBreakLoop.Group(breakLoop.ToKVs(breakLoop.From(loop.Of(users...)),
		func(u User) ([]string, error) {
			return slice.Convert(u.roles, func(r Role) string { return strings.ToLower(r.name) }), nil
		},
		func(u User) ([]string, error) { return []string{u.name, strings.ToLower(u.name)}, nil },
	).Next)

	assert.Equal(t, m["admin"], slice.Of("Bob", "bob"))
	assert.Equal(t, m["manager"], slice.Of("Bob", "bob", "Alice", "alice"))
	assert.Equal(t, m[""], slice.Of("Tom", "tom", "", ""))

}
