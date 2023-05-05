package test

import (
	"testing"

	"github.com/m4gshm/gollections/op/use"
	"github.com/stretchr/testify/assert"
)

func Test_Use(t *testing.T) {
	result := use.If(1, false).Else(2)
	assert.Equal(t, 2, result)

	result = use.If(1, true).Else(2)
	assert.Equal(t, 1, result)
}
