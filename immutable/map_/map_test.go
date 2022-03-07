package map_

import (
	"sort"
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	dict := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))

	assert.Equal(t, 4, dict.Len())
	assert.Equal(t, 4, len(dict.Collect()))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)

	it := dict.Head()
	for key, val, ok := it.Next(); ok; key, val, ok = it.Next() {
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = dict.Keys().Collect()
	values = dict.Values().Collect()
	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_Iterate_Keys(t *testing.T) {
	dict := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(dict.Collect()))

	expectedK := slice.Of(1, 2, 3, 4)

	keys := []int{}
	for it, key, ok := dict.K().First(); ok; key, ok = it.Next() {
		keys = append(keys, key)
	}

	sort.Ints(keys)
	assert.Equal(t, expectedK, keys)

}

func Test_Map_Iterate_Values(t *testing.T) {
	ordered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Collect()))

	expectedV := slice.Of("1", "2", "3", "4")

	values := []string{}
	for it, val, ok := ordered.V().First(); ok; val, ok = it.Next() {
		values = append(values, val)
	}

	sort.Strings(values)
	assert.Equal(t, expectedV, values)
}
