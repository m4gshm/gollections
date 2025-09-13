package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/op"
)

// Test Sum function with different numeric types
func Test_Sum(t *testing.T) {
	// Test with integers
	assert.Equal(t, 8, op.Sum(3, 5))
	assert.Equal(t, 0, op.Sum(-5, 5))
	assert.Equal(t, -10, op.Sum(-3, -7))

	// Test with floats
	assert.Equal(t, 3.7, op.Sum(1.2, 2.5))
	assert.Equal(t, 0.0, op.Sum(-1.5, 1.5))

	// Test with strings (concatenation)
	assert.Equal(t, "hello world", op.Sum("hello ", "world"))
	assert.Equal(t, "test", op.Sum("", "test"))
	assert.Equal(t, "test", op.Sum("test", ""))
}

// Test Sub function with different numeric types
func Test_Sub(t *testing.T) {
	// Test with integers
	assert.Equal(t, 2, op.Sub(5, 3))
	assert.Equal(t, -10, op.Sub(-5, 5))
	assert.Equal(t, 4, op.Sub(-3, -7))

	// Test with floats
	assert.Equal(t, -1.3, op.Sub(1.2, 2.5))
	assert.Equal(t, -3.0, op.Sub(-1.5, 1.5))
}

func Test_Min(t *testing.T) {
	// Test with integers
	assert.Equal(t, 5, op.Min(5, 5))
	assert.Equal(t, 5, op.Min(5, 6))
	assert.Equal(t, 5, op.Min(6, 5))
	assert.Equal(t, -10, op.Min(-5, -10))

	// Test with floats
	assert.Equal(t, 1.2, op.Min(1.2, 2.5))
	assert.Equal(t, -2.5, op.Min(1.2, -2.5))

	// Test with strings (lexicographic ordering)
	assert.Equal(t, "A", op.Min("a", "A"))
	assert.Equal(t, "apple", op.Min("apple", "banana"))
	assert.Equal(t, "", op.Min("", "test"))
}

func Test_Max(t *testing.T) {
	// Test with integers
	assert.Equal(t, 5, op.Max(5, 5))
	assert.Equal(t, 6, op.Max(5, 6))
	assert.Equal(t, 6, op.Max(6, 5))
	assert.Equal(t, -5, op.Max(-5, -10))

	// Test with floats
	assert.Equal(t, 2.5, op.Max(1.2, 2.5))
	assert.Equal(t, 1.2, op.Max(1.2, -2.5))

	// Test with strings (lexicographic ordering)
	assert.Equal(t, "a", op.Max("a", "A"))
	assert.Equal(t, "banana", op.Max("apple", "banana"))
	assert.Equal(t, "test", op.Max("", "test"))
}

func Test_IfElse(t *testing.T) {
	// Test with basic types
	assert.Equal(t, 5, op.IfElse(true, 5, 6))
	assert.Equal(t, 6, op.IfElse(false, 5, 6))

	// Test with strings
	assert.Equal(t, "true", op.IfElse(true, "true", "false"))
	assert.Equal(t, "false", op.IfElse(false, "true", "false"))

	// Test with different types
	assert.Equal(t, 0, op.IfElse(false, 42, 0))
	assert.Equal(t, 42, op.IfElse(true, 42, 0))
}

// Test IfElseErr function
func Test_IfElseErr(t *testing.T) {
	testErr := errors.New("test error")

	// Test with true condition
	result, err := op.IfElseErr(true, 42, testErr)
	assert.Equal(t, 42, result)
	assert.NoError(t, err)

	// Test with false condition
	result, err = op.IfElseErr(false, 42, testErr)
	assert.Equal(t, 0, result) // zero value for int
	assert.Error(t, err)
	assert.Equal(t, testErr, err)

	// Test with string type
	strResult, err := op.IfElseErr(true, "success", testErr)
	assert.Equal(t, "success", strResult)
	assert.NoError(t, err)

	strResult, err = op.IfElseErr(false, "success", testErr)
	assert.Equal(t, "", strResult) // zero value for string
	assert.Error(t, err)
}

func Test_IfElseDelay(t *testing.T) {
	assert.Equal(t, 5, op.IfElse(true, func() int { return 5 }, func() int { return 6 })())
	assert.Equal(t, 6, op.IfElse(false, func() int { return 5 }, func() int { return 6 })())
}

func Test_IfGetElse(t *testing.T) {
	// Test basic functionality
	assert.Equal(t, 5, op.IfGetElse(true, func() int { return 5 }, func() int { return 6 }))
	assert.Equal(t, 6, op.IfGetElse(false, func() int { return 5 }, func() int { return 6 }))

	// Test with side effects (verify correct function is called)
	called := 0
	trueFn := func() int { called = 1; return 5 }
	falseFn := func() int { called = 2; return 6 }

	// True case - should call trueFn
	called = 0
	result := op.IfGetElse(true, trueFn, falseFn)
	assert.Equal(t, 5, result)
	assert.Equal(t, 1, called)

	// False case - should call falseFn
	called = 0
	result = op.IfGetElse(false, trueFn, falseFn)
	assert.Equal(t, 6, result)
	assert.Equal(t, 2, called)
}

// Test IfGetElseGetErr function
func Test_IfGetElseGetErr(t *testing.T) {
	testErr := errors.New("test error")

	// Test with true condition
	result, err := op.IfGetElseGetErr(true, func() int { return 42 }, func() error { return testErr })
	assert.Equal(t, 42, result)
	assert.NoError(t, err)

	// Test with false condition
	result, err = op.IfGetElseGetErr(false, func() int { return 42 }, func() error { return testErr })
	assert.Equal(t, 0, result) // zero value for int
	assert.Error(t, err)
	assert.Equal(t, testErr, err)

	// Test with side effects
	called := 0
	trueFn := func() string { called = 1; return "success" }
	falseFn := func() error { called = 2; return testErr }

	// True case
	called = 0
	strResult, err := op.IfGetElseGetErr(true, trueFn, falseFn)
	assert.Equal(t, "success", strResult)
	assert.NoError(t, err)
	assert.Equal(t, 1, called)

	// False case
	called = 0
	strResult, err = op.IfGetElseGetErr(false, trueFn, falseFn)
	assert.Equal(t, "", strResult) // zero value for string
	assert.Error(t, err)
	assert.Equal(t, 2, called)
}

// Test Get function
func Test_Get(t *testing.T) {
	// Test with different return types
	assert.Equal(t, 42, op.Get(func() int { return 42 }))
	assert.Equal(t, "hello", op.Get(func() string { return "hello" }))
	assert.Equal(t, true, op.Get(func() bool { return true }))

	// Test with side effects
	called := false
	getterFn := func() int {
		called = true
		return 100
	}

	result := op.Get(getterFn)
	assert.Equal(t, 100, result)
	assert.True(t, called, "Getter function should be called")
}

// Test Compare function
func Test_Compare(t *testing.T) {
	// Test with integers
	assert.Equal(t, -1, op.Compare(5, 10))
	assert.Equal(t, 0, op.Compare(5, 5))
	assert.Equal(t, 1, op.Compare(10, 5))

	// Test with negative integers
	assert.Equal(t, -1, op.Compare(-10, -5))
	assert.Equal(t, 0, op.Compare(-5, -5))
	assert.Equal(t, 1, op.Compare(-5, -10))

	// Test with floats
	assert.Equal(t, -1, op.Compare(1.2, 2.5))
	assert.Equal(t, 0, op.Compare(2.5, 2.5))
	assert.Equal(t, 1, op.Compare(2.5, 1.2))

	// Test with strings (lexicographic ordering)
	assert.Equal(t, 1, op.Compare("apple", "Apple")) // lowercase comes after uppercase in ASCII
	assert.Equal(t, 0, op.Compare("apple", "apple"))
	assert.Equal(t, -1, op.Compare("apple", "banana"))
	assert.Equal(t, 1, op.Compare("banana", "apple"))

	// Test edge cases
	assert.Equal(t, 1, op.Compare("test", ""))
	assert.Equal(t, -1, op.Compare("", "test"))
	assert.Equal(t, 0, op.Compare("", ""))
}

// Additional edge case tests
func Test_EdgeCases(t *testing.T) {
	// Test Sum with zero values
	assert.Equal(t, 0, op.Sum(0, 0))
	assert.Equal(t, 5, op.Sum(0, 5))
	assert.Equal(t, 5, op.Sum(5, 0))

	// Test Sub with zero values
	assert.Equal(t, 0, op.Sub(0, 0))
	assert.Equal(t, -5, op.Sub(0, 5))
	assert.Equal(t, 5, op.Sub(5, 0))
	assert.Equal(t, 0, op.Sub(10, 10))

	// Test Min/Max with zero
	assert.Equal(t, 0, op.Min(0, 1))
	assert.Equal(t, -1, op.Min(0, -1))
	assert.Equal(t, 1, op.Max(0, 1))
	assert.Equal(t, 0, op.Max(0, -1))

	// Test Compare with zero
	assert.Equal(t, 0, op.Compare(0, 0))
	assert.Equal(t, -1, op.Compare(0, 1))
	assert.Equal(t, 1, op.Compare(1, 0))
	assert.Equal(t, 1, op.Compare(0, -1))
	assert.Equal(t, -1, op.Compare(-1, 0))

	// Test IfElse with complex types
	var ptr1, ptr2 *int
	val1, val2 := 1, 2
	ptr1, ptr2 = &val1, &val2
	result := op.IfElse(true, ptr1, ptr2)
	assert.Equal(t, ptr1, result)
	assert.Equal(t, 1, *result)

	// Test IfElseErr with nil error
	intResult, err := op.IfElseErr(false, 42, nil)
	assert.Equal(t, 0, intResult)
	assert.NoError(t, err)
}
