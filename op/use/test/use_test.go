package test

import (
	"testing"

	"github.com/m4gshm/gollections/op/use"
	"github.com/stretchr/testify/assert"
)

func Test_UseIf(t *testing.T) {
	result := use.If(1, false).Else(2)
	assert.Equal(t, 2, result)

	result = use.If(1, true).Else(2)
	assert.Equal(t, 1, result)

	result = use.If_[int]{Condition: false, Then: 1}.Else(2)
	assert.Equal(t, 2, result)
}

func Test_UseIfElseGet(t *testing.T) {
	result := use.If(1, false).ElseGet(func() int { return 2 })
	assert.Equal(t, 2, result)

	result = use.If(1, true).ElseGet(func() int { return 2 })
	assert.Equal(t, 1, result)
}

func Test_UseOneIf(t *testing.T) {
	result := use.One(1).If(false).Else(2)
	assert.Equal(t, 2, result)

	result = use.One(1).If(true).Else(2)
	assert.Equal(t, 1, result)

	result = use.One_[int]{Value: 1}.If(true).Else(2)
	assert.Equal(t, 1, result)
}
