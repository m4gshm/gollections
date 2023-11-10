package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

import "github.com/m4gshm/gollections/slice"

func Test_SliceOf(t *testing.T) {

	var s = slice.Of(1, 3, -1, 2, 0) //[]int{1, 3, -1, 2, 0}

	assert.Equal(t, []int{1, 3, -1, 2, 0}, s)
}
