package test

import (
	"testing"

	oMapIter "github.com/m4gshm/gollections/immutable/ordered/map_/iter"
	"github.com/m4gshm/gollections/map_/iter"
	"github.com/stretchr/testify/assert"
)

func Test_Key_Zero_Safety(t *testing.T) {
	var it iter.KeyIter[int, string]

	_, ok := it.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, it.Cap())

}

func Test_OrderedEmbedMapKVIter_Safety(t *testing.T) {
	var it oMapIter.OrderedMapIter[int, string]

	_, _, ok := it.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, it.Cap())
}
