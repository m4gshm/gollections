package loopexamples

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/loop"
	"github.com/stretchr/testify/assert"
)

func Test_Convert(t *testing.T) {

	var s []string = loop.Convert(loop.Of(1, 3, 5, 7, 9, 11), strconv.Itoa).Slice()
	//[]string{"1", "3", "5", "7", "9", "11"}

	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, s)
}
