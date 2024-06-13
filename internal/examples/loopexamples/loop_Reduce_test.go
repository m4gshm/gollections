package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
)

func Test_Loop_ReduceSum(t *testing.T) {

	var sum, ok = loop.Reduce(loop.Of(1, 2, 3, 4, 5, 6), func(i1, i2 int) int { return i1 + i2 })
	//21, true

	assert.True(t, ok)
	assert.Equal(t, 21, sum)
}
