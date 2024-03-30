package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
)

func Test_LoopOf_LoopS(t *testing.T) {

	var (
		ints    = loop.Of(1, 2, 3)
		strings = loop.S([]string{"a", "b", "c"})
	)

	assert.Equal(t, []int{1, 2, 3}, ints.Slice())
	assert.Equal(t, []string{"a", "b", "c"}, strings.Slice())
}
