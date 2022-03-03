package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/mutable/vector"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/sum"
	"github.com/m4gshm/gollections/walk/group"
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

func Test_VectorIterateOverRange(t *testing.T) {
	expected := slice.Of(1, 2, 3, 4)
	v := vector.Of(1, 2, 3, 4)
	result := make([]int, v.Len())
	for i, val := range *v {
		result[i] = val
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

func Test_Vector_Sort(t *testing.T) {
	var (
		elements = vector.Of(3, 1, 5, 6, 8, 0, -2)
		sorted   = elements.Sort(func(e1, e2 int) bool { return e1 < e2 })
	)
	assert.Equal(t, vector.Of(-2, 0, 1, 3, 5, 6, 8), sorted)
	assert.Equal(t, vector.Of(-2, 0, 1, 3, 5, 6, 8), elements)
	assert.Equal(t, elements, sorted)
}

func Test_Vector_SortStructByField(t *testing.T) {
	var (
		anonymous = &user{"Anonymous", 0}
		cherlie   = &user{"Cherlie", 25}
		alise     = &user{"Alise", 20}
		bob       = &user{"Bob", 19}

		elements     = vector.Of(anonymous, cherlie, alise, bob)
		sortedByName = vector.Sort(elements.Copy(), (*user).Name)
		sortedByAge  = vector.Sort(elements.Copy(), (*user).Age)
	)
	assert.Equal(t, vector.Of(alise, anonymous, bob, cherlie), sortedByName)
	assert.Equal(t, vector.Of(anonymous, bob, alise, cherlie), sortedByAge)
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }

func Test_Vector_AddAndDelete(t *testing.T) {
	vec := vector.New[int](0)
	added := vec.AddAll(range_.Of(0, 1000))
	assert.Equal(t, added, true)
	deleted := false
	for i := vec.Head(); i.HasNext(); {
		deleted = i.DeleteNext()
	}
	assert.Equal(t, added, true)
	assert.Equal(t, deleted, true)
	assert.True(t, vec.IsEmpty())

	added = vec.AddAll(range_.Of(0, 10000))
	for i := vec.Tail(); i.HasPrev(); {
		deleted = i.DeletePrev()
	}
	assert.Equal(t, added, true)
	assert.Equal(t, deleted, true)
	assert.True(t, vec.IsEmpty())
}

func Test_Vector_Add(t *testing.T) {
	vec := vector.New[int](0)
	added := vec.Add(1, 1, 2, 4, 3, 1)
	assert.Equal(t, added, true)
	added = vec.Add(1)
	assert.Equal(t, added, true)
	assert.Equal(t, slice.Of(1, 1, 2, 4, 3, 1, 1), vec.Collect())
}

func Test_Vector_Add_And_Iterate(t *testing.T) {
	vec := vector.New[int](0)
	it, v, ok := vec.First()
	//no a first element
	assert.False(t, ok)
	//no more elements
	assert.False(t, it.HasNext())
	vec.Add(1)
	//exists one more
	assert.True(t, it.HasNext())
	v, ok = it.Get()
	//but the cursor points out of the range
	assert.False(t, ok)
	//starts itaration
	v, ok = it.Next()
	//success
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	//no more
	assert.False(t, it.HasNext())
	//no prev
	assert.False(t, it.HasPrev())
	//only the current one
	v, ok = it.Get()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
}

func Test_Vector_Delete_And_Iterate(t *testing.T) {
	vec := vector.Of(2)
	it, v, ok := vec.First()
	//success
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	//no more
	assert.False(t, it.HasNext())
	//no prev
	assert.False(t, it.HasPrev())
	//only the current one
	v, ok = it.Get()
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	//delete the element
	vec.DeleteOne(0)

	//no more
	assert.False(t, it.HasNext())
	//no prev
	assert.False(t, it.HasPrev())

	//no the current one
	_, ok = it.Get()
	assert.False(t, ok)

}

func Test_Vector_DeleteOne(t *testing.T) {
	vec := vector.Of("1", "1", "2", "4", "3", "1")
	r := vec.Delete(3)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("1", "1", "2", "3", "1"), vec.Collect())
	r = vec.Delete(5)
	assert.Equal(t, r, false)
}

func Test_Vector_DeleteMany(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	r := vec.Delete(3, 0, 5)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("1", "2", "4", "6"), vec.Collect())
	r = vec.Delete(5, 4)
	assert.Equal(t, r, false)
}

func Test_Vector_DeleteManyFromTail(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	r := vec.Delete(4, 5, 6)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("0", "1", "2", "3"), vec.Collect())
}

func Test_Vector_DeleteManyFromHead(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	r := vec.Delete(0, 1, 2)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("3", "4", "5", "6"), vec.Collect())
}

func Test_Vector_DeleteManyFromMiddle(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	r := vec.Delete(4, 3)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("0", "1", "2", "5", "6"), vec.Collect())
}

func Test_Vector_Set(t *testing.T) {
	vec := vector.Of("1", "1", "2", "4", "3", "1")
	added := vec.Set(10, "11")
	assert.Equal(t, added, true)
	assert.Equal(t, slice.Of("1", "1", "2", "4", "3", "1", "", "", "", "", "11"), vec.Collect())
}

func Test_Vector_DeleteByIterator(t *testing.T) {
	vec := vector.Of(1, 1, 2, 4, 3, 1)
	iter := vec.BeginEdit()

	i := 0
	var v int
	var ok, deleted bool
	for v, ok = iter.Next(); ok; v, ok = iter.Next() {
		i++
		deleted = iter.Delete()
	}

	_, _, _ = v, ok, deleted

	assert.Equal(t, 6, i)
	assert.Equal(t, 0, len(vec.Collect()))
}

func Test_Vector_DeleteByIterator_Reverse(t *testing.T) {
	vec := vector.Of(1, 1, 2, 4, 3, 1)
	iter := vec.Tail()

	i := 0
	var v int
	var ok, deleted bool
	for v, ok = iter.Prev(); ok; v, ok = iter.Prev() {
		i++
		deleted = iter.Delete()
	}

	_, _, _ = v, ok, deleted

	assert.Equal(t, 6, i)
	assert.Equal(t, 0, len(vec.Collect()))
}

func Test_Vector_FilterMapReduce(t *testing.T) {
	s := vector.Of(1, 1, 2, 4, 3, 4).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(sum.Of[int])
	assert.Equal(t, 20, s)
}

func Test_Vector_Group(t *testing.T) {
	groups := group.Of(vector.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 1, 3, 1, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}
