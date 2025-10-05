package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/break/op"
)

func Test_Min(t *testing.T) {
	r, err := op.Min(5, 5)
	assert.Equal(t, 5, r)
	assert.NoError(t, err)

	r, err = op.Min(5, 6)
	assert.Equal(t, 5, r)
	assert.NoError(t, err)

	r2, err := op.Min("a", "A")
	assert.Equal(t, "A", r2)
	assert.NoError(t, err)
}

func Test_Max(t *testing.T) {
	r, err := op.Max(5, 5)
	assert.Equal(t, 5, r)
	assert.NoError(t, err)

	r, err = op.Max(5, 6)
	assert.Equal(t, 6, r)
	assert.NoError(t, err)

	r2, err := op.Max("a", "A")
	assert.Equal(t, "a", r2)
	assert.NoError(t, err)
}

func Test_IfElse(t *testing.T) {
	r, err := op.IfElse(true, 5, 6)
	assert.Equal(t, 5, r)
	assert.NoError(t, err)

	r, err = op.IfElse(false, 5, 6)
	assert.Equal(t, 6, r)
	assert.NoError(t, err)
}

func Test_IfElseDelay(t *testing.T) {
	r, err := op.IfElse(true, func() int { return 5 }, func() int { return 6 })
	assert.Equal(t, 5, r())
	assert.NoError(t, err)
	r, err = op.IfElse(false, func() int { return 5 }, func() int { return 6 })
	assert.Equal(t, 6, r())
	assert.NoError(t, err)
}

func Test_IfGetElseGet(t *testing.T) {
	r, err := op.IfGetElseGet(true, func() (int, error) { return 5, nil }, func() (int, error) { return 6, nil })
	assert.Equal(t, 5, r)
	assert.NoError(t, err)
	r, err = op.IfGetElseGet(false, func() (int, error) { return 5, nil }, func() (int, error) { return 6, nil })
	assert.Equal(t, 6, r)
	assert.NoError(t, err)
}
