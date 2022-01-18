package map_

import (
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	opdered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(opdered.Collect()))

	expectedK := slice.Of(1, 2, 4, 3)
	expectedV := slice.Of("1", "2", "4", "3")

	keys := make([]int, 0)
	values := make([]string, 0)
	for it := opdered.Begin(); it.HasNext(); {
		key, val, _ := it.Get()
		keys = append(keys, key)
		values = append(values, val)
	}
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	assert.Equal(t, slice.Of(1, 2, 4, 3), opdered.Keys().Collect())
	assert.Equal(t, slice.Of("1", "2", "4", "3"), opdered.Values().Collect())
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

	assert.Equal(t, slice.Of(1, 2, 4, 3), d.Keys().Collect())
	assert.Equal(t, slice.Of("1", "2", "4", "3"), d.Values().Collect())
}
