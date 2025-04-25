package seqexamples

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_Usage_Seq(t *testing.T) {

	even := func(i int) bool { return i%2 == 0 }
	strSeq := seq.Convert(seq.Filter(seq.Of(1, 2, 3, 4), even), strconv.Itoa)

	// iterate over sequence
	for s := range strSeq {
		fmt.Println(s)
	}

	// or reduce
	var oneString string = seq.Sum(strSeq) // 24

	// or collect
	var strings []string = seq.Slice(strSeq) //[2 4]

	assert.Equal(t, "24", oneString)
	assert.Equal(t, []string{"2", "4"}, strings)

}
