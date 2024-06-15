package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func Test_Slice_AccumSum(t *testing.T) {

	var sum = slice.Accum(100, slice.Of(1, 2, 3, 4, 5, 6), op.Sum)
	//121

	assert.Equal(t, 121, sum)
}
