package map_

import (
	"sort"
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	ordered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Collect()))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)
	for it := ordered.Begin(); it.HasNext(); {
		key, val := it.Next()
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = ordered.Keys().Collect()
	sort.Ints(keys)
	values = ordered.Values().Collect()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_Add(t *testing.T) {
	d := New[int, string](4)
	s, _ := d.Set(1, "1")
	assert.Equal(t, s, true)
	s, _ = d.Set(2, "2")
	assert.Equal(t, s, true)
	s, _ = d.Set(4, "4")
	assert.Equal(t, s, true)
	s, _ = d.Set(3, "3")
	assert.Equal(t, s, true)
	s, _ = d.Set(1, "11")
	assert.Equal(t, s, false)

	keys := d.Keys().Collect()
	sort.Ints(keys)
	values := d.Values().Collect()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}
