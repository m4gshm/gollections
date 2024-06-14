package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/op/delay/string_/join"
	"github.com/m4gshm/gollections/slice"
)

func Test_Slice_AccumSum(t *testing.T) {

	var sum = slice.Accum("Steady", slice.Of("Ready", "Go"), join.NonEmpty(", "))
	//"Steady, Ready, Go"

	assert.Equal(t, "Steady, Ready, Go", sum)
}
