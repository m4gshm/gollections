package test

import (
	"errors"
	"testing"

	"github.com/m4gshm/gollections/expr/get"
	"github.com/stretchr/testify/assert"
)

var (
	getOne = func() int { return 1 }
)

func Test_GetIfElse(t *testing.T) {
	result := get.If(true, getOne).Else(2)
	assert.Equal(t, 1, result)

	result = get.If(false, getOne).Else(2)
	assert.Equal(t, 2, result)
}
