package test

import (
	"errors"
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

func Test_AccumSum(t *testing.T) {
	s := breakLoop.Of(1, 3, 5, 7, 9, 11)
	r, err := breakLoop.Accum(100, s, op.Sum[int])
	assert.Equal(t, 100+1+3+5+7+9+11, r)
	assert.NoError(t, err)
}

func Test_AccummSum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, err := loop.Accumm(100, s, func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errors.New("stop")
		}
		return i1 + i2, nil
	})
	assert.Equal(t, 100+1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceSum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, ok, err := breakLoop.ReduceOK(breakLoop.From(s), op.Sum[int])
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_ReduceeSum(t *testing.T) {
	s := loop.Of(1, 3, 5, 7, 9, 11)
	r, ok, err := breakLoop.ReduceeOK(breakLoop.From(s), func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errors.New("stop")
		}
		return i1 + i2, nil
	})
	assert.ErrorContains(t, err, "stop")
	assert.True(t, ok)
	assert.Equal(t, 1+3+5+7+9, r)
}

func Test_ReduceeSumFirstErr(t *testing.T) {
	var tru breakLoop.Loop[int] = func() (int, bool, error) {
		return 1, true, errors.New("first-err")
	}
	r, ok, err := breakLoop.ReduceeOK(tru, func(i1, i2 int) (int, error) {
		return i1 + i2, nil
	})
	assert.ErrorContains(t, err, "first-err")
	assert.True(t, ok)
	assert.Equal(t, 1, r)

	var fals breakLoop.Loop[int] = func() (int, bool, error) {
		return 2, false, errors.New("first-err")
	}

	r, ok, err = breakLoop.ReduceeOK(fals, func(i1, i2 int) (int, error) {
		return i1 + i2, nil
	})
	assert.ErrorContains(t, err, "first-err")
	assert.False(t, ok)
	assert.Equal(t, 2, r)
}

func Test_Firstt(t *testing.T) {
	result, ok, err := loop.Firstt(loop.Of(1, 2, 3, 4, 5, 6), func(i int) (bool, error) {
		return more.Than(5)(i), nil
	})

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	result, ok, err = loop.Firstt(loop.Of(1, 2, 3, 4, 5, 6), func(_ int) (bool, error) { return true, errors.New("abort") })

	assert.True(t, ok)
	assert.Equal(t, 1, result)
	assert.ErrorContains(t, err, "abort")

	_, ok, err = loop.Firstt(loop.Of(1, 2, 3, 4, 5, 6), func(_ int) (bool, error) { return false, errors.New("abort") })

	assert.False(t, ok)
	assert.ErrorContains(t, err, "abort")

	// _, ok, _ = loop.Firstt(loop.Of(1, 2, 3, 4, 5, 6), nil)
	// assert.False(t, ok)

	_, ok, _ = loop.Firstt(nil, func(_ int) (bool, error) { return false, errors.New("abort") })
	assert.False(t, ok)
}

func Test_ReduceeEmptyLoop(t *testing.T) {
	s := breakLoop.Of[int]()
	r, ok, err := breakLoop.ReduceOK(s, op.Sum[int])
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, 0, r)
}

func Test_ReduceeNilLoop(t *testing.T) {
	var s breakLoop.Loop[int]
	r, ok, err := breakLoop.ReduceOK(s, op.Sum[int])
	assert.NoError(t, err)
	assert.False(t, ok)
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
	o, _ := breakLoop.Slice(r)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, o)
}

func Test_IterWitErr(t *testing.T) {
	s := breakLoop.From(loop.Of("1", "3", "5", "7eee", "9", "11"))
	r := []int{}
	var outErr error
	for {
		it := breakLoop.Conv(s, strconv.Atoi)
		i, ok, err := it()
		if err != nil {
			assert.False(t, ok)
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
	for {
		it := breakLoop.Conv(s, strconv.Atoi)
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

func Test_NotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = breakLoop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = breakLoop.NotNil(source)
		expected = []*entity{{"first"}, {"third"}, {"fifth"}}
	)
	o, _ := breakLoop.Slice(result)
	assert.Equal(t, expected, o)
}

func Test_ConvertPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = breakLoop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = breakLoop.PtrVal(source)
		expected = []entity{{"first"}, {}, {"third"}, {}, {"fifth"}}
	)
	o, _ := breakLoop.Slice(result)
	assert.Equal(t, expected, o)
}

func Test_ConvertNotnilPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = breakLoop.Of([]*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}...)
		result   = breakLoop.NoNilPtrVal(source)
		expected = []entity{{"first"}, {"third"}, {"fifth"}}
	)
	o, _ := breakLoop.Slice(result)
	assert.Equal(t, expected, o)
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
	r := breakLoop.FilterAndConvert(breakLoop.From(s), even, strconv.Itoa)
	o, _ := breakLoop.Slice(r)
	assert.Equal(t, []string{"4", "8"}, o)
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := breakLoop.ConvOK(breakLoop.From(s), func(i int) (string, bool, error) { return strconv.Itoa(i), even(i), nil })
	o, _ := breakLoop.Slice(r)
	assert.Equal(t, []string{"4", "8"}, o)
}

func Test_Flatt(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.Flat(breakLoop.From(md), as.Is)
	e := []int{1, 2, 3, 4, 5, 6}
	o, _ := breakLoop.Slice(f)
	assert.Equal(t, e, o)
}

func Test_FlattFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FilterAndFlat(breakLoop.From(md), func(from []int) bool { return len(from) > 1 }, as.Is)
	e := []int{1, 2, 3, 5, 6}
	o, _ := breakLoop.Slice(f)
	assert.Equal(t, e, o)
}

func Test_FlattElemFilter(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FlattAndFilter(breakLoop.From(md), as.Is, even)
	e := []int{2, 4, 6}
	o, _ := breakLoop.Slice(f)
	assert.Equal(t, e, o)
}

func Test_FilterAndFlattFilt(t *testing.T) {
	md := loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := breakLoop.FilterFlatFilter(breakLoop.From(md), func(from []int) bool { return len(from) > 1 }, as.Is, even)
	e := []int{2, 6}
	o, _ := breakLoop.Slice(f)
	assert.Equal(t, e, o)
}

func Test_Filter(t *testing.T) {
	s := loop.Of(1, 3, 4, 5, 7, 8, 9, 11)
	f := breakLoop.Filter(breakLoop.From(s), even)
	e := []int{4, 8}
	o, _ := breakLoop.Slice(f)
	assert.Equal(t, e, o)
}

func Test_Filtering(t *testing.T) {
	r := breakLoop.Filt(breakLoop.From(loop.Of(1, 2, 3, 4, 5, 6)), func(i int) (bool, error) { return i%2 == 0, nil })
	o, _ := breakLoop.Slice(r)
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
	m, _ := breakKvLoop.Group(breakLoop.KeyValue(breakLoop.From(loop.Of(users...)), User.Name, User.Age))

	assert.Equal(t, m["Alice"], slice.Of(35))
	assert.Equal(t, m["Bob"], slice.Of(26))
	assert.Equal(t, m["Tom"], slice.Of(18))

	g, _ := breakLoop.Group(breakLoop.From(loop.Of(users...)), User.Name, User.Age)
	assert.Equal(t, m, g)
}

func Test_Keyer(t *testing.T) {
	m, _ := breakKvLoop.Group(breakLoop.ExtraKey(breakLoop.From(loop.Of(users...)), User.Name))

	assert.Equal(t, m["Alice"], slice.Of(users[1]))
	assert.Equal(t, m["Bob"], slice.Of(users[0]))
	assert.Equal(t, m["Tom"], slice.Of(users[2]))

	g := loop.Group(loop.Of(users...), User.Name, as.Is)
	assert.Equal(t, m, g)
}

func Test_Valuer(t *testing.T) {
	bob, bobRoles, _, _ := breakLoop.ExtraValue(breakLoop.From(loop.Of(users...)), User.Roles)()

	assert.Equal(t, bob, users[0])
	assert.Equal(t, bobRoles, users[0].roles)
}

func Test_MultiValuer(t *testing.T) {
	l := breakLoop.ExtraVals(breakLoop.From(loop.Of(users...)), User.Roles)
	bob, bobRole, _, _ := l()
	bob2, bobRole2, _, _ := l()

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
	))

	assert.Equal(t, m["admin"], slice.Of("Bob", "bob"))
	assert.Equal(t, m["manager"], slice.Of("Bob", "bob", "Alice", "alice"))
	assert.Equal(t, m[""], slice.Of("Tom", "tom", "", ""))
}
