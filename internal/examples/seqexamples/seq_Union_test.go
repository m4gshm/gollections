package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_Union(t *testing.T) {

	var result []int

	seq1 := seq.Of(1, 3, 5)
	seq2 := seq.Of(7, 9, 11)
	for i := range seq.Union(seq1, seq2) {
		result = append(result, i)
	}
	//[]int{1, 3, 5, 7, 9, 11}

	assert.Equal(t, []int{1, 3, 5, 7, 9, 11}, result)
}
