package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop/range_"
)

func Test_RangeClosed(t *testing.T) {

	var increasing = range_.Closed(-1, 3).Slice()    //[]int{-1, 0, 1, 2, 3}
	var decreasing = range_.Closed('e', 'a').Slice() //[]rune{'e', 'd', 'c', 'b', 'a'}
	var one = range_.Closed(1, 1).Slice()            //[]int{1}

	assert.Equal(t, []int{-1, 0, 1, 2, 3}, increasing)
	assert.Equal(t, []rune{'e', 'd', 'c', 'b', 'a'}, decreasing)
	assert.Equal(t, []int{1}, one)
}
