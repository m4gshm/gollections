package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
)

func Test_Loop_ReduceOKSum(t *testing.T) {

	adder := func(i1, i2 int) int { return i1 + i2 }

	sum, ok := loop.ReduceOK(loop.Of(1, 2, 3, 4, 5, 6), adder)
	//21, true

	emptyLoop := loop.Of[int]()
	sum, ok = loop.ReduceOK(emptyLoop, adder)
	//0, false

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
}
