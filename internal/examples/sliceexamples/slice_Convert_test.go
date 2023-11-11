package sliceexamples

import (
	"testing"

	"strconv"

	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Convert(t *testing.T) {

	var s []string = slice.Convert([]int{1, 3, 5, 7, 9, 11}, strconv.Itoa)
	//[]string{"1", "3", "5", "7", "9", "11"}

	assert.Equal(t, slice.Of("1", "3", "5", "7", "9", "11"), s)
}
