//go:build goexperiment.rangefunc

package test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/immutable/map_"
	"github.com/m4gshm/gollections/k"
	"github.com/m4gshm/gollections/slice"
)

func Test_Map_Iterate_All(t *testing.T) {
	dict := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)

	for key, val := range dict.All {
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = dict.Keys().Slice()
	values = dict.Values().Slice()
	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}
