package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/slice"
)

func Test_Slice_ReduceSumOfNilSlice(t *testing.T) {

	var ints []int
	var sum = slice.Reduce(ints, func(i1, i2 int) int { return i1 + i2 })
	//0

	assert.Equal(t, 0, sum)
}
