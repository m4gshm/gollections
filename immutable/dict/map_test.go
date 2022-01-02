package omap

import (
	"testing"

	"github.com/m4gshm/container/K"
	"github.com/m4gshm/container/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	opdered := Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, opdered.Len())
	assert.Equal(t, 4, len(opdered.Values()))

	expectedK := slice.Of(1, 2, 4, 3)
	expectedV := slice.Of("1", "2", "4", "3")

	keys := make([]int, 0)
	values := make([]string, 0)
	for it := opdered.Begin(); it.HasNext(); {
		e := it.Get()
		keys = append(keys, e.Key())
		values = append(values, e.Value())
	}
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)
}
