package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/predicate/exclude"
	"github.com/m4gshm/gollections/predicate/one"
)

func Test_OneOf(t *testing.T) {

	var f1 = loop.Filter(loop.Of(1, 3, 5, 7, 9, 11), one.Of(1, 7).Or(one.Of(11))).Slice()
	//[]int{1, 7, 11}

	var f2 = loop.Filter(loop.Of(1, 3, 5, 7, 9, 11), exclude.All(1, 7, 11)).Slice()
	//[]int{3, 5, 9}

	assert.Equal(t, []int{1, 7, 11}, f1)
	assert.Equal(t, []int{3, 5, 9}, f2)
}
