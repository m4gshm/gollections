package errorable

import (
	"errors"
	"testing"

	"github.com/m4gshm/gollections/expr/get"
	"github.com/stretchr/testify/assert"
)

var (
	err    = errors.New("test")
	getErr = func() (int, error) { return 0, err }
)

func Test_GetIfElse(t *testing.T) {
	result, err := get.If_(true, getErr).Else(2)
	assert.Error(t, err)

	result, err = get.If_(false, getErr).Else(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, result)
}
