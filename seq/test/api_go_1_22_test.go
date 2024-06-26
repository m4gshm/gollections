//go:build goexperiment.rangefunc

package test

import (
	"iter"
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/stretchr/testify/assert"
)

func Test_SeqOfNil(t *testing.T) {
	var in, out []int

	iter := false
	for e := range seq.Of(in...) {
		iter = true
		out = append(out, e)
	}

	assert.Nil(t, out)
	assert.False(t, iter)
}

func Test_ConvertNilSeq(t *testing.T) {
	var in iter.Seq[int]
	var out []int

	iter := false
	for e := range seq.Convert(in, as.Is) {
		iter = true
		out = append(out, e)
	}

	assert.Nil(t, out)
	assert.False(t, iter)
}

func Test_AllFiltered(t *testing.T) {
	from := seq.Of(1, 2, 3, 5, 7, 8, 9, 11)

	s := []int{}

	for e := range seq.Filter(from, func(e int) bool { return e%2 == 0 }) {
		s = append(s, e)
	}

	assert.Equal(t, slice.Of(2, 8), sort.Asc(s))
}

func Test_AllConverted(t *testing.T) {
	from := set.Of(1, 2, 3, 5, 7, 8, 9, 11)
	s := []string{}

	for e := range seq.Convert(from.All, strconv.Itoa) {
		s = append(s, e)
	}

	assert.Equal(t, slice.Of("1", "2", "3", "5", "7", "8", "9", "11"), s)
}

func Test_AllConv(t *testing.T) {
	from := set.Of("1", "2", "3", "5", "_7", "8", "9", "11")
	i := []int{}

	for v, err := range seq.Conv(from.All, strconv.Atoi) {
		if err == nil {
			i = append(i, v)
		}
	}

	assert.Equal(t, slice.Of(1, 2, 3, 5, 8, 9, 11), i)
}
