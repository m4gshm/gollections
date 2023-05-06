package test

import (
	"testing"

	"github.com/m4gshm/gollections/op/get"
	"github.com/stretchr/testify/assert"
)

var (
	getOne = func() int { return 1 }
	getTwo = func() int { return 2 }
)

func Test_GetIfElse(t *testing.T) {

	result := get.If(false, getOne).Else(2)
	assert.Equal(t, 2, result)

	result = get.If(true, getOne).Else(2)
	assert.Equal(t, 1, result)

	result = get.When[int]{Condition: false, Then: getOne}.Else(2)
	assert.Equal(t, 2, result)
}

func Test_GetIfElseGet(t *testing.T) {
	result := get.If(false, getOne).ElseGet(getTwo)
	assert.Equal(t, 2, result)

	result = get.If(true, getOne).ElseGet(getTwo)
	assert.Equal(t, 1, result)
}

func Test_GetThisIfElse(t *testing.T) {
	result := get.This(getOne).If(false).Else(2)
	assert.Equal(t, 2, result)

	result = get.This(getOne).If(true).Else(2)
	assert.Equal(t, 1, result)

	result = get.ThisOne[int]{Value: getOne}.If(true).Else(2)
	assert.Equal(t, 1, result)
}

func Test_GetOneIfElse(t *testing.T) {
	result := get.One(getOne).If(false).Else(2)
	assert.Equal(t, 2, result)

	result = get.One(getOne).If(true).Else(2)
	assert.Equal(t, 1, result)
}
