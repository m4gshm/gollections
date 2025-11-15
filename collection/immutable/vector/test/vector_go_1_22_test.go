//go:build goexperiment.rangefunc

package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/immutable/vector"
	"github.com/m4gshm/gollections/slice"
)

func Test_VectorIterate_All(t *testing.T) {
	vec := vector.Of(1, 1, 2, 4, 3, 1)
	expected := slice.Of(1, 1, 2, 4, 3, 1)

	out := []int{}

	for v := range vec.All {
		out = append(out, v)
	}

	assert.Len(t, out, len(expected))

	for i := range out {
		assert.Equal(t, out[i], expected[i])
	}
}
