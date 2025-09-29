package vector

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/collection/immutable/vector"
	"github.com/m4gshm/gollections/seq"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func Test_VectorIterate(t *testing.T) {
	expected := slice.Of(1, 2, 3, 4)
	v := vector.Of(1, 2, 3, 4)
	result := make([]int, v.Len())
	i := 0
	for it := range v.All {
		result[i] = it
		i++
	}
	assert.Equal(t, expected, result)
}

func Test_Vector_Sort(t *testing.T) {
	var (
		elements = vector.Of(3, 1, 5, 6, 8, 0, -2)
		sorted   = elements.Sort(op.Compare)
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
		strings  = seq.Slice(seq.Filter(vector.Convert(ints, strconv.Itoa), func(s string) bool { return len(s) == 1 }))
		strings2 = vector.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }).Slice()
	)
	assert.Equal(t, slice.Of("3", "1", "5", "6", "8", "0"), strings)
	assert.Equal(t, strings, strings2)
}

func Test_Vector_Flatt(t *testing.T) {
	var (
		deepInts    = vector.Of(vector.Of(3, 1), vector.Of(5, 6, 8, 0, -2))
		ints        = vector.Flat(deepInts, immutable.Vector[int].Slice)
		c           = seq.Convert(ints, strconv.Itoa)
		stringsPipe = c.Filter(func(s string) bool { return len(s) == 1 })
	)
	assert.Equal(t, slice.Of("3", "1", "5", "6", "8", "0"), stringsPipe.Slice())
}

func Test_Vector_Zero(t *testing.T) {
	var vec immutable.Vector[int]

	vec.IsEmpty()
	vec.Len()

	vec.ForEach(nil)
	vec.TrackEach(nil)

	vec.Slice()

	_, ok := vec.Head()

	assert.False(t, ok)

	_, ok = vec.Tail()

	assert.False(t, ok)
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }
