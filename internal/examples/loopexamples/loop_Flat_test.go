package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/loop"
)

func Test_Flat(t *testing.T) {

	var i []int = loop.Flat(loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...), as.Is).Slice()
	//[]int{1, 2, 3, 4, 5, 6}

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, i)
}
