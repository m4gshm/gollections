package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_Loop_ReduceSum(t *testing.T) {

	var sum = seq.Reduce(seq.Of(1, 2, 3, 4, 5, 6), func(i1, i2 int) int { return i1 + i2 })
	//21

	assert.Equal(t, 21, sum)
}
