package test

import (
	"errors"
	"testing"

	"github.com/m4gshm/gollections/expr/use"
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

func Test_UseIfElse(t *testing.T) {
	result := use.If(true, 1).Else(2)
	assert.Equal(t, 1, result)

	result = use.If(false, 1).Else(2)
	assert.Equal(t, 2, result)
}

func Test_UseEval(t *testing.T) {
	result, ok := use.If(true, 1).If(true, 2).Eval()
	assert.True(t, ok)
	assert.Equal(t, 1, result)

	result, ok = use.If(false, 1).If(true, 2).Eval()
	assert.True(t, ok)
	assert.Equal(t, 2, result)

	result, ok = use.If(false, 1).If(false, 2).Eval()
	assert.False(t, ok)
	assert.Equal(t, 0, result)
}

func Test_UseElseZero(t *testing.T) {
	result := use.If(true, 1).If(true, 2).ElseZero()
	assert.Equal(t, 1, result)

	result = use.If(false, 1).If(true, 2).ElseZero()
	assert.Equal(t, 2, result)

	result = use.If(false, 1).If(false, 2).ElseZero()
	assert.Equal(t, 0, result)
}


func Test_UseIfIfElse(t *testing.T) {
	result := use.If(true, 1).If(true, 2).Else(3)
	assert.Equal(t, 1, result)

	result = use.If(false, 1).If(true, 2).Else(3)
	assert.Equal(t, 2, result)

	result = use.If(false, 1).If(false, 2).Else(3)
	assert.Equal(t, 3, result)
}

func Test_UseIfElseGet(t *testing.T) {
	result := use.If(false, 1).ElseGet(getTwo)
	assert.Equal(t, 2, result)

	result = use.If(true, 1).ElseGet(getTwo)
	assert.Equal(t, 1, result)
}

func Test_UseIfElseErr(t *testing.T) {
	result, e := use.If(false, 1).ElseErr(err)
	assert.Error(t, e)

	result, e = use.If(true, 1).ElseErr(err)
	assert.NoError(t, e)
	assert.Equal(t, 1, result)

	result = use.If(true, 1).If(true, 2).Else(3)
	assert.Equal(t, 1, result)
}

func Test_UseIfElseGetErr(t *testing.T) {
	result, _ := use.If(false, 1).ElseGetErr(getTwoErr)
	assert.Equal(t, 2, result)

	result, _ = use.If(true, 1).ElseGetErr(getTwoErr)
	assert.Equal(t, 1, result)
}

func Test_UseIfIfGetElseGet(t *testing.T) {
	result := use.If(false, 1).IfGet(false, getTwo).ElseGet(getThree)
	assert.Equal(t, 3, result)

	result = use.If(false, 1).IfGet(true, getTwo).ElseGet(getThree)
	assert.Equal(t, 2, result)

	result = use.If(true, 1).IfGet(true, getTwo).ElseGet(getThree)
	assert.Equal(t, 1, result)
}

func Test_UseIfIfGetErrElseGetErr(t *testing.T) {
	result, _ := use.If(false, 1).IfGetErr(false, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 3, result)

	result, _ = use.If(false, 1).IfGetErr(true, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 2, result)

	result, _ = use.If(true, 1).IfGetErr(true, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 1, result)
}

func Test_UseIfOtherElseGet(t *testing.T) {
	result := use.If(false, 1).Other(getFalse, getTwo).ElseGet(getThree)
	assert.Equal(t, 3, result)

	result = use.If(false, 1).Other(getTrue, getTwo).ElseGet(getThree)
	assert.Equal(t, 2, result)

	result = use.If(true, 1).Other(getTrue, getTwo).ElseGet(getThree)
	assert.Equal(t, 1, result)
}

func Test_UseIfOtherErrElseGetErr(t *testing.T) {
	result, _ := use.If(false, 1).OtherErr(getFalse, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 3, result)

	result, _ = use.If(false, 1).OtherErr(getTrue, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 2, result)

	result, _ = use.If(true, 1).OtherErr(getTrue, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 1, result)
}
