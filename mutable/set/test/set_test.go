package test

import (
	"testing"

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

	expected := slice.Of(1, 2, 4, 3)
	assert.Equal(t, expected, values)

	iterSlice := it.Slice[int](set.Begin())
	assert.Equal(t, expected, iterSlice)

	out := make([]int, 0)
	for it := set.Begin(); it.HasNext(); {
		n, _ := it.Get()
		out = append(out, n)
	}
	assert.Equal(t, expected, out)

	out = make([]int, 0)
	_ = set.ForEach(func(v int) { out = append(out, v) })
}

func Test_Set_Add(t *testing.T) {
	set := set.New[int](0)
	added, _ := set.Add(1, 2, 4, 3)
	assert.Equal(t, added, true)
	added, _ = set.Add(1)
	assert.Equal(t, added, false)

	values := set.Collect()

	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_Set_Delete(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := set.Collect()

	for _, v := range values {
		_, _ = set.Delete(v)
	}

	assert.Equal(t, 0, len(set.Collect()))
}

func Test_Set_DeleteByIterator(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	iter := set.Begin()

	i := 0
	for iter.HasNext() {
		i++
		_, _ = iter.Delete()
	}

	assert.Equal(t, 4, i)
	assert.Equal(t, 0, len(set.Collect()))
}

func Test_Set_FilterMapReduce(t *testing.T) {
	sum := set.Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 12, sum)

	sum = it.Pipe[int](set.Of(1, 1, 2, 4, 3, 1).Begin()).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 12, sum)
}

func Test_Set_Group(t *testing.T) {
	groups := group.Of(set.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}
