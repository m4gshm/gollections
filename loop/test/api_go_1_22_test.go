//go:build goexperiment.rangefunc

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
)

func Test_IterAll(t *testing.T) {

	r := []int{}
	for v := range loop.Of(1, 3, 5, 7, 9, 11).All {
		r = append(r, v)
	}

	assert.Equal(t, []int{1, 3, 5, 7, 9, 11}, r)
}
