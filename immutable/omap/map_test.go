package omap

import (
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	ordered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Collect()))

	expectedK := slice.Of(1, 2, 4, 3)
	expectedV := slice.Of("1", "2", "4", "3")

	keys := []int{}
	values := []string{}
	for it, key, val, ok := ordered.First(); ok; key, val, ok = it.Next() {
		keys = append(keys, key)
		values = append(values, val)
	}
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	assert.Equal(t, slice.Of(1, 2, 4, 3), ordered.Keys().Collect())
	assert.Equal(t, slice.Of("1", "2", "4", "3"), ordered.Values().Collect())
}

func Test_Map_Iterate_Keys(t *testing.T) {
	ordered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Collect()))

	expectedK := slice.Of(1, 2, 4, 3)

	keys := []int{}
	for it, key, ok := ordered.K().First(); ok; key, ok = it.Next() {
		keys = append(keys, key)
	}
	assert.Equal(t, expectedK, keys)

}

func Test_Map_Iterate_Values(t *testing.T) {
	ordered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Collect()))

	expectedV := slice.Of("1", "2", "4", "3")

	values := []string{}
	for it, val, ok := ordered.V().First(); ok; val, ok = it.Next() {
		values = append(values, val)
	}

	assert.Equal(t, expectedV, values)
}
