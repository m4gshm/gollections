package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/predicate/more"
)

func Test_First(t *testing.T) {

	result, ok := loop.First(loop.Of(1, 3, 5, 7, 9, 11), more.Than(5)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 7, result)
}
