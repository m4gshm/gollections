package sliceexamples

import (
	"testing"

	"github.com/m4gshm/gollections/slice/range_"
	"github.com/stretchr/testify/assert"
)

func Test_RangeOf(t *testing.T) {

	var increasing = range_.Of(-1, 3)    //[]int{-1, 0, 1, 2}
	var decreasing = range_.Of('e', 'a') //[]rune{'e', 'd', 'c', 'b'}
	var nothing = range_.Of(1, 1)        //nil

	assert.Equal(t, []int{-1, 0, 1, 2}, increasing)
	assert.Equal(t, []rune{'e', 'd', 'c', 'b'}, decreasing)
	assert.Nil(t, nothing)
}
