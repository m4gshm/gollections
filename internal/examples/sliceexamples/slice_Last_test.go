package sliceexamples

import (
	"testing"

	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Last(t *testing.T) {

	result, ok := slice.Last([]int{1, 3, 5, 7, 9, 11}, less.Than(9)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 7, result)
}
