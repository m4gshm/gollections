package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
)

func Test_Loop_Accum(t *testing.T) {

	var sum = loop.Accum(100, loop.Of(1, 2, 3, 4, 5, 6), op.Sum)
	//121

	assert.Equal(t, 121, sum)
}
