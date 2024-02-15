//go:build goexperiment.rangefunc

package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/immutable/vector"
	"github.com/m4gshm/gollections/slice"
)

func Test_VectorIterate_go_1_22(t *testing.T) {
	vec := vector.Of(1, 1, 2, 4, 3, 1)
	expected := slice.Of(1, 1, 2, 4, 3, 1)

	out := []int{}

	for _, v := range vec.All {
		out = append(out, v)
	}

	assert.Equal(t, len(expected), len(out))

	for i := range out {
		assert.Equal(t, out[i], expected[i])
	}
}
