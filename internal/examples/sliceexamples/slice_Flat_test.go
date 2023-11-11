package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/slice"
)

func Test_Slice_Flat(t *testing.T) {

	var i []int = slice.Flat([][]int{{1, 2, 3}, {4}, {5, 6}}, as.Is[[]int])
	//[]int{1, 2, 3, 4, 5, 6}

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, i)
}
