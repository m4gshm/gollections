package boilerplate

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Slice_Filter_Conver_Reduce(t *testing.T) {

	even := func(i int) bool { return i%2 == 0 }
	result := slice.Reduce(
		slice.Convert(
			slice.Filter(slice.Of(1, 2, 3, 4), even),
			strconv.Itoa,
		),
		op.Sum[string],
	)

	assert.Equal(t, "24", result)

}
