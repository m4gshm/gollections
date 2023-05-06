package test

import (
	"testing"

	"github.com/m4gshm/gollections/op/get"
	"github.com/stretchr/testify/assert"
)

var (
	getOne  = func() int { return 1 }
	getTwo  = func() int { return 2 }
	isFalse = func() bool { return false }
	isTrue  = func() bool { return true }
)

func Test_GetIf(t *testing.T) {

	result := get.If(getOne, false).Else(2)
	assert.Equal(t, 2, result)

	result = get.If(getOne, true).Else(2)
	assert.Equal(t, 1, result)

	result = get.If_[int]{Condition: false, Then: getOne}.Else(2)
	assert.Equal(t, 2, result)
}

func Test_GetIfElseGet(t *testing.T) {
	result := get.If(getOne, false).ElseGet(getTwo)
	assert.Equal(t, 2, result)

	result = get.If(getOne, true).ElseGet(getTwo)
	assert.Equal(t, 1, result)
}

func Test_GetOneIf(t *testing.T) {
	result := get.One(getOne).If(false).Else(2)
	assert.Equal(t, 2, result)

	result = get.One(getOne).If(true).Else(2)
	assert.Equal(t, 1, result)

	result = get.One_[int]{Getter: getOne}.If(true).Else(2)
	assert.Equal(t, 1, result)
}

func Test_GetOneIfCalc(t *testing.T) {
	result := get.One(getOne).If_(isFalse).ElseGet(getTwo)
	assert.Equal(t, 2, result)

	result = get.One(getOne).If_(isTrue).Else(2)
	assert.Equal(t, 1, result)

	result = get.One_[int]{Getter: getOne}.If_(isTrue).ElseGet(getTwo)
	assert.Equal(t, 1, result)
}
