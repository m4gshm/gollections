package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/predicate/exclude"
	"github.com/m4gshm/gollections/predicate/one"
	"github.com/m4gshm/gollections/slice"
)

func Test_OneOf(t *testing.T) {

	var f1 = slice.Filter([]int{1, 3, 5, 7, 9, 11}, one.Of(1, 7).Or(one.Of(11))) //[]int{1, 7, 11}
	var f2 = slice.Filter([]int{1, 3, 5, 7, 9, 11}, exclude.All(1, 7, 11))       //[]int{3, 5, 9}

	assert.Equal(t, slice.Of(1, 7, 11), f1)
	assert.Equal(t, slice.Of(3, 5, 9), f2)
}
