package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

import "github.com/m4gshm/gollections/op/sum"

func Test_Slice_Sum(t *testing.T) {

	var sum = sum.Of(1, 2, 3, 4, 5, 6) //21

	assert.Equal(t, 21, sum)
}
