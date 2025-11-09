package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func Test_IfGetElseGet(t *testing.T) {
	assert.Equal(t, 5, op.IfGetElseGet(true, func() int { return 5 }, func() int { return 6 }))
	assert.Equal(t, 6, op.IfGetElseGet(false, func() int { return 5 }, func() int { return 6 }))
}

func Test_SumAndSub(t *testing.T) {
	// integers
	assert.Equal(t, 5, op.Sum(2, 3))
	assert.Equal(t, -1, op.Sub(2, 3))
	// floats
	assert.InDelta(t, 5.5, op.Sum(2.5, 3.0), 1e-9)
	assert.InDelta(t, -0.5, op.Sub(2.5, 3.0), 1e-9)
	// complex numbers (only addition/subtraction defined for complex in constraints.Complex)
	c1 := complex(1, 2)
	c2 := complex(3, 4)
	assert.Equal(t, complex(4, 6), op.Sum(c1, c2))
	assert.Equal(t, complex(-2, -2), op.Sub(c1, c2))
}

func Test_IfElseErrfAndErr(t *testing.T) {
	// IfElseErrf success
	v, err := op.IfElseErrf(true, 10, "should not happen")
	require.NoError(t, err)
	assert.Equal(t, 10, v)
	// failure case with format
	_, err = op.IfElseErrf(false, 0, "%s %d", "value", 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "value 5")

	// IfElseErr success
	v2, err := op.IfElseErr(true, 20, errors.New("boom"))
	require.NoError(t, err)
	assert.Equal(t, 20, v2)
	// failure case
	_, err = op.IfElseErr(false, 0, errors.New("error occurred"))
	require.Error(t, err)
	assert.Equal(t, "error occurred", err.Error())
}

func Test_IfElseGetAndGetters(t *testing.T) {
	// IfElseGet with functions
	res := op.IfElseGet(true, 1, func() int { return 2 })
	assert.Equal(t, 1, res)
	res = op.IfElseGet(false, 1, func() int { return 2 })
	assert.Equal(t, 2, res)

	// IfGetElse
	res = op.IfGetElse(true, func() int { return 3 }, 4)
	assert.Equal(t, 3, res)
	res = op.IfGetElse(false, func() int { return 3 }, 4)
	assert.Equal(t, 4, res)

	// IfGetElseGet
	res = op.IfGetElseGet(true, func() int { return 5 }, func() int { return 6 })
	assert.Equal(t, 5, res)
	res = op.IfGetElseGet(false, func() int { return 5 }, func() int { return 6 })
	assert.Equal(t, 6, res)
}

func Test_IfElseGetWithErrAndGettersErr(t *testing.T) {
	// IfElseGetWithErr success
	v, err := op.IfElseGetWithErr(true, 7, func() (int, error) { return 8, nil })
	require.NoError(t, err)
	assert.Equal(t, 7, v)
	// failure case
	_, err = op.IfElseGetWithErr(false, 0, func() (int, error) { return 9, errors.New("fail") })
	require.Error(t, err)
	assert.Equal(t, "fail", err.Error())

	// IfGetElseGetErr success
	v2, err := op.IfGetElseGetErr(true, func() int { return 10 }, func() error { return nil })
	require.NoError(t, err)
	assert.Equal(t, 10, v2)
	// failure case
	_, err = op.IfGetElseGetErr(false, func() int { return 11 }, func() error { return errors.New("oops") })
	require.Error(t, err)
	assert.Equal(t, "oops", err.Error())
}

func Test_CompareAndGet(t *testing.T) {
	// Compare ordered types
	assert.Equal(t, -1, op.Compare(1, 2))
	assert.Equal(t, 0, op.Compare("a", "a"))
	assert.Equal(t, 1, op.Compare(5.0, 4.0))

	// Get simply returns value from getter
	val := op.Get(func() string { return "hello" })
	assert.Equal(t, "hello", val)
}

func Test_IfElseGetErr(t *testing.T) {
	// success case
	v, err := op.IfElseGetErr(true, 42, func() error { return errors.New("should not happen") })
	require.NoError(t, err)
	assert.Equal(t, 42, v)
	// failure case
	_, err = op.IfElseGetErr(false, 0, func() error { return errors.New("boom") })
	require.Error(t, err)
	assert.Equal(t, "boom", err.Error())
}
