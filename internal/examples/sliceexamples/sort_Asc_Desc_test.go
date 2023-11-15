package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/slice/sort"
)

func Test_SortAscDesc(t *testing.T) {

	var ascendingSorted = sort.Asc([]int{1, 3, -1, 2, 0})   //[]int{-1, 0, 1, 2, 3}
	var descendingSorted = sort.Desc([]int{1, 3, -1, 2, 0}) //[]int{3, 2, 1, 0, -1}

	assert.Equal(t, []int{-1, 0, 1, 2, 3}, ascendingSorted)
	assert.Equal(t, []int{3, 2, 1, 0, -1}, descendingSorted)
}
