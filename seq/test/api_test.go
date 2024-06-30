package test

import (
	"testing"

	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/seq"
	"github.com/stretchr/testify/assert"
)

func Test_ReduceSum(t *testing.T) {

	sum, ok := seq.ReduceOK(seq.Of(1, 2, 3, 4, 5, 6), op.Sum)

	assert.True(t, ok)
	assert.Equal(t, 21, sum)
}

func Test_ReduceContains(t *testing.T) {

	ok := seq.Contains(seq.Of(1, 2, 3, 4, 5, 6), 5)

	assert.True(t, ok)
}

func Test_ReduceFirst(t *testing.T) {

	result, ok := seq.First(seq.Of(1, 2, 3, 4, 5, 6), more.Than(5)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 6, result)

}
