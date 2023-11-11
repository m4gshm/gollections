package sliceexamples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/slice"
)

func Test_Conv(t *testing.T) {

	result, err := slice.Conv(slice.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi)
	//[]int{1, 3, 5}, ErrSyntax

	assert.Equal(t, slice.Of(1, 3, 5), result)
	assert.ErrorIs(t, err, strconv.ErrSyntax)
}
