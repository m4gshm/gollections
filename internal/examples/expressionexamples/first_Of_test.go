package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/m4gshm/gollections/expr/first"
	"github.com/m4gshm/gollections/predicate/more"
)

func Test_First(t *testing.T) {

	result, ok := first.Of(1, 3, 5, 7, 9, 11).By(more.Than(5)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 7, result)
}
