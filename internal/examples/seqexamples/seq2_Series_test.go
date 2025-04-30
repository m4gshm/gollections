package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq2"
)

func Test_Series(t *testing.T) {

	var numbers, factorials []int
	for i, n := range seq2.Series(1, func(i int, prev int) (int, bool) {
		return i * prev, i <= 5
	}) {
		numbers = append(numbers, i)
		factorials = append(factorials, n)
	}
	//[]int{0, 1, 2, 3, 4, 5}
	//[]int{1, 1, 2, 6, 24, 120}

	assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, numbers)
	assert.Equal(t, []int{1, 1, 2, 6, 24, 120}, factorials)
}
