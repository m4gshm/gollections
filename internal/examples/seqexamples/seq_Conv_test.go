package seqexamples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seqe"
)

func Test_Conv(t *testing.T) {

	sequence := seq.Conv(seq.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi)
	result, err := seqe.Slice(sequence)
	//[]int{1, 3, 5}, ErrSyntax

	assert.Equal(t, []int{1, 3, 5}, result)
	assert.ErrorIs(t, err, strconv.ErrSyntax)

	var out []int
	for v, err := range sequence {
		//ignore error
		_ = err
		out = append(out, v)
	}

	assert.Equal(t, []int{1, 3, 5, 0, 9, 11}, out)
}
