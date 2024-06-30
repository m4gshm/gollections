package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_Loop_ReduceOKSum(t *testing.T) {

	adder := func(i1, i2 int) int { return i1 + i2 }

	sum, ok := seq.ReduceOK(seq.Of(1, 2, 3, 4, 5, 6), adder)
	//21, true

	emptyLoop := seq.Of[int]()
	sum, ok = seq.ReduceOK(emptyLoop, adder)
	//0, false

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
}
