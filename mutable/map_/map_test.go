package map_

import (
	"sort"
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	unordered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(unordered.Collect()))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)
	for it, key, val, ok := unordered.First(); ok; key, val, ok = it.Next() {
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = unordered.Keys().Collect()
	sort.Ints(keys)
	values = unordered.Values().Collect()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_IterateOverRange(t *testing.T) {
	unordered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(unordered.Collect()))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)
	for key, val := range unordered {
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = unordered.Keys().Collect()
	sort.Ints(keys)
	values = unordered.Values().Collect()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_Add(t *testing.T) {
	d := New[int, string](4)
	s := d.SetVerify(1, "1")
	assert.Equal(t, s, true)

	d.Set(2, "2")
	d.Set(4, "4")
	d.Set(3, "3")

	s = d.SetVerify(1, "11")
	assert.Equal(t, s, false)

	keys := d.Keys().Collect()
	sort.Ints(keys)
	values := d.Values().Collect()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}
