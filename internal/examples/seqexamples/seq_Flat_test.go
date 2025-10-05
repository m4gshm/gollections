package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/seq"
)

func Test_Flat(t *testing.T) {

	twoDimensions := [][]int{{1, 2, 3}, {4}, {5, 6}}
	var i []int = seq.Slice(seq.Flat(seq.Of(twoDimensions...), as.Is))
	//[]int{1, 2, 3, 4, 5, 6}

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, i)
}
