package seqexamples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
)

func Test_Conv(t *testing.T) {

	result, err := seq2.Slice(seq.Conv(seq.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi))
	//[]int{1, 3, 5}, ErrSyntax

	assert.Equal(t, []int{1, 3, 5}, result)
	assert.ErrorIs(t, err, strconv.ErrSyntax)
}
