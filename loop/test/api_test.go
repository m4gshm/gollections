package test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert/as"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	kvloopgroup "github.com/m4gshm/gollections/kv/loop/group"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/conv"
	"github.com/m4gshm/gollections/loop/convert"
	"github.com/m4gshm/gollections/loop/filter"
	"github.com/m4gshm/gollections/loop/first"
	"github.com/m4gshm/gollections/loop/flat"
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

func Test_ReduceeSum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, err := loop.Reducee(s, func(i1, i2 int) (int, error) { return i1 + i2, nil })
	assert.Equal(t, 1+3+5+7+9+11, r)
	assert.Nil(t, err)
}

func Test_EmptyLoop(t *testing.T) {
	s := loop.Of[int]()
	r := loop.Reduce(s, op.Sum[int])
	assert.Equal(t, 0, r)
}

func Test_NilLoop(t *testing.T) {
	var s loop.Loop[int] = nil
	r := loop.Reduce(s, op.Sum[int])
	assert.Equal(t, 0, r)
}

func Test_ConvertAndReduce(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r := convert.AndReduce(s, func(i int) int { return i * i }, op.Sum[int])
	assert.Equal(t, 1+3*3+5*5+7*7+9*9+11*11, r)
}

func Test_ConvAndReduce(t *testing.T) {
	s := loop.Of("1", "3", "5", "7", "9", "11")
	r, err := conv.AndReduce(s, strconv.Atoi, op.Sum[int])
	assert.NoError(t, err)
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

func Test_FirstConverted(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, ok := first.Converted(s, func(i int) bool { return i > 5 }, strconv.Itoa)
	assert.True(t, ok)
	assert.Equal(t, "7", r)
}

func Test_NotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = loop.NotNil(source)
		expected = []*entity{{"first"}, {"third"}, {"fifth"}}
	)
	assert.Equal(t, expected, loop.Slice(result))
}

func Test_ConvertPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = loop.PtrVal(source)
		expected = []entity{{"first"}, {}, {"third"}, {}, {"fifth"}}
	)
	assert.Equal(t, expected, loop.Slice(result))
}

func Test_ConvertNotnilPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = loop.NoNilPtrVal(source)
		expected = []entity{{"first"}, {"third"}, {"fifth"}}
	)
	assert.Equal(t, expected, loop.Slice(result))
}

func Test_Convert(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r := loop.Convert(s, strconv.Itoa)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, loop.Slice(r))
}

func Test_IterWitErr(t *testing.T) {
	s := loop.Of("1", "3", "5", "7eee", "9", "11")
	r := []int{}
	var outErr error
	for {
		it := loop.Conv(s, strconv.Atoi)
		i, ok, err := it()
		if !ok && err == nil {
			break
		}
		if err != nil {
			outErr = err
			break
		}
		r = append(r, i)
	}

	assert.Error(t, outErr)
	assert.Equal(t, []int{1, 3, 5}, r)

	s = loop.Of("1", "3", "5", "7eee", "9", "11")
	r = []int{}
	//ignore err
	for {
		it := loop.Conv(s, strconv.Atoi)
		i, ok, err := it()
		if !ok && err == nil {
			break
		}
		if err == nil {
			r = append(r, i)
		}
	}
	assert.Equal(t, []int{1, 3, 5, 9, 11}, r)
}

func Test_IterStart(t *testing.T) {
	l := loop.Convert(loop.Of(1, 3, 5, 7, 9, 11), strconv.Itoa)
	r := []string{}

	for {
		i, ok := l()
		if !ok {
			break
		}
		r = append(r, i)
	}
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, r)
}

func Test_ConvertNotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = loop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = convert.NotNil(source, func(e *entity) string { return e.val })
		expected = []string{"first", "third", "fifth"}
	)
	assert.Equal(t, expected, loop.Slice(result))
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
	assert.Equal(t, expected, loop.Slice(result))
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
	assert.Equal(t, expected, loop.Slice(result))
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.FilterAndConvert(s, even, strconv.Itoa)
	assert.Equal(t, []string{"4", "8"}, loop.Slice(r))
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.ConvertCheck(s, func(i int) (string, bool) { return strconv.Itoa(i), even(i) })
	assert.Equal(t, []string{"4", "8"}, loop.Slice(r))
}

func Test_Flat(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.Flat(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, loop.Slice(f))
}

func Test_FlatAndConvert(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := flat.AndConvert(md, func(i []int) []int { return i }, strconv.Itoa)
	e := []string{"1", "2", "3", "4", "5", "6"}
	assert.Equal(t, e, loop.Slice(f))
}

func Test_FlatFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FilterAndFlat(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i })
	e := []int{1, 2, 3, 5, 6}
	assert.Equal(t, e, loop.Slice(f))
}

func Test_FlattElemFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FlattAndFilter(md, func(i []int) []int { return i }, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, loop.Slice(f))
}

func Test_FilterAndFlattFilt(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := loop.FilterFlatFilter(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i }, even)
	e := []int{2, 6}
	assert.Equal(t, e, loop.Slice(f))
}

func Test_Filter(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.Filter(s, even)
	assert.Equal(t, slice.Of(4, 8), loop.Slice(r))
}

func Test_FilterConvertFilter(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := filter.ConvertFilter(s, even, func(i int) int { return i * 2 }, even)
	assert.Equal(t, slice.Of(8, 16), loop.Slice(r))
}

func Test_Filt(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	l := loop.Filt(s, func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errors.New("abort"), nil) })
	r, err := breakLoop.Slice(l)
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4, 8), r)
}

func Test_Filt2(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	l := loop.Filt(s, func(i int) (bool, error) {
		ok := i <= 7
		return ok && even(i), op.IfElse(ok, nil, errors.New("abort"))
	})
	r, err := breakLoop.Slice(l)
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4), r)
}

func Test_FiltAndConv(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := loop.FiltAndConv(s, func(v int) (bool, error) { return v%2 == 0, nil }, func(i int) (int, error) { return i * 2, nil })
	o, _ := breakLoop.Slice(r)
	assert.Equal(t, slice.Of(8, 16), o)
}

func Test_Filtering(t *testing.T) {
	r := loop.Filter(loop.Of(1, 2, 3, 4, 5, 6), func(i int) bool { return i%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, loop.Slice(r))
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
	m := kvloop.Group(loop.KeyValue(loop.Of(users...), User.Name, User.Age))

	assert.Equal(t, m["Alice"], slice.Of(35))
	assert.Equal(t, m["Bob"], slice.Of(26))
	assert.Equal(t, m["Tom"], slice.Of(18))

	g := loop.Group(loop.Of(users...), User.Name, User.Age)
	assert.Equal(t, m, g)
}

func Test_Keyer(t *testing.T) {
	m := kvloop.Group(loop.ExtraKey(loop.Of(users...), User.Name))

	assert.Equal(t, m["Alice"], slice.Of(users[1]))
	assert.Equal(t, m["Bob"], slice.Of(users[0]))
	assert.Equal(t, m["Tom"], slice.Of(users[2]))

	g := loop.Group(loop.Of(users...), User.Name, as.Is)
	assert.Equal(t, m, g)
}

func Test_Valuer(t *testing.T) {
	bob, bobRoles, _ := loop.ExtraValue(loop.Of(users...), User.Roles)()

	assert.Equal(t, bob, users[0])
	assert.Equal(t, bobRoles, users[0].roles)
}

func Test_MultiValuer(t *testing.T) {
	l := loop.ExtraVals(loop.Of(users...), User.Roles)
	bob, bobRole, _ := l()
	bob2, bobRole2, _ := l()

	assert.Equal(t, bob, users[0])
	assert.Equal(t, bob2, users[0])
	assert.Equal(t, bobRole, users[0].roles[0])
	assert.Equal(t, bobRole2, users[0].roles[1])
}

func Test_MultipleKeyValuer(t *testing.T) {
	m := kvloop.Group(loop.KeysValues(loop.Of(users...),
		func(u User) []string {
			return slice.Convert(u.roles, func(r Role) string { return strings.ToLower(r.name) })
		},
		func(u User) []string { return []string{u.name, strings.ToLower(u.name)} },
	))

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

func Test_Sequence(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), loop.Slice(loop.Sequence(-1, func(prev int) (int, bool) { return prev + 1, prev < 3 })))
}

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")
	result := loop.Slice(loop.OfIndexed(len(indexed), func(i int) string { return indexed[i] }))
	assert.Equal(t, indexed, result)
}

func Test_ConvertIndexed(t *testing.T) {
	indexed := slice.Of(10, 11, 12, 13, 14)
	result := loop.Slice(convert.FromIndexed(len(indexed), func(i int) int { return indexed[i] }, strconv.Itoa))
	assert.Equal(t, slice.Of("10", "11", "12", "13", "14"), result)
}

func Test_ConvIndexed(t *testing.T) {
	indexed := slice.Of("10", "11", "12", "13", "14")
	result, err := breakLoop.Slice(conv.FromIndexed(len(indexed), func(i int) string { return indexed[i] }, strconv.Atoi))
	assert.NoError(t, err)
	assert.Equal(t, slice.Of(10, 11, 12, 13, 14), result)
}

func Test_Containt(t *testing.T) {
	assert.True(t, loop.Contains(loop.Of(1, 2, 3), 3))
	assert.False(t, loop.Contains(loop.Of(1, 2, 3), 0))
}

func Test_New(t *testing.T) {
	source := []string{"one", "two", "three"}
	i := 0
	l := loop.New(source, func(s []string) bool { return i < len(s) }, func(s []string) string { o := s[i]; i++; return o })

	assert.Equal(t, source, loop.Slice(l))
}

func Test_For(t *testing.T) {
	var out []int
	err := loop.For(loop.Of(1, 2, 3, 4), func(i int) error {
		if i == 3 {
			return c.Break
		}
		out = append(out, i)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, slice.Of(1, 2), out)
}

func Test_ForEachFiltered(t *testing.T) {
	var out []int
	loop.ForEachFiltered(loop.Of(1, 2, 3, 4), even, func(i int) { out = append(out, i) })
	assert.Equal(t, slice.Of(2, 4), out)
}

func Test_FlatValues(t *testing.T) {
	g := kvloopgroup.Of(loop.KeyValues(loop.Of(users...), func(u User) string { return u.name }, func(u User) []int { return slice.Of(u.age) }))

	assert.Equal(t, g["Bob"], slice.Of(26))
}

func Test_FlatKeys(t *testing.T) {
	g := kvloopgroup.Of(loop.KeysValue(loop.Of(users...), func(u User) []string { return slice.Of(u.name) }, func(u User) int { return u.age }))
	assert.Equal(t, g["Alice"], slice.Of(35))
}
