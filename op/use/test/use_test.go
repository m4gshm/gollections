package test

import (
	"testing"

	"github.com/m4gshm/gollections/op/use"
	"github.com/stretchr/testify/assert"
)

var (
	one = func() int { return 1 }
	two = func() int { return 2 }
)

func Test_UseIfElse(t *testing.T) {
	result := use.If(false, 1).Else(2)
	assert.Equal(t, 2, result)

	result = use.If(true, 1).Else(2)
	assert.Equal(t, 1, result)

	result = use.When[int]{Condition: false, Then: 1}.Else(2)
	assert.Equal(t, 2, result)
}

func Test_UseIfElseGet(t *testing.T) {
	result := use.If(false, 1).ElseGet(func() int { return 2 })
	assert.Equal(t, 2, result)

	result = use.If(true, 1).ElseGet(func() int { return 2 })
	assert.Equal(t, 1, result)
}

func Test_UseThisIfElse(t *testing.T) {

	result := use.This(1).If(false).Else(2)
	assert.Equal(t, 2, result)

	result = use.This(1).If(true).Else(2)
	assert.Equal(t, 1, result)

	result = use.ThisOne[int]{Value: 1}.If(true).Else(2)
	assert.Equal(t, 1, result)
}

func Test_UseOneIfElse(t *testing.T) {
	result := use.One(1).If(false).Else(2)
	assert.Equal(t, 2, result)

	result = use.One(1).If(true).Else(2)
	assert.Equal(t, 1, result)
}
