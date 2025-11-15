package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/vector"
	"github.com/m4gshm/gollections/seq"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
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

func Test_Vector_Empty_IAll(t *testing.T) {
	expected := []int{}
	v := &mutable.Vector[int]{}
	result := make([]int, v.Len())
	for _, it := range v.IAll {
		result = append(result, it)
	}
	assert.Equal(t, expected, result)
}

func Test_Vector_Empty_All(t *testing.T) {
	expected := []int{}
	v := &mutable.Vector[int]{}
	result := make([]int, v.Len())
	for it := range v.All {
		result = append(result, it)
	}
	assert.Equal(t, expected, result)
}

func Test_VectorIterateOverRange(t *testing.T) {
	expected := slice.Of(1, 2, 3, 4)
	v := vector.Of(1, 2, 3, 4)
	result := make([]int, v.Len())
	copy(result, *v)
	assert.Equal(t, expected, result)
}

func Test_Vector_Sort(t *testing.T) {
	var (
		elements = vector.Of(3, 1, 5, 6, 8, 0, -2)
		sorted   = elements.Sort(op.Compare)
	)
	assert.Equal(t, vector.Of(-2, 0, 1, 3, 5, 6, 8), sorted)
	assert.Equal(t, vector.Of(-2, 0, 1, 3, 5, 6, 8), elements)
	assert.Same(t, elements, sorted)
}

func Test_Vector_SortStructByField(t *testing.T) {
	var (
		anonymous = &user{"Anonymous", 0}
		cherlie   = &user{"Cherlie", 25}
		alise     = &user{"Alise", 20}
		bob       = &user{"Bob", 19}

		elements     = vector.Of(anonymous, cherlie, alise, bob)
		sortedByName = vector.Sort(elements.Clone(), (*user).Name)
		sortedByAge  = vector.Sort(elements.Clone(), (*user).Age)
	)
	assert.Equal(t, vector.Of(alise, anonymous, bob, cherlie), sortedByName)
	assert.Equal(t, vector.Of(anonymous, bob, alise, cherlie), sortedByAge)
}

func Test_Vector_Nil(t *testing.T) {
	var vec *mutable.Vector[int]
	var nils []int
	vec.Add(1, 2, 3)
	vec.Add(nils...)
	vec.AddOne(4)
	vec.AddAll(vec.All)

	vec.Delete(1, 2, 3)
	vec.Delete(nils...)
	vec.DeleteOne(4)

	vec.IsEmpty()
	vec.Len()

	vec.ForEach(nil)
	vec.TrackEach(nil)

	assert.Equal(t, nils, vec.Slice())

	_, ok := vec.Head()
	assert.False(t, ok)

	_, ok = vec.Tail()
	assert.False(t, ok)
}

func Test_Vector_Zero(t *testing.T) {
	var vec mutable.Vector[string]

	var nilValues []string
	vec.Add("a", "b", "c")
	vec.Add(nilValues...)
	vec.AddOne("d")
	vec.AddAll(vec.All)

	vec.Delete(0, 1, 2)
	var nilIndexes []int
	vec.Delete(nilIndexes...)
	vec.DeleteOne(0)

	e := vec.IsEmpty()
	assert.False(t, e)

	l := vec.Len()
	assert.Equal(t, 4, l)

	vec.ForEach(func(_ string) {})
	vec.TrackEach(func(_ int, _ string) {})

	assert.Equal(t, slice.Of("a", "b", "c", "d"), vec.Slice())

	head, ok := vec.Head()
	assert.True(t, ok)
	assert.Equal(t, "a", head)

	tail, ok := vec.Tail()
	assert.True(t, ok)
	assert.Equal(t, "d", tail)
}

func Test_Vector_new(t *testing.T) {
	var vec = new(mutable.Vector[string])

	var nilValues []string
	vec.Add("a", "b", "c")
	vec.Add(nilValues...)
	vec.AddOne("d")
	vec.AddAll(vec.All)

	vec.Delete(0, 1, 2)
	var nilIndexes []int
	vec.Delete(nilIndexes...)
	vec.DeleteOne(0)

	e := vec.IsEmpty()
	assert.False(t, e)

	l := vec.Len()
	assert.Equal(t, 4, l)

	vec.ForEach(func(_ string) {})
	vec.TrackEach(func(_ int, _ string) {})

	assert.Equal(t, slice.Of("a", "b", "c", "d"), vec.Slice())

	head, ok := vec.Head()
	assert.True(t, ok)
	assert.Equal(t, "a", head)

	tail, ok := vec.Tail()
	assert.True(t, ok)
	assert.Equal(t, "d", tail)
}

func Test_Vector_AddAllOfSelf(t *testing.T) {
	vec := vector.Of(1, 2, 3)
	vec.AddAll(vec.All)
	assert.Equal(t, vector.Of(1, 2, 3, 1, 2, 3), vec)
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }

func Test_Vector_AddAndDelete(t *testing.T) {
	vec := vector.NewCap[rune](0)
	vec.Add(range_.Of('a', 'a'+rune(1000))...)
	assert.Equal(t, 1000, vec.Len())
	vec.Delete(range_.Of(0, 1000)...)
	assert.True(t, vec.IsEmpty())
}

func Test_Vector_Add(t *testing.T) {
	vec := vector.NewCap[int](0)
	vec.Add(1, 1, 2, 4, 3, 1)
	assert.Equal(t, slice.Of(1, 1, 2, 4, 3, 1), vec.Slice())
	vec.Add(1)
	assert.Equal(t, slice.Of(1, 1, 2, 4, 3, 1, 1), vec.Slice())
}

func Test_Vector_AddAll(t *testing.T) {
	vec := vector.NewCap[int](0)
	vec.AddAll(seq.Of(1, 1, 2, 4, 3, 1))
	assert.Equal(t, slice.Of(1, 1, 2, 4, 3, 1), vec.Slice())
	vec.AddAll(seq.Of(1))
	assert.Equal(t, slice.Of(1, 1, 2, 4, 3, 1, 1), vec.Slice())
}

func Test_Vector_DeleteOne(t *testing.T) {
	vec := vector.Of("1", "1", "2", "4", "3", "1")
	vec.DeleteOne(3)
	assert.Equal(t, slice.Of("1", "1", "2", "3", "1"), vec.Slice())
	r := vec.DeleteActualOne(5)
	assert.False(t, r)
}

func Test_Vector_DeleteMany(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	vec.Delete(3, 0, 5)
	assert.Equal(t, slice.Of("1", "2", "4", "6"), vec.Slice())
	r := vec.DeleteActual(5, 4)
	assert.False(t, r)
}

func Test_Vector_DeleteManyFromTail(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	vec.Delete(4, 5, 6)
	assert.Equal(t, slice.Of("0", "1", "2", "3"), vec.Slice())
}

func Test_Vector_DeleteManyFromHead(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	vec.Delete(0, 1, 2)
	assert.Equal(t, slice.Of("3", "4", "5", "6"), vec.Slice())
}

func Test_Vector_DeleteManyFromMiddle(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	vec.Delete(4, 3)
	assert.Equal(t, slice.Of("0", "1", "2", "5", "6"), vec.Slice())
}

func Test_Vector_Set(t *testing.T) {
	vec := vector.Of("1", "1", "2", "4", "3", "1")
	vec.Set(10, "11")
	assert.Equal(t, slice.Of("1", "1", "2", "4", "3", "1", "", "", "", "", "11"), vec.Slice())
}

func Test_Vector_FilterMapReduce(t *testing.T) {
	s := vector.Of(1, 1, 2, 4, 3, 4).Filter(func(i int) bool { return i%2 == 0 }).Convert(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 20, s)
}
