package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/sum"
)

func Test_Sum(t *testing.T) {

	var sum = sum.Of(loop.Of(1, 2, 3, 4, 5, 6)) //21

	assert.Equal(t, 21, sum)
}
