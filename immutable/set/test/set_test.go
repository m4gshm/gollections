package test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/sum"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_Set_Iterate(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := set.Collect()

	assert.Equal(t, 4, len(values))

	expected := slice.Of(1, 2, 3, 4)
	sort.Ints(values)
	assert.Equal(t, expected, values)

	iterSlice := it.Slice(set.Begin())
	sort.Ints(iterSlice)
	assert.Equal(t, expected, iterSlice)

	out := make(map[int]int, 0)
	for it := set.Begin(); it.HasNext(); {
		n := it.Next()
		out[n] = n
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
	s := set.Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(sum.Of[int])
	assert.Equal(t, 12, s)

	s = it.Pipe(set.Of(1, 1, 2, 4, 3, 1).Begin()).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(sum.Of[int])
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
	groups := it.Group(set.Of(0, 1, 1, 2, 4, 3, 1, 6, 7).Begin(), func(e int) bool { return e%2 == 0 }).Collect()

	assert.Equal(t, len(groups), 2)
	fg := groups[false]

	sort.Ints(fg)
	tg := groups[true]
	sort.Ints(tg)

	assert.Equal(t, []int{1, 3, 7}, fg)
	assert.Equal(t, []int{0, 2, 4, 6}, tg)
}

func Test_Set_Sort(t *testing.T) {
	var (
		elements = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		sorted   = elements.Sort(func(e1, e2 int) bool { return e1 < e2 })
	)
	assert.Equal(t, oset.Of(-2, 0, 1, 3, 5, 6, 8), sorted)
}
