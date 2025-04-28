package seqexamples

import (
	"testing"

	"github.com/m4gshm/gollections/seq"

	"github.com/stretchr/testify/assert"
)

func Test_Top(t *testing.T) {

	var i []int = seq.Slice(seq.Top(4, seq.Of(1, 3, 5, 7, 9, 11)))
	//[]int{1, 3, 5, 7}

	assert.Equal(t, []int{1, 3, 5, 7}, i)
}
