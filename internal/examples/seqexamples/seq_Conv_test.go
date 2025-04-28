package seqexamples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_Conv(t *testing.T) {

	var result []int
	for i, err := range seq.Conv(seq.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi) {
		if err != nil {
			//ErrSyntax
			break
		}
		result = append(result, i)
	}
	//[]int{1, 3, 5}

	assert.Equal(t, []int{1, 3, 5}, result)
}
