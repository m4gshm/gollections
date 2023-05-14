package errorable

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
	result, err := get.If_(true, getErr).Else(2)
	assert.Error(t, err)

	result, err = get.If_(false, getErr).Else(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, result)
}

func Test_GetIfElseGet(t *testing.T) {
	result, _ := get.If_(false, getOneErr).ElseGet(getTwo)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).ElseGet(getTwo)
	assert.Equal(t, 1, result)
}

func Test_GetIfElseErr(t *testing.T) {
	result, e := get.If_(false, getOneErr).ElseErr(err)
	assert.Error(t, e)

	result, e = get.If_(true, getOneErr).ElseErr(err)
	assert.NoError(t, e)
	assert.Equal(t, 1, result)
}

func Test_GetIfIfElse(t *testing.T) {
	result, _ := get.If_(false, getOneErr).If(false, 2).Else(3)
	assert.Equal(t, 3, result)

	result, _ = get.If_(false, getOneErr).If(true, 2).Else(3)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).If(true, 2).Else(3)
	assert.Equal(t, 1, result)
}

func Test_GetIfElseGetErr(t *testing.T) {
	result, _ := get.If_(false, getOneErr).ElseGetErr(getTwoErr)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).ElseGetErr(getTwoErr)
	assert.Equal(t, 1, result)
}

func Test_GetIfIfGetElseGet(t *testing.T) {
	result, _ := get.If_(false, getOneErr).IfGet(false, getTwo).ElseGet(getThree)
	assert.Equal(t, 3, result)

	result, _ = get.If_(false, getOneErr).IfGet(true, getTwo).ElseGet(getThree)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).IfGet(true, getTwo).ElseGet(getThree)
	assert.Equal(t, 1, result)
}

func Test_GetIfIfGetErrElseGetErr(t *testing.T) {
	result, _ := get.If_(false, getOneErr).IfGetErr(false, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 3, result)

	result, _ = get.If_(false, getOneErr).IfGetErr(true, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).IfGetErr(true, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 1, result)
}

func Test_GetIfOtherElseGet(t *testing.T) {
	result, _ := get.If_(false, getOneErr).Other(getFalse, getTwo).ElseGet(getThree)
	assert.Equal(t, 3, result)

	result, _ = get.If_(false, getOneErr).Other(getTrue, getTwo).ElseGet(getThree)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).Other(getTrue, getTwo).ElseGet(getThree)
	assert.Equal(t, 1, result)
}

func Test_GetIfOtherErrElseGetErr(t *testing.T) {
	result, _ := get.If_(false, getOneErr).OtherErr(getFalse, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 3, result)

	result, _ = get.If_(false, getOneErr).OtherErr(getTrue, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).OtherErr(getTrue, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 1, result)
}
