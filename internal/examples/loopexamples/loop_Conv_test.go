package loopexamples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
)

func Test_Conv(t *testing.T) {

	result, err := loop.Conv(loop.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi).Slice()
	//[]int{1, 3, 5}, ErrSyntax

	assert.Equal(t, []int{1, 3, 5}, result)
	assert.ErrorIs(t, err, strconv.ErrSyntax)
}
