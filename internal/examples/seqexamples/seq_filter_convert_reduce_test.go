package seqexamples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_Usage_Seq(t *testing.T) {

	even := func(i int) bool { return i%2 == 0 }
	sequence := seq.Convert(seq.Filter(seq.Of(1, 2, 3, 4), even), strconv.Itoa)
	var result []string = seq.ToSlice(sequence) //[2 4]

	assert.Equal(t, []string{"2", "4"}, result)

}
