//go:build goexperiment.rangefunc

package test

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/seq2"
	"github.com/stretchr/testify/assert"
)

func Test_OverFiltered(t *testing.T) {
	integers := slice.Of(1, 2, 3, 5, 7, 8, 9, 11)

	s := []int{}

	for _, e := range seq2.Filtered(integers, func(e int) bool { return e%2 == 0 }) {
		s = append(s, e)
	}

	assert.Equal(t, slice.Of(2, 8), s)
}

func Test_OverConverted(t *testing.T) {
	integers := slice.Of(1, 2, 3, 5, 7, 8, 9, 11)
	s := []string{}

	for _, e := range seq2.Converted(integers, strconv.Itoa) {
		s = append(s, e)
	}

	assert.Equal(t, slice.Of("1", "2", "3", "5", "7", "8", "9", "11"), s)
}
