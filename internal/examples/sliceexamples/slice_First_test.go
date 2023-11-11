package sliceexamples

import (
	"testing"

	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_First(t *testing.T) {

	result, ok := slice.First([]int{1, 3, 5, 7, 9, 11}, more.Than(5)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 7, result)
}
