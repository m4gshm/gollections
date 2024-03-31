//go:build goexperiment.rangefunc

package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/break/loop"
)

func Test_IterAll(t *testing.T) {
	var (
		r    []int
		rerr error
	)
	for v, err := range loop.Conv(loop.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi).All {
		if rerr = err; err == nil {
			r = append(r, v)
		}
	}

	assert.Equal(t, []int{1, 3, 5}, r)
	assert.Error(t, rerr, "invalid syntax")
}
