package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

import "github.com/m4gshm/gollections/slice/range_"

func Test_RangeClosed(t *testing.T) {

	var increasing = range_.Closed(-1, 3)    //[]int{-1, 0, 1, 2, 3}
	var decreasing = range_.Closed('e', 'a') //[]rune{'e', 'd', 'c', 'b', 'a'}
	var one = range_.Closed(1, 1)            //[]int{1}

	assert.Equal(t, []int{-1, 0, 1, 2, 3}, increasing)
	assert.Equal(t, []rune{'e', 'd', 'c', 'b', 'a'}, decreasing)
	assert.Equal(t, []int{1}, one)
}
