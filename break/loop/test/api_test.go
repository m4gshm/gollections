package test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	breakKvLoop "github.com/m4gshm/gollections/break/kv/loop"
	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/convert"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/slice"
)

func Test_ReduceSum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, _ := breakLoop.Reduce(breakLoop.From(s), op.Sum[int])
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_EmptyLoop(t *testing.T) {
	s := breakLoop.Of[int]()
	r, _ := breakLoop.Reduce(s, op.Sum[int])
	assert.Equal(t, 0, r)
}

func Test_Sum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, _ := breakLoop.Sum(breakLoop.From(s))
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_Convert(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r := breakLoop.Convert(breakLoop.From(s), strconv.Itoa)
	o, _ := breakLoop.Slice(r.Next)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, o)
}

func Test_IterWitErr(t *testing.T) {
	s := breakLoop.From(loop.Of("1", "3", "5", "7eee", "9", "11"))
	r := []int{}
	var outErr error
	for it, i, ok, err := breakLoop.Conv(s, strconv.Atoi).Start(); ok || err != nil; i, ok, err = it.Next() {
		if err != nil {
			outErr = err
			break
		}
		r = append(r, i)
	}

	assert.Error(t, outErr)
	assert.Equal(t, []int{1, 3, 5}, r)

	s = breakLoop.From(loop.Of("1", "3", "5", "7eee", "9", "11"))
	r = []int{}
	//ignore err
	for it, i, ok, err := breakLoop.Conv(s, strconv.Atoi).Start(); ok || err != nil; i, ok, err = it.Next() {
		if err == nil {
			r = append(r, i)
		}
	}
	assert.Equal(t, []int{1, 3, 5, 9, 11}, r)
}

func Test_NotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = breakLoop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = breakLoop.NotNil(source)
		expected = []*entity{{"first"}, {"third"}, {"fifth"}}
	)
	o, _ := breakLoop.Slice(result.Next)
	assert.Equal(t, expected, o)
}

func Test_ConvertPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = breakLoop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = breakLoop.PtrVal(source)
		expected = []entity{{"first"}, {}, {"third"}, {}, {"fifth"}}
	)
	o, _ := breakLoop.Slice(result.Next)
	assert.Equal(t, expected, o)
}

func Test_ConvertNotnilPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = breakLoop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = breakLoop.NoNilPtrVal(source)
		expected = []entity{{"first"}, {"third"}, {"fifth"}}
	)
	o, _ := breakLoop.Slice(result.Next)
	assert.Equal(t, expected, o)
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
	r := breakLoop.FilterAndConvert(breakLoop.From(s), even, strconv.Itoa)
	o, _ := breakLoop.Slice(r.Next)
	assert.Equal(t, []string{"4", "8"}, o)
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := breakLoop.ConvCheck(breakLoop.From(s), func(i int) (string, bool, error) { return strconv.Itoa(i), even(i), nil })
	o, _ := breakLoop.Slice(r.Next)
	assert.Equal(t, []string{"4", "8"}, o)
}

func Test_Flatt(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.Flat(breakLoop.From(md), as.Is)
	e := []int{1, 2, 3, 4, 5, 6}
	o, _ := breakLoop.Slice(f.Next)
	assert.Equal(t, e, o)
}

func Test_FlattFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FilterAndFlat(breakLoop.From(md), func(from []int) bool { return len(from) > 1 }, as.Is)
	e := []int{1, 2, 3, 5, 6}
	o, _ := breakLoop.Slice(f.Next)
	assert.Equal(t, e, o)
}

func Test_FlattElemFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FlattAndFilter(breakLoop.From(md), as.Is, even)
	e := []int{2, 4, 6}
	o, _ := breakLoop.Slice(f.Next)
	assert.Equal(t, e, o)
}

func Test_FilterAndFlattFilt(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FilterFlatFilter(breakLoop.From(md), func(from []int) bool { return len(from) > 1 }, as.Is, even)
	e := []int{2, 6}
	o, _ := breakLoop.Slice(f.Next)
	assert.Equal(t, e, o)
}

func Test_Filter(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	f := breakLoop.Filter(breakLoop.From(s), even)
	e := []int{4, 8}
	o, _ := breakLoop.Slice(f.Next)
	assert.Equal(t, e, o)
}

func Test_Filtering(t *testing.T) {
	r := breakLoop.Filt(breakLoop.From(loop.Of(1, 2, 3, 4, 5, 6)), func(i int) (bool, error) { return i%2 == 0, nil })
	o, _ := breakLoop.Slice(r.Next)
	assert.Equal(t, []int{2, 4, 6}, o)
}

func Test_MatchAny(t *testing.T) {
	elements := loop.Of(1, 2, 3, 4)

	ok, _ := breakLoop.HasAny(breakLoop.From(elements), eq.To(4))
	assert.True(t, ok)

	noOk, _ := breakLoop.HasAny(breakLoop.From(elements), more.Than(5))
	assert.False(t, noOk)
}

type Role struct {
	name string
}

type User struct {
	name  string
	age   int
	roles []Role
}

func (u User) Name() string  { return u.name }
func (u User) Age() int      { return u.age }
func (u User) Roles() []Role { return u.roles }

var users = []User{
	{name: "Bob", age: 26, roles: []Role{{"Admin"}, {"manager"}}},
	{name: "Alice", age: 35, roles: []Role{{"Manager"}}},
	{name: "Tom", age: 18}, {},
}

func Test_KeyValuer(t *testing.T) {
	m, _ := breakKvLoop.Group(breakLoop.KeyValue(breakLoop.From(loop.Of(users...)), User.Name, User.Age).Next)

	assert.Equal(t, m["Alice"], slice.Of(35))
	assert.Equal(t, m["Bob"], slice.Of(26))
	assert.Equal(t, m["Tom"], slice.Of(18))

	g, _ := breakLoop.Group(breakLoop.From(loop.Of(users...)), User.Name, User.Age)
	assert.Equal(t, m, g)
}

func Test_Keyer(t *testing.T) {
	m, _ := breakKvLoop.Group(breakLoop.ExtraKey(breakLoop.From(loop.Of(users...)), User.Name).Next)

	assert.Equal(t, m["Alice"], slice.Of(users[1]))
	assert.Equal(t, m["Bob"], slice.Of(users[0]))
	assert.Equal(t, m["Tom"], slice.Of(users[2]))

	g := loop.Group(loop.Of(users...), User.Name, as.Is)
	assert.Equal(t, m, g)
}

func Test_Valuer(t *testing.T) {
	bob, bobRoles, _, _ := breakLoop.ExtraValue(breakLoop.From(loop.Of(users...)), User.Roles).Next()

	assert.Equal(t, bob, users[0])
	assert.Equal(t, bobRoles, users[0].roles)
}

func Test_MultiValuer(t *testing.T) {
	l := breakLoop.ExtraVals(breakLoop.From(loop.Of(users...)), User.Roles)
	bob, bobRole, _, _ := l.Next()
	bob2, bobRole2, _, _ := l.Next()

	assert.Equal(t, bob, users[0])
	assert.Equal(t, bob2, users[0])
	assert.Equal(t, bobRole, users[0].roles[0])
	assert.Equal(t, bobRole2, users[0].roles[1])
}

func Test_MultipleKeyValuer(t *testing.T) {
	m, _ := breakKvLoop.Group(breakLoop.KeysValues(breakLoop.From(loop.Of(users...)),
		func(u User) ([]string, error) {
			return slice.Convert(u.roles, func(r Role) string { return strings.ToLower(r.name) }), nil
		},
		func(u User) ([]string, error) { return []string{u.name, strings.ToLower(u.name)}, nil },
	).Next)

	assert.Equal(t, m["admin"], slice.Of("Bob", "bob"))
	assert.Equal(t, m["manager"], slice.Of("Bob", "bob", "Alice", "alice"))
	assert.Equal(t, m[""], slice.Of("Tom", "tom", "", ""))
}
