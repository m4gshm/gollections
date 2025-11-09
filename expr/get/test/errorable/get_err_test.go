package errorable

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m4gshm/gollections/expr/get"
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
	_, err := get.If_(true, getErr).Else(2)
	require.Error(t, err)

	result, err := get.If_(false, getErr).Else(2)
	require.NoError(t, err)
	assert.Equal(t, 2, result)
}

func Test_UseIfOtherElseGet(t *testing.T) {
	result, _ := get.If_(false, getOneErr).Other(getFalse, getTwo).ElseGet(getThree)
	assert.Equal(t, 3, result)

	result, _ = get.If_(false, getOneErr).Other(getTrue, getTwo).ElseGet(getThree)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).Other(getTrue, getTwo).ElseGet(getThree)
	assert.Equal(t, 1, result)
}

func Test_UseIfOtherErrElseGetErr(t *testing.T) {
	result, _ := get.If_(false, getOneErr).OtherErr(getFalse, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 3, result)

	result, _ = get.If_(false, getOneErr).OtherErr(getTrue, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 2, result)

	result, _ = get.If_(true, getOneErr).OtherErr(getTrue, getTwoErr).ElseGetErr(getThreeErr)
	assert.Equal(t, 1, result)
}
