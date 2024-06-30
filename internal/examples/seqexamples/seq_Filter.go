package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/predicate/exclude"
	"github.com/m4gshm/gollections/predicate/one"
	"github.com/m4gshm/gollections/seq"
)

func Test_OneOf(t *testing.T) {

	var f1 = seq.Slice(seq.Filter(seq.Of(1, 3, 5, 7, 9, 11), one.Of(1, 7).Or(one.Of(11))))
	//[]int{1, 7, 11}

	var f2 = seq.Slice(seq.Filter(seq.Of(1, 3, 5, 7, 9, 11), exclude.All(1, 7, 11)))
	//[]int{3, 5, 9}

	assert.Equal(t, []int{1, 7, 11}, f1)
	assert.Equal(t, []int{3, 5, 9}, f2)
}
