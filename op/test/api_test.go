package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/op"
)

func Test_Min(t *testing.T) {
	assert.Equal(t, 5, op.Min(5, 5))
	assert.Equal(t, 5, op.Min(5, 6))
	assert.Equal(t, "A", op.Min("a", "A"))
}

func Test_Max(t *testing.T) {
	assert.Equal(t, 5, op.Max(5, 5))
	assert.Equal(t, 6, op.Max(5, 6))
	assert.Equal(t, "a", op.Max("a", "A"))
}

func Test_IfElse(t *testing.T) {
	assert.Equal(t, 5, op.IfElse(true, 5, 6))
	assert.Equal(t, 6, op.IfElse(false, 5, 6))
}

func Test_IfElseDelay(t *testing.T) {
	assert.Equal(t, 5, op.IfElse(true, func() int { return 5 }, func() int { return 6 })())
	assert.Equal(t, 6, op.IfElse(false, func() int { return 5 }, func() int { return 6 })())
}

func Test_IfDoElse(t *testing.T) {
	assert.Equal(t, 5, op.IfGetElseGet(true, func() int { return 5 }, func() int { return 6 }))
	assert.Equal(t, 6, op.IfGetElseGet(false, func() int { return 5 }, func() int { return 6 }))
}
