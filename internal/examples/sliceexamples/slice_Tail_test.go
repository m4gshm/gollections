package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/slice"
)

func Test_Tail(t *testing.T) {

	result, ok := slice.Tail([]int{1, 3, 5, 7, 9, 11}) //11, true

	assert.True(t, ok)
	assert.Equal(t, 11, result)
}
