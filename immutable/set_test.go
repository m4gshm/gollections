package immutable

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

	iterSlice := iter.ToSlice(set.Begin())
	assert.Equal(t, expected, iterSlice)

	out := make([]int, 0)
	for it := set.Begin(); it.HasNext(); {
		out = append(out, it.Get())
	}
	assert.Equal(t, expected, out)
}
