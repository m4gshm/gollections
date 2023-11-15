package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/vector"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_Vector_From(t *testing.T) {
	set := vector.From(iter.Of(1, 1, 2, 2, 3, 4, 3, 2, 1).Next)
	assert.Equal(t, slice.Of(1, 1, 2, 2, 3, 4, 3, 2, 1), set.Slice())
}

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
	vec.AddAll(vec)

	vec.Delete(1, 2, 3)
	vec.Delete(nils...)
	vec.DeleteOne(4)

	vec.IsEmpty()
	vec.Len()

	vec.For(nil)
	vec.ForEach(nil)
	vec.Track(nil)
	vec.TrackEach(nil)

	assert.Equal(t, nils, vec.Slice())

	head := vec.Head()
	assert.False(t, head.HasNext())
	assert.False(t, head.HasPrev())

	_, ok := head.Get()
	assert.False(t, ok)
	_, ok = head.Next()
	assert.False(t, ok)
	head.Cap()

	tail := vec.Tail()
	assert.False(t, tail.HasNext())
	assert.False(t, tail.HasPrev())

	_, ok = tail.Get()
	assert.False(t, ok)
	_, ok = tail.Next()
	assert.False(t, ok)
	tail.Cap()
}

func Test_Vector_Zero(t *testing.T) {
	var vec mutable.Vector[string]

	var nilValues []string
	vec.Add("a", "b", "c")
	vec.Add(nilValues...)
	vec.AddOne("d")
	vec.AddAll(&vec)

	vec.Delete(0, 1, 2)
	var nilIndexes []int
	vec.Delete(nilIndexes...)
	vec.DeleteOne(0)

	e := vec.IsEmpty()
	assert.False(t, e)

	l := vec.Len()
	assert.Equal(t, 4, l)

	vec.For(func(s string) error { return nil })
	vec.ForEach(func(s string) {})
	vec.Track(func(i int, s string) error { return nil })
	vec.TrackEach(func(i int, s string) {})

	assert.Equal(t, slice.Of("a", "b", "c", "d"), vec.Slice())

	head := vec.Head()
	assert.True(t, head.HasNext())
	assert.False(t, head.HasPrev())

	_, ok := head.Get()
	assert.False(t, ok)
	fv, ok := head.Next()
	assert.True(t, ok)
	assert.Equal(t, "a", fv)
	c := head.Cap()
	assert.Equal(t, 4, c)

	tail := vec.Tail()
	assert.False(t, tail.HasNext())
	assert.True(t, tail.HasPrev())

	_, ok = tail.Get()
	assert.False(t, ok)
	tv, ok := tail.Prev()
	assert.True(t, ok)
	assert.Equal(t, "d", tv)
	c = tail.Cap()
	assert.Equal(t, 4, c)
}

func Test_Vector_new(t *testing.T) {
	var vec = new(mutable.Vector[string])

	var nilValues []string
	vec.Add("a", "b", "c")
	vec.Add(nilValues...)
	vec.AddOne("d")
	vec.AddAll(vec)

	vec.Delete(0, 1, 2)
	var nilIndexes []int
	vec.Delete(nilIndexes...)
	vec.DeleteOne(0)

	e := vec.IsEmpty()
	assert.False(t, e)

	l := vec.Len()
	assert.Equal(t, 4, l)

	vec.For(func(s string) error { return nil })
	vec.ForEach(func(s string) {})
	vec.Track(func(i int, s string) error { return nil })
	vec.TrackEach(func(i int, s string) {})

	assert.Equal(t, slice.Of("a", "b", "c", "d"), vec.Slice())

	head := vec.Head()
	assert.True(t, head.HasNext())
	assert.False(t, head.HasPrev())

	_, ok := head.Get()
	assert.False(t, ok)
	fv, ok := head.Next()
	assert.True(t, ok)
	assert.Equal(t, "a", fv)
	c := head.Cap()
	assert.Equal(t, 4, c)

	tail := vec.Tail()
	assert.False(t, tail.HasNext())
	assert.True(t, tail.HasPrev())

	_, ok = tail.Get()
	assert.False(t, ok)
	tv, ok := tail.Prev()
	assert.True(t, ok)
	assert.Equal(t, "d", tv)
	c = tail.Cap()
	assert.Equal(t, 4, c)
}

func Test_Vector_AddAllOfSelf(t *testing.T) {
	vec := vector.Of(1, 2, 3)
	vec.AddAll(vec)
	assert.Equal(t, vector.Of(1, 2, 3, 1, 2, 3), vec)
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }

func Test_Vector_AddAndDelete(t *testing.T) {
	vec := vector.NewCap[int](0)
	vec.Add(range_.Closed(0, 1000)...)
	deleted := false
	for i := vec.Head(); i.HasNext(); {
		deleted = i.DeleteNext()
	}
	assert.Equal(t, deleted, true)
	assert.True(t, vec.IsEmpty())

	vec.Add(range_.Closed(0, 10000)...)
	for i := vec.Tail(); i.HasPrev(); {
		deleted = i.DeletePrev()
	}
	assert.Equal(t, deleted, true)
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
	vec.AddAll(vector.Of(1, 1, 2, 4, 3, 1))
	assert.Equal(t, slice.Of(1, 1, 2, 4, 3, 1), vec.Slice())
	vec.AddAll(vector.Of(1))
	assert.Equal(t, slice.Of(1, 1, 2, 4, 3, 1, 1), vec.Slice())
}

func Test_Vector_Add_And_Iterate(t *testing.T) {
	vec := vector.NewCap[int](0)
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

	it.Delete()

	//no the current one
	_, ok = it.Get()
	assert.False(t, ok)

	//no more
	assert.False(t, it.HasNext())
	//no prev
	assert.False(t, it.HasPrev())

	assert.True(t, vec.IsEmpty())

	// add values to vector
	vec.Add(1, 3)

	//it must to point before the first
	_, ok = it.Get()
	assert.False(t, ok)
	assert.True(t, it.HasNext())
	assert.False(t, it.HasPrev())

	v, ok = it.Next()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	assert.True(t, it.HasNext())
	assert.False(t, it.HasPrev())

	v, ok = it.Next()
	assert.True(t, ok)
	assert.Equal(t, 3, v)

	assert.False(t, it.HasNext())
	assert.True(t, it.HasPrev())

	v, ok = it.Prev()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	assert.True(t, it.HasNext())
	assert.False(t, it.HasPrev())

	//delete the first one
	it.Delete()

	//second must remains
	assert.False(t, it.HasNext())
	assert.False(t, it.HasPrev())
	v, ok = it.Get()
	assert.True(t, ok)
	assert.Equal(t, 3, v)

	assert.Equal(t, []int{3}, vec.Slice())
}

func Test_Vector_DeleteOne(t *testing.T) {
	vec := vector.Of("1", "1", "2", "4", "3", "1")
	vec.DeleteOne(3)
	assert.Equal(t, slice.Of("1", "1", "2", "3", "1"), vec.Slice())
	r := vec.DeleteActualOne(5)
	assert.Equal(t, r, false)
}

func Test_Vector_DeleteMany(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	vec.Delete(3, 0, 5)
	assert.Equal(t, slice.Of("1", "2", "4", "6"), vec.Slice())
	r := vec.DeleteActual(5, 4)
	assert.Equal(t, r, false)
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

func Test_Vector_DeleteByIterator(t *testing.T) {
	vec := vector.Of(1, 1, 2, 4, 3, 1)
	iterator := vec.IterEdit()

	i := 0
	var v int
	var ok bool
	for v, ok = iterator.Next(); ok; v, ok = iterator.Next() {
		i++
		iterator.Delete()
	}

	_, _ = v, ok

	assert.Equal(t, 6, i)
	assert.Equal(t, 0, len(vec.Slice()))
}

func Test_Vector_DeleteByIterator_Reverse(t *testing.T) {
	vec := vector.Of(1, 1, 2, 4, 3, 1)
	iterator := vec.Tail()

	i := 0
	var v int
	var ok bool
	for v, ok = iterator.Prev(); ok; v, ok = iterator.Prev() {
		i++
		iterator.Delete()
	}

	_, _ = v, ok

	assert.Equal(t, 6, i)
	assert.Equal(t, 0, len(vec.Slice()))
}

func Test_Vector_FilterMapReduce(t *testing.T) {
	s := vector.Of(1, 1, 2, 4, 3, 4).Filter(func(i int) bool { return i%2 == 0 }).Convert(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 20, s)
}

func Test_Vector_Group(t *testing.T) {
	groups := group.Of(vector.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 1, 3, 1, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}
