package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/seq"
)

func Test_First(t *testing.T) {

	result, ok := seq.First(seq.Of(1, 3, 5, 7, 9, 11), more.Than(5)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 7, result)

	//or
	result, ok = seq.Of(1, 3, 5, 7, 9, 11).First(more.Than(5)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 7, result)
}
