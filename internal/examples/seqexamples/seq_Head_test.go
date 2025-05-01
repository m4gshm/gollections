package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_Head(t *testing.T) {

	result, ok := seq.Head(seq.Of(1, 3, 5, 7, 9, 11)) //1, true

	assert.True(t, ok)
	assert.Equal(t, 1, result)
}
