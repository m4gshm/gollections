//go:build goexperiment.rangefunc

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/slice"
)

func Test_Set_Iterate_go_1_22(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	expected := slice.Of(1, 2, 3, 4)

	out := make(map[int]int, 0)

	for v := range set.All {
		out[v] = v
	}

	assert.Equal(t, len(expected), len(out))
	for k := range out {
		assert.True(t, set.Contains(k))
	}
}
