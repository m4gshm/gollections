package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_Range(t *testing.T) {

	var numbers []int
	for n := range seq.Range(5, -2) {
		numbers = append(numbers, n)
	}
	//[]int{5, 4, 3, 2, 1, 0, -1}

	assert.Equal(t, []int{5, 4, 3, 2, 1, 0, -1}, numbers)
}
