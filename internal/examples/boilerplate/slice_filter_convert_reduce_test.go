package boilerplate

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func Test_Slice_Filter_Conver_Reduce(t *testing.T) {

	data, err := slice.Conv(slice.Of("1", "2", "3", "4", "_", "6"), strconv.Atoi)
	//[1 2 3 4], invalid syntax

	even := func(i int) bool { return i%2 == 0 }
	result := slice.Reduce(slice.Convert(slice.Filter(data, even), strconv.Itoa), op.Sum) //24

	assert.ErrorIs(t, err, strconv.ErrSyntax)
	assert.Equal(t, "24", result)

}
