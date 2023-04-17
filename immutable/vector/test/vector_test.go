package vector

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/slice"
)

func Test_VectorIterate(t *testing.T) {
	expected := slice.Of(1, 2, 3, 4)
	v := vector.Of(1, 2, 3, 4)
	result := make([]int, v.Len())
	i := 0
	for it := v.Head(); it.HasNext(); {
		result[i] = it.GetNext()
		i++
	}
	assert.Equal(t, expected, result)
}

func Test_VectorIterate2(t *testing.T) {
	expected := slice.Of(1, 2, 3, 4)
	v := vector.Of(1, 2, 3, 4)
	result := make([]int, v.Len())
	i := 0
	for it, v, ok := v.First(); ok; v, ok = it.Next() {
		result[i] = v
		i++
	}
	assert.Equal(t, expected, result)
}

func Test_VectorIterate3(t *testing.T) {
	expected := slice.Of(1, 2, 3, 4)
	v := vector.Of(1, 2, 3, 4)
	result := make([]int, v.Len())
	i := 0
	it := v.Head()
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		result[i] = v
		i++
	}
	assert.Equal(t, expected, result)
}

func Test_VectorReverseIteration(t *testing.T) {
	expected := slice.Of(4, 3, 2, 1)
	v := vector.Of(1, 2, 3, 4)
	result := make([]int, v.Len())
	i := 0
	for it := v.Tail(); it.HasPrev(); {
		result[i] = it.GetPrev()
		i++
	}
	assert.Equal(t, expected, result)
}

func Test_VectorReverseIteration2(t *testing.T) {
	expected := slice.Of(4, 3, 2, 1)
	v := vector.Of(1, 2, 3, 4)
	result := make([]int, v.Len())
	i := 0
	for it, v, ok := v.Last(); ok; v, ok = it.Prev() {
		result[i] = v
		i++
	}
	assert.Equal(t, expected, result)
}

func Test_Vector_Sort(t *testing.T) {
	var (
		elements = vector.Of(3, 1, 5, 6, 8, 0, -2)
		sorted   = elements.Sort(func(e1, e2 int) bool { return e1 < e2 })
	)
	assert.Equal(t, vector.Of(-2, 0, 1, 3, 5, 6, 8), sorted)
}

func Test_Vector_SortStructByField(t *testing.T) {
	var (
		anonymous = &user{"Anonymous", 0}
		cherlie   = &user{"Cherlie", 25}
		alise     = &user{"Alise", 20}
		bob       = &user{"Bob", 19}

		elements     = vector.Of(anonymous, cherlie, alise, bob)
		sortedByName = vector.Sort(elements, (*user).Name)
		sortedByAge  = vector.Sort(elements, (*user).Age)
	)
	assert.Equal(t, vector.Of(alise, anonymous, bob, cherlie), sortedByName)
	assert.Equal(t, vector.Of(anonymous, bob, alise, cherlie), sortedByAge)
}

func Test_Vector_Convert(t *testing.T) {
	var (
		ints     = vector.Of(3, 1, 5, 6, 8, 0, -2)
		strings  = it.ToSlice(it.Filter(vector.Convert(ints, strconv.Itoa), func(s string) bool { return len(s) == 1 }))
		strings2 = vector.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }).Slice()
	)
	assert.Equal(t, slice.Of("3", "1", "5", "6", "8", "0"), strings)
	assert.Equal(t, strings, strings2)
}

func Test_Vector_Flatt(t *testing.T) {
	var (
		deepInts    = vector.Of(vector.Of(3, 1), vector.Of(5, 6, 8, 0, -2))
		ints        = vector.Flatt(deepInts, immutable.Vector[int].Collect)
		stringsPipe = iterable.Filter(iterable.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }), func(s string) bool { return len(s) == 1 })
	)
	assert.Equal(t, slice.Of("3", "1", "5", "6", "8", "0"), stringsPipe.Slice())
}

func Test_Vector_DoubleConvert(t *testing.T) {
	var (
		ints               = vector.Of(3, 1, 5, 6, 8, 0, -2)
		stringsPipe        = vector.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 })
		prefixedStrinsPipe = iterable.Convert(stringsPipe, func(s string) string { return "_" + s })
	)
	assert.Equal(t, slice.Of("_3", "_1", "_5", "_6", "_8", "_0"), prefixedStrinsPipe.Slice())

	//second call do nothing
	var no []string
	assert.Equal(t, no, stringsPipe.Slice())
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }
