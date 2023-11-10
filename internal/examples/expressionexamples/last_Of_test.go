package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/m4gshm/gollections/expr/last"
	"github.com/m4gshm/gollections/predicate/less"
)

func Test_Last(t *testing.T) {

	result, ok := last.Of(1, 3, 5, 7, 9, 11).By(less.Than(9)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 7, result)
}