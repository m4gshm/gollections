package test

import (
	"testing"

	omap "github.com/m4gshm/gollections/immutable/ordered/map_"
	"github.com/m4gshm/gollections/map_"
	"github.com/stretchr/testify/assert"
)

func Test_Key_Zero_Safety(t *testing.T) {
	var it map_.KeyIter[int, string]

	_, ok := it.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, it.Cap())

}

func Test_OrderedMapIter_Safety(t *testing.T) {
	var it omap.Iter[int, string]

	_, _, ok := it.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, it.Cap())
}
