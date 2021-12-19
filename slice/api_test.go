package slice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConvertByAndChain(t *testing.T) {

	var toString Converter[int, string] = func(i int) string { return fmt.Sprintf("%d", i) }
	var addTail Converter[string, string] = func(s string) string { return s + "_tail" }

	converted := Convert(Of(1, 2, 3), And(toString, addTail))

	assert.Equal(t, Of("1_tail", "2_tail", "3_tail"), converted)
}

func Test_SpreadSlices(t *testing.T) {
	var (
		multiDimension [][]int
		oneDimension   []int
	)

	multiDimension = Of(Of(1, 2, 3), Of(4, 5, 6))

	oneDimension = Spread(multiDimension, AsIs[[]int])

	assert.Equal(t, Of(1, 2, 3, 4, 5, 6), oneDimension)
}
