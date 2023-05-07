package test

import "testing"
import "github.com/m4gshm/gollections/op/delay/string_"
import "github.com/stretchr/testify/assert"

func TestConcat(t *testing.T) {

	result := string_.Of("1", "2", "3")
	assert.Equal(t, "123", result())
}
