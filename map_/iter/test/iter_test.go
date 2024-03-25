package test

import (
	"testing"

	"github.com/m4gshm/gollections/map_"
	"github.com/stretchr/testify/assert"
)

func Test_Key_Zero_Safety(t *testing.T) {
	var it map_.KeyIter[int, string]

	_, ok := it.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, it.Size())

}

func Test_OrderedMapIter_Safety(t *testing.T) {
	var it map_.Iter[int, string]

	_, _, ok := it.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, it.Size())
}
