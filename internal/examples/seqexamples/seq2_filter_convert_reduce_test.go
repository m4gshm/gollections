package seqexamples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
)

func Test_Usage_Seq2_Errorable(t *testing.T) {

	intSeq := seq.Conv(seq.Of("1", "2", "3", "ddd4", "5"), strconv.Atoi)
	ints, err := seq2.Slice(intSeq) //[1 2 3], invalid syntax

	assert.Equal(t, []int{1, 2, 3}, ints)
	assert.ErrorContains(t, err, "invalid syntax")

}
