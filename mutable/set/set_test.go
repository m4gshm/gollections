package set

import (
	"testing"

	"github.com/m4gshm/container/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/walk"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Iterate(t *testing.T) {
	set := Of(1, 1, 2, 4, 3, 1)
	values := set.Elements()

	assert.Equal(t, 4, set.Len())
	assert.Equal(t, 4, len(values))

	expected := slice.Of(1, 2, 4, 3)
	assert.Equal(t, expected, values)

	iterSlice := iter.Slice[int](set.Begin())
	assert.Equal(t, expected, iterSlice)

	out := make([]int, 0)
	for it := set.Begin(); it.HasNext(); {
		out = append(out, it.Get())
	}
	assert.Equal(t, expected, out)

	out = make([]int, 0)
	set.ForEach(func(v int) { out = append(out, v) })
}

func Test_Set_Add(t *testing.T) {
	set := New[int](0)
	assert.Equal(t, set.Add(1), true)
	assert.Equal(t, set.Add(2), true)
	assert.Equal(t, set.Add(4), true)
	assert.Equal(t, set.Add(3), true)
	assert.Equal(t, set.Add(1), false)

	values := set.Elements()

	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_Set_Delete(t *testing.T) {
	set := Of(1, 1, 2, 4, 3, 1)
	values := set.Elements()

	for _, v := range values {
		set.Delete(v)
	}

	assert.Equal(t, 0, set.Len())
}

func Test_Set_DeleteByIterator(t *testing.T) {
	set := Of(1, 1, 2, 4, 3, 1)
	iter := set.Begin()

	i := 0
	for iter.HasNext() {
		i++
		iter.Delete()
	}

	assert.Equal(t, 4, i)
	assert.Equal(t, 0, set.Len())
}

func Test_Set_FilterMapReduce(t *testing.T) {
	sum := Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	//no sum, already computer stream
	assert.Equal(t, 12, sum)

	sum = iter.Pipe[int](Of(1, 1, 2, 4, 3, 1).Begin()).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	//no sum, already computer stream
	assert.Equal(t, 12, sum)
}

func Test_Set_Group(t *testing.T) {
	groups := walk.Group(Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}