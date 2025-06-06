package seqexamples

import (
	"testing"

	"github.com/m4gshm/gollections/seq"

	"github.com/stretchr/testify/assert"
)

func Test_Skip(t *testing.T) {

	var i []int = seq.Slice(seq.Skip(4, seq.Of(1, 3, 5, 7, 9, 11)))
	//[]int{9, 11}

	assert.Equal(t, []int{9, 11}, i)
}
