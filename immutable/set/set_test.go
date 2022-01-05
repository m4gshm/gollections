package set

import (
	"testing"

	"github.com/m4gshm/container/it"
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

	iterSlice := it.Slice(set.Begin())
	assert.Equal(t, expected, iterSlice)

	out := make([]int, 0)
	for it := set.Begin(); it.HasNext(); {
		out = append(out, it.Get())
	}
	assert.Equal(t, expected, out)

	out = make([]int, 0)
	set.ForEach(func(v int) { out = append(out, v) })

	assert.Equal(t, expected, out)
}

func Test_Set_Contains(t *testing.T) {
	set := Of(1, 1, 2, 4, 3, 1)
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.True(t, set.Contains(4))
	assert.True(t, set.Contains(3))
	assert.False(t, set.Contains(0))
	assert.False(t, set.Contains(-1))
}

func Test_Set_FilterMapReduce(t *testing.T) {
	sum := Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	//no sum, already computer stream
	assert.Equal(t, 12, sum)

	sum = it.Pipe[int](Of(1, 1, 2, 4, 3, 1).Begin()).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	//no sum, already computer stream
	assert.Equal(t, 12, sum)
}

func Test_Set_Group(t *testing.T) {
	groups := walk.Group(Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}
