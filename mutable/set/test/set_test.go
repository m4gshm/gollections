package test

import (
	"sort"
	"testing"

	cgroup "github.com/m4gshm/gollections/c/group"
	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/mutable/set"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"

	"github.com/m4gshm/gollections/walk/group"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Iterate(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := set.Collect()

	assert.Equal(t, 4, len(values))

	expected := slice.Of(1, 2, 3, 4)
	sort.Ints(values)
	assert.Equal(t, expected, values)

	iterSlice := it.ToSlice(set.Begin())
	sort.Ints(iterSlice)
	assert.Equal(t, expected, iterSlice)

	out := make(map[int]int, 0)
	it := set.Begin()
	for v, ok := it.Next(); ok; v, ok = it.Next() {
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
	added = set.AddNewOne(1)
	assert.Equal(t, added, false)

	values := set.Collect()
	sort.Ints(values)

	assert.Equal(t, slice.Of(1, 2, 3, 4), values)
}

func Test_Set_Delete(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := set.Collect()

	for _, v := range values {
		set.Delete(v)
	}

	assert.Equal(t, 0, len(set.Collect()))
}

func Test_Set_DeleteByIterator(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	iter := set.BeginEdit()

	i := 0
	for _, ok := iter.Next(); ok; _, ok = iter.Next() {
		i++
		iter.Delete()
	}

	assert.Equal(t, 4, i)
	assert.Equal(t, 0, len(set.Collect()))
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

func Test_Set_Group_By_Walker(t *testing.T) {
	groups := group.Of(set.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	fg := groups[false]
	sort.Ints(fg)
	tg := groups[true]
	sort.Ints(tg)
	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, fg)
	assert.Equal(t, []int{0, 2, 4, 6}, tg)
}

func Test_Set_Group_By_Iterator(t *testing.T) {
	groups := cgroup.Of(set.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 }).Collect()

	assert.Equal(t, len(groups), 2)
	fg := groups[false]

	sort.Ints(fg)
	tg := groups[true]
	sort.Ints(tg)

	assert.Equal(t, []int{1, 3, 7}, fg)
	assert.Equal(t, []int{0, 2, 4, 6}, tg)
}
