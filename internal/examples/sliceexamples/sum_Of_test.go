package sliceexamples

import (
	"testing"

	"github.com/m4gshm/gollections/op/sum"
	"github.com/stretchr/testify/assert"
)

func Test_Slice_Sum(t *testing.T) {

	var sum = sum.Of(1, 2, 3, 4, 5, 6) //21

	assert.Equal(t, 21, sum)
}
