package test

import (
	"fmt"
	"testing"

	"github.com/m4gshm/gollections/error_"
	"github.com/stretchr/testify/assert"
)

type testError string

var _ error = testError("test error")

func (t testError) Error() string {
	return string(t)
}

func Test_error_As(t *testing.T) {
	e := fmt.Errorf("wrap: %w", testError("test error"))

	err, ok := error_.As[testError](e)

	assert.True(t, ok)
	assert.Equal(t, testError("test error"), err)
}
