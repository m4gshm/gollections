package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/seq"
)

func Test_Flat(t *testing.T) {

	var i []int = seq.Slice(seq.Flat(seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...), as.Is))
	//[]int{1, 2, 3, 4, 5, 6}

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, i)
}
