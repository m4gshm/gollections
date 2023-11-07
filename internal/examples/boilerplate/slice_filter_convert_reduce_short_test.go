package boilerplate

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/filter"
	"github.com/m4gshm/gollections/slice/sum"
)

func Test_Slice_FilterAndConver_Reduce(t *testing.T) {

	data := slice.Of(1, 2, 3, 4)
	even := func(i int) bool { return i%2 == 0 }

	result := sum.Of(filter.AndConvert(data, even, strconv.Itoa))

	assert.Equal(t, "24", result)

}
