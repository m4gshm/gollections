package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/slice"
)

func Test_Top3(t *testing.T) {

	result := slice.Top(3, []int{1, 3, 5, 7, 9, 11}) //[]int{1, 3, 5}

	assert.Equal(t, []int{1, 3, 5}, result)

}

func Test_Top10(t *testing.T) {

	result := slice.Top(10, []int{1, 3, 5, 7, 9, 11}) //[]int{1, 3, 5, 7, 9, 11}

	assert.Equal(t, []int{1, 3, 5, 7, 9, 11}, result)

}
