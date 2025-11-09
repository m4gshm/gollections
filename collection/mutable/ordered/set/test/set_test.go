package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/collection/mutable/ordered/set"
	"github.com/m4gshm/gollections/seq"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/slice"
)

func Test_Set_FromSeq(t *testing.T) {
	set := set.FromSeq(seq.Of(1, 1, 2, 2, 3, 4, 3, 2, 1))
	assert.Equal(t, slice.Of(1, 2, 3, 4), set.Slice())
}

func Test_Set_Iterate(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := set.Slice()

	assert.Len(t, values, 4)

	expected := slice.Of(1, 2, 4, 3)
	assert.Equal(t, expected, values)

	iterSlice := seq.Slice(set.All)
	assert.Equal(t, expected, iterSlice)

	out := make([]int, 0)
	for v := range set.All {
		out = append(out, v)
	}
	assert.Equal(t, expected, out)

	out = make([]int, 0)
	set.ForEach(func(v int) { out = append(out, v) })
}

func Test_Set_AddVerify(t *testing.T) {
	set := set.NewCap[int](0)
	added := set.AddNew(1, 2, 4, 3)
	assert.True(t, added)
	added = set.AddOneNew(1)
	assert.False(t, added)

	values := set.Slice()

	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_Set_AddAll(t *testing.T) {
	set := set.NewCap[int](0)
	set.AddAll(seq.Of(1, 2))
	set.AddAll(seq.Of(4, 3))

	values := set.Slice()

	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_Set_AddAllNew(t *testing.T) {
	set := set.NewCap[int](0)
	added := set.AddAllNew(seq.Of(1, 2))
	assert.True(t, added)
	//4, 3 are new
	added = set.AddAllNew(seq.Of(1, 4, 3))
	assert.True(t, added)
	added = set.AddAllNew(seq.Of(2, 4, 3))
	assert.False(t, added)

	values := set.Slice()
	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_Set_Delete(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := set.Slice()

	for _, v := range values {
		set.Delete(v)
	}

	assert.Empty(t, set.Slice())
}

func Test_Set_FilterMapReduce(t *testing.T) {
	s := set.Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Convert(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 12, s)
}

func Test_Set_Convert(t *testing.T) {
	var (
		ints     = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		strings  = seq.Slice(seq.Filter(set.Convert(ints, strconv.Itoa), func(s string) bool { return len(s) == 1 }))
		strings2 = set.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }).Slice()
	)
	assert.Equal(t, slice.Of("3", "1", "5", "6", "8", "0"), strings)
	assert.Equal(t, strings, strings2)
}

func Test_Set_Flatt(t *testing.T) {
	var (
		ints        = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		fints       = set.Flat(ints, func(i int) []int { return slice.Of(i) })
		stringsPipe = seq.Filter(seq.Convert(fints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }), func(s string) bool { return len(s) == 1 })
	)
	assert.Equal(t, slice.Of("3", "1", "5", "6", "8", "0"), stringsPipe.Slice())
}

func Test_Set_Nil(t *testing.T) {
	var set *ordered.Set[int]
	var nils []int
	set.Add(1, 2, 3)
	set.Add(nils...)
	set.AddOne(4)
	set.AddAll(set.All)

	set.Delete(1, 2, 3)
	set.Delete(nils...)
	set.DeleteOne(4)

	set.IsEmpty()
	set.Len()

	set.ForEach(nil)

	set.Slice()

	_, ok := set.Head()
	assert.False(t, ok)
}

func Test_Set_Empty_All(t *testing.T) {
	set := &ordered.Set[int]{}
	assert.Empty(t, seq.Slice(set.All))
}

func Test_Set_Empty_IAll(t *testing.T) {
	set := &ordered.Set[int]{}
	out := []int{}
	for _, v := range set.IAll {
		out = append(out, v)
	}
	assert.Empty(t, out)
}

func Test_Set_Zero(t *testing.T) {
	var mset ordered.Set[int]
	var nils []int
	mset.Add(1, 2, 3)
	assert.False(t, mset.IsEmpty())
	mset.Add(nils...)
	mset.AddOne(4)
	mset.AddAll(mset.All)

	assert.Equal(t, slice.Of(1, 2, 3, 4), mset.Slice())

	mset.Delete(1, 2, 3)
	mset.Delete(nils...)
	mset.DeleteOne(4)

	assert.True(t, mset.IsEmpty())
	assert.Equal(t, 0, mset.Len())

	mset.ForEach(nil)

	_, ok := mset.Head()
	assert.False(t, ok)
}

func Test_Set_new(t *testing.T) {
	var mset = new(ordered.Set[int])
	var nils []int
	mset.Add(1, 2, 3)
	assert.False(t, mset.IsEmpty())
	mset.Add(nils...)
	mset.AddOne(4)
	mset.AddAll(mset.All)

	s := mset.Slice()
	assert.Equal(t, slice.Of(1, 2, 3, 4), s)

	mset.Delete(1, 2, 3)
	mset.Delete(nils...)
	mset.DeleteOne(4)

	assert.True(t, mset.IsEmpty())
	assert.Equal(t, 0, mset.Len())

	mset.ForEach(nil)

	_, ok := mset.Head()
	assert.False(t, ok)
}

func Test_Set_CopyByValue(t *testing.T) {
	set := set.Of[int]()

	set.Add(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)

	s := copySet(*set)
	set2 := &s

	set2.Add(17, 18, 19, 20)

	c := set.Contains(20)
	assert.True(t, c)

	heq := set.HasAny(eq.To(20))
	assert.True(t, heq)
}

func copySet[T comparable](set ordered.Set[T]) ordered.Set[T] {
	return set
}

func Test_Set_Sort(t *testing.T) {
	ints := set.Of(3, 1, 5, 6, 8, 0, -2)
	sorted := ints.Sort(op.Compare)

	assert.Same(t, ints, sorted)

	expected := ordered.NewSet(-2, 0, 1, 3, 5, 6, 8)
	assert.Equal(t, expected, sorted)

	ssorted := ints.StableSort(op.Compare)
	assert.Equal(t, expected, ssorted)
}
