package seqexamples

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/seq"
	"github.com/stretchr/testify/assert"
)

func Test_Convert(t *testing.T) {

	var s []string = seq.Slice(seq.Convert(seq.Of(1, 3, 5, 7, 9, 11), strconv.Itoa))
	//[]string{"1", "3", "5", "7", "9", "11"}

	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, s)
}
