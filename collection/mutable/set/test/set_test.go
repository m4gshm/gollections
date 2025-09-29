package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/collection/mutable/set"
	"github.com/m4gshm/gollections/seq"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
)

func Test_Set_FromSeq(t *testing.T) {
	set := set.FromSeq(seq.Of(1, 1, 2, 2, 3, 4, 3, 2, 1))
	assert.Equal(t, slice.Of(1, 2, 3, 4), sort.Asc(set.Slice()))
}

func Test_Set_Iterate(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := sort.Asc(set.Slice())

	assert.Equal(t, 4, len(values))

	expected := slice.Of(1, 2, 3, 4)
	assert.Equal(t, expected, values)

	loopSlice := sort.Asc(seq.Slice(set.All))
	assert.Equal(t, expected, loopSlice)

	out := make(map[int]int, 0)
	for v := range set.All {
		out[v] = v
	}

	assert.Equal(t, len(expected), len(out))
	for k := range out {
		assert.True(t, set.Contains(k))
	}

	out = make(map[int]int, 0)
	set.ForEach(func(n int) { out[n] = n })

	assert.Equal(t, len(expected), len(out))
	for k := range out {
		assert.True(t, set.Contains(k))
	}
}

func Test_Set_AddVerify(t *testing.T) {
	set := set.NewCap[int](0)
	added := set.AddNew(1, 2, 4, 3)
	assert.Equal(t, added, true)
	added = set.AddOneNew(1)
	assert.Equal(t, added, false)

	values := sort.Asc(set.Slice())

	assert.Equal(t, slice.Of(1, 2, 3, 4), values)
}

func Test_Set_AddAll(t *testing.T) {
	set := set.NewCap[int](0)
	set.AddAll(seq.Of(1, 2))
	set.AddAll(seq.Of(4, 3))

	values := sort.Asc(set.Slice())

	assert.Equal(t, slice.Of(1, 2, 3, 4), values)
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

	values := sort.Asc(set.Slice())
	assert.Equal(t, slice.Of(1, 2, 3, 4), values)
}

func Test_Set_Delete(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := set.Slice()

	for _, v := range values {
		set.Delete(v)
	}

	assert.Equal(t, 0, len(set.Slice()))
}

func Test_Set_Contains(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.True(t, set.Contains(4))
	assert.True(t, set.Contains(3))
	assert.False(t, set.Contains(0))
	assert.False(t, set.Contains(-1))
}

func Test_Set_FilterMapReduce(t *testing.T) {
	s := set.Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Convert(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 12, s)
}

func Test_Set_Convert(t *testing.T) {
	var (
		ints     = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		strings  = sort.Asc(seq.Slice(seq.Filter(set.Convert(ints, strconv.Itoa), func(s string) bool { return len(s) == 1 })))
		strings2 = sort.Asc(set.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }).Slice())
	)
	assert.Equal(t, slice.Of("0", "1", "3", "5", "6", "8"), strings)
	assert.Equal(t, strings, strings2)
}

func Test_Set_Flatt(t *testing.T) {
	var (
		ints        = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		fints       = set.Flat(ints, func(i int) []int { return slice.Of(i) })
		stringsPipe = seq.Convert(fints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 })
	)
	assert.Equal(t, slice.Of("0", "1", "3", "5", "6", "8"), sort.Asc(stringsPipe.Slice()))
}

func Test_Set_Nil(t *testing.T) {
	var set *mutable.Set[int]
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

	_ = set.For(nil)
	set.ForEach(nil)

	set.Slice()

	_, ok := set.Head()
	assert.False(t, ok)
}

func Test_Set_Zero(t *testing.T) {
	var mset mutable.Set[int]
	var nils []int
	mset.Add(1, 2, 3)
	assert.False(t, mset.IsEmpty())
	mset.Add(nils...)
	mset.AddOne(4)
	mset.AddAll(mset.All)

	assert.Equal(t, *set.Of(1, 2, 3, 4), mset)
	assert.Equal(t, slice.Of(1, 2, 3, 4), sort.Asc(mset.Slice()))

	mset.Delete(1, 2, 3)
	mset.Delete(nils...)
	mset.DeleteOne(4)

	assert.True(t, mset.IsEmpty())
	assert.Equal(t, 0, mset.Len())

	mset.For(nil)
	mset.ForEach(nil)

	_, ok := mset.Head()
	assert.False(t, ok)
}

func Test_Set_new(t *testing.T) {
	var mset = new(mutable.Set[int])
	var nils []int
	mset.Add(1, 2, 3)
	assert.False(t, mset.IsEmpty())
	mset.Add(nils...)
	mset.AddOne(4)
	mset.AddAll(mset.All)

	assert.Equal(t, set.Of(1, 2, 3, 4), mset)
	assert.Equal(t, slice.Of(1, 2, 3, 4), sort.Asc(mset.Slice()))

	mset.Delete(1, 2, 3)
	mset.Delete(nils...)
	mset.DeleteOne(4)

	assert.True(t, mset.IsEmpty())
	assert.Equal(t, 0, mset.Len())

	mset.For(nil)
	mset.ForEach(nil)

	_, ok := mset.Head()
	assert.False(t, ok)
}

func Test_Set_Sort(t *testing.T) {
	ints := set.Of(3, 1, 5, 6, 8, 0, -2)
	sorted := ints.Sort(op.Compare)
	ssorted := ints.StableSort(op.Compare)
	assert.NotSame(t, ints, sorted)

	assert.Equal(t, ordered.NewSet(-2, 0, 1, 3, 5, 6, 8), sorted)
	assert.Equal(t, sorted, ssorted)
}
