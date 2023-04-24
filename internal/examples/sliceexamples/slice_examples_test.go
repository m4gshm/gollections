package sliceexamples

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/last"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/exclude"
	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/predicate/not"
	"github.com/m4gshm/gollections/predicate/one"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/clone/sort"
	sliceConvert "github.com/m4gshm/gollections/slice/convert"
	"github.com/m4gshm/gollections/slice/filter"
	"github.com/m4gshm/gollections/slice/flatt"
	"github.com/m4gshm/gollections/slice/group"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/slice/reverse"
	"github.com/m4gshm/gollections/sum"
)

type User struct {
	name  string
	age   int
	roles []Role
}

type Role struct {
	name string
}

func (u Role) Name() string {
	return u.name
}

func (u User) Name() string {
	return u.name
}

func (u User) Age() int {
	return u.age
}

func (u User) Roles() []Role {
	return u.roles
}

var users = []User{
	{name: "Bob", age: 26, roles: []Role{{"Admin"}, {"manager"}}},
	{name: "Alice", age: 35, roles: []Role{{"Manager"}}},
	{name: "Tom", age: 18},
}

func Test_GroupBySeveralKeysAndConvertMapValues(t *testing.T) {
	usersByRole := group.InMultiple(users, func(u User) []string {
		return sliceConvert.AndConvert(u.Roles(), Role.Name, strings.ToLower)
	})
	namesByRole := map_.ConvertValues(usersByRole, func(u []User) []string {
		return slice.Convert(u, User.Name)
	})

	assert.Equal(t, namesByRole[""], []string{"Tom"})
	assert.Equal(t, namesByRole["manager"], []string{"Bob", "Alice"})
	assert.Equal(t, namesByRole["admin"], []string{"Bob"})
}

func Test_FindFirsManager(t *testing.T) {
	alice, _ := slice.First(users, func(user User) bool {
		roles := slice.Convert(user.Roles(), Role.Name)
		return slice.Contains(roles, "Manager")
	})

	assert.Equal(t, "Alice", alice.Name())
}

func Test_AggregateFilteredRoles(t *testing.T) {
	roles := flatt.AndConvert(users, User.Roles, Role.Name)
	roleNamesExceptManager := slice.Filter(roles, not.Eq("Manager"))

	assert.Equal(t, slice.Of("Admin", "manager"), roleNamesExceptManager)
}

func Test_SortStructs(t *testing.T) {
	var users = []User{
		{name: "Bob", age: 26},
		{name: "Alice", age: 35},
		{name: "Tom", age: 18},
	}
	var (
		//sorted
		byName = sort.By(users, User.Name)
		byAge  = sort.By(users, User.Age)
	)
	assert.Equal(t, []User{
		{name: "Alice", age: 35},
		{name: "Bob", age: 26},
		{name: "Tom", age: 18},
	}, byName)
	assert.Equal(t, []User{
		{name: "Tom", age: 18},
		{name: "Bob", age: 26},
		{name: "Alice", age: 35},
	}, byAge)
}

func Test_SortStructsByLess(t *testing.T) {
	var users = []User{
		{name: "Bob", age: 26},
		{name: "Alice", age: 35},
		{name: "Tom", age: 18},
	}
	var (
		//sorted
		byName       = sort.ByLess(users, func(u1, u2 User) bool { return u1.name < u2.name })
		byAgeReverse = sort.ByLess(users, func(u1, u2 User) bool { return u1.age > u2.age })
	)
	assert.Equal(t, []User{
		{name: "Alice", age: 35},
		{name: "Bob", age: 26},
		{name: "Tom", age: 18},
	}, byName)
	assert.Equal(t, []User{
		{name: "Alice", age: 35},
		{name: "Bob", age: 26},
		{name: "Tom", age: 18},
	}, byAgeReverse)
}

func Test_SortInt(t *testing.T) {
	source := []int{1, 3, -1, 2, 0}
	sorted := sort.Of(source)
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, sorted)
}

func Test_Reverse(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, slice.Reverse([]int{3, 2, 1, 0, -1}))
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, reverse.Of([]int{3, 2, 1, 0, -1}))
}

func Test_Clone(t *testing.T) {
	type entity struct{ val string }
	var (
		entities = []*entity{{"first"}, {"second"}, {"third"}}
		c        = clone.Of(entities)
	)

	assert.Equal(t, entities, c)
	assert.NotSame(t, entities, c)

	for i := range entities {
		assert.Same(t, entities[i], c[i])
	}
}

func Test_DeepClone(t *testing.T) {
	type entity struct{ val string }
	var (
		entities = []*entity{{"first"}, {"second"}, {"third"}}
		c        = clone.Deep(entities, clone.Ptr[entity])
	)

	assert.Equal(t, entities, c)
	assert.NotSame(t, entities, c)

	for i := range entities {
		assert.Equal(t, entities[i], c[i])
		assert.NotSame(t, entities[i], c[i])
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
		result   = filter.AndConvert(source, even, strconv.Itoa)
		expected = []string{"4", "8"}
	)
	assert.Equal(t, expected, result)
}

func Test_FilterConverted(t *testing.T) {
	var (
		source = []int{1, 3, 4, 5, 7, 8, 9, 11}
		result = sliceConvert.AndFilter(source, strconv.Itoa, func(s string) bool {
			return len(s) == 2
		})
		expected = []string{"11"}
	)
	assert.Equal(t, expected, result)
}

func Test_ConvertNilSafe(t *testing.T) {
	type entity struct{ val *string }
	var (
		first  = "first"
		third  = "third"
		fifth  = "fifth"
		source = []*entity{{&first}, {}, {&third}, nil, {&fifth}}
		result = sliceConvert.NilSafe(source, func(e *entity) *string {
			return e.val
		})
		expected = []*string{&first, &third, &fifth}
	)
	assert.Equal(t, expected, result)
}

func Test_ConvertFilteredWithIndexInPlace(t *testing.T) {
	var (
		source = slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
		result = sliceConvert.CheckIndexed(source, func(index int, elem int) (string, bool) {
			return strconv.Itoa(index + elem), even(elem)
		})
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

func Test_FilterNotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = []*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}
		result   = slice.NotNil(source)
		expected = []*entity{{"first"}, {"third"}, {"fifth"}}
	)
	assert.Equal(t, expected, result)
}

func Test_Flatt(t *testing.T) {
	var (
		source   = [][]int{{1, 2, 3}, {4}, {5, 6}}
		result   = slice.Flatt(source, as.Is[[]int])
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
		result   = slice.Flatt(source, as.Is[[]int])
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
		counter = 0
		result  = slice.Generate(func() (int, bool) {
			counter++
			return counter, counter < 4
		})
		expected = slice.Of(1, 2, 3)
	)
	assert.Equal(t, expected, result)
}
