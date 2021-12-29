package mutable

import (
	"testing"

	"github.com/m4gshm/container/iter"
	"github.com/m4gshm/container/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Iterate(t *testing.T) {
	set := NewOrderedSet(1, 1, 2, 4, 3, 1)
	values := set.Values()

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
	set := NewOrderedSet[int]()
	assert.Equal(t, set.Add(1), true)
	assert.Equal(t, set.Add(2), true)
	assert.Equal(t, set.Add(4), true)
	assert.Equal(t, set.Add(3), true)
	assert.Equal(t, set.Add(1), false)

	values := set.Values()

	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_Set_Delete(t *testing.T) {
	set := NewOrderedSet(1, 1, 2, 4, 3, 1)
	values := set.Values()

	for _, v := range values {
		set.Delete(v)
	}

	assert.Equal(t, 0, set.Len())
}

func Test_Set_DeleteByIterator(t *testing.T) {
	set := NewOrderedSet(1, 1, 2, 4, 3, 1)
	iter := set.Begin()

	i := 0
	for iter.HasNext() {
		i++
		iter.Delete()
	}

	assert.Equal(t, 4, i)
	assert.Equal(t, 0, set.Len())
}
