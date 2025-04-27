package sliceexamples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/one"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/convert"
	"github.com/m4gshm/gollections/slice/filter"

	"github.com/m4gshm/gollections/slice/reverse"
)

func Test_SortStructs(t *testing.T) {

	var users = []User{
		{name: "Bob", age: 26},
		{name: "Alice", age: 35},
		{name: "Tom", age: 18},
		{name: "Chris", age: 41},
	}
	var byName = slice.Sort(slice.Clone(users), func(u1, u2 User) int { return op.Compare(u1.name, u2.name) })
	var byAgeReverse = slice.Sort(slice.Clone(users), func(u1, u2 User) int { return -op.Compare(u1.age, u2.age) })

	assert.Equal(t, []User{
		{name: "Alice", age: 35},
		{name: "Bob", age: 26},
		{name: "Chris", age: 41},
		{name: "Tom", age: 18},
	}, byName)

	assert.Equal(t, []User{
		{name: "Chris", age: 41},
		{name: "Alice", age: 35},
		{name: "Bob", age: 26},
		{name: "Tom", age: 18},
	}, byAgeReverse)

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
	assert.NotSame(t, &entities, &c)

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
	assert.NotSame(t, &entities, &c)

	for i := range entities {
		assert.Equal(t, entities[i], c[i])
		assert.NotSame(t, entities[i], c[i])
	}
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
		result = convert.AndFilter(source, strconv.Itoa, func(s string) bool {
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
		result = convert.NilSafe(source, func(e *entity) *string {
			return e.val
		})
		expected = []*string{&first, &third, &fifth}
	)
	assert.Equal(t, expected, result)
}

func Test_ConvertFilteredWithIndexInPlace(t *testing.T) {
	var (
		source = slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
		result = convert.CheckIndexed(source, func(index int, elem int) (string, bool) {
			return strconv.Itoa(index + elem), even(elem)
		})
		expected = []string{"6", "13"}
	)
	assert.Equal(t, expected, result)
}

func Test_Slice_Filter(t *testing.T) {

	var even = slice.Filter([]int{1, 2, 3, 4, 5, 6}, func(v int) bool { return v%2 == 0 })
	//[]int{2, 4, 6}

	assert.Equal(t, []int{2, 4, 6}, even)
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

func Test_ConvertPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = []*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}
		result   = slice.ToValues(source)
		expected = []entity{{"first"}, {}, {"third"}, {}, {"fifth"}}
	)
	assert.Equal(t, expected, result)
}

func Test_ConvertNotnilPointersToValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = []*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}
		result   = slice.GetValues(source)
		expected = []entity{{"first"}, {"third"}, {"fifth"}}
	)
	assert.Equal(t, expected, result)
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

type Rows[T any] struct {
	row    []T
	cursor int
}

func (r *Rows[T]) Next() bool         { return r.cursor < len(r.row) }
func (r *Rows[T]) Scan(dest *T) error { *dest = r.row[r.cursor]; r.cursor++; return nil }

func Test_OffNextPush(t *testing.T) {
	var (
		rows        = &Rows[int]{slice.Of(1, 2, 3), 0}
		result, err = slice.OfNextPush(rows.Next, rows.Scan)
		expected    = slice.Of(1, 2, 3)
	)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}
