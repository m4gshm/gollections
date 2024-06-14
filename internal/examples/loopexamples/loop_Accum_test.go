package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op/delay/string_/join"
)

func Test_Loop_Accum(t *testing.T) {

	var sum = loop.Accum("Steady", loop.Of("Ready", "Go"), join.NonEmpty(", "))
	//"Steady, Ready, Go"

	assert.Equal(t, "Steady, Ready, Go", sum)
}
