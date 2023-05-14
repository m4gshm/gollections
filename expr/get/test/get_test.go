package test

import (
	"errors"
	"testing"

	"github.com/m4gshm/gollections/expr/get"
	"github.com/stretchr/testify/assert"
)

var (
	err         = errors.New("test")
	getErr      = func() (int, error) { return 0, err }
	getOne      = func() int { return 1 }
	getTwo      = func() int { return 2 }
	getThree    = func() int { return 3 }
	getOneErr   = func() (int, error) { return 1, nil }
	getTwoErr   = func() (int, error) { return 2, nil }
	getThreeErr = func() (int, error) { return 3, nil }
	getTrue     = func() bool { return true }
	getFalse    = func() bool { return false }
)

func Test_GetIfElse(t *testing.T) {
	result := get.If(true, getOne).Else(2)
	assert.Equal(t, 1, result)

	result = get.If(false, getOne).Else(2)
	assert.Equal(t, 2, result)
}

func Test_GetIfElseGet(t *testing.T) {
	result := get.If(false, getOne).ElseGet(getTwo)
	assert.Equal(t, 2, result)

	result = get.If(true, getOne).ElseGet(getTwo)
	assert.Equal(t, 1, result)
}

func Test_GetIfElseErr(t *testing.T) {
	result, e := get.If(false, getOne).ElseErr(err)
	assert.Error(t, e)

	result, e = get.If(true, getOne).ElseErr(err)
	assert.NoError(t, e)
	assert.Equal(t, 1, result)
}

func Test_GetIfIfElse(t *testing.T) {
	result := get.If(false, getOne).If(false, 2).Else(3)
	assert.Equal(t, 3, result)

	result = get.If(false, getOne).If(true, 2).Else(3)
	assert.Equal(t, 2, result)

	result = get.If(true, getOne).If(true, 2).Else(3)
	assert.Equal(t, 1, result)
}

func Test_GetIfElseGetErr(t *testing.T) {
	result, _ := get.If(false, getOne).ElseGetErr(getTwoErr)
	assert.Equal(t, 2, result)

	result, _ = get.If(true, getOne).ElseGetErr(getTwoErr)
	assert.Equal(t, 1, result)
}

func Test_GetIfIfGetElseGet(t *testing.T) {
	result := get.If(false, getOne).IfGet(false, getTwo).ElseGet(getThree)
	assert.Equal(t, 3, result)

	result = get.If(false, getOne).IfGet(true, getTwo).ElseGet(getThree)
	assert.Equal(t, 2, result)

	result = get.If(true, getOne).IfGet(true, getTwo).ElseGet(getThree)
	assert.Equal(t, 1, result)
}

func Test_GetIfIfGetErrElseGetErr(t *testing.T) {
	result, _ := get.If(false, getOne).IfGetErr(false, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 3, result)

	result, _ = get.If(false, getOne).IfGetErr(true, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 2, result)

	result, _ = get.If(true, getOne).IfGetErr(true, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 1, result)
}

func Test_GetIfOtherElseGet(t *testing.T) {
	result := get.If(false, getOne).Other(getFalse, getTwo).ElseGet(getThree)
	assert.Equal(t, 3, result)

	result = get.If(false, getOne).Other(getTrue, getTwo).ElseGet(getThree)
	assert.Equal(t, 2, result)

	result = get.If(true, getOne).Other(getTrue, getTwo).ElseGet(getThree)
	assert.Equal(t, 1, result)
}

func Test_GetIfOtherErrElseGetErr(t *testing.T) {
	result, _ := get.If(false, getOne).OtherErr(getFalse, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 3, result)

	result, _ = get.If(false, getOne).OtherErr(getTrue, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 2, result)

	result, _ = get.If(true, getOne).OtherErr(getTrue, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 1, result)
}
