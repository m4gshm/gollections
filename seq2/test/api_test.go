package test

import (
	"testing"

	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/stretchr/testify/assert"
)

func Test_Map(t *testing.T) {
	s := seq2.Of("first", "second", "third")
	m := seq2.Map(s)

	assert.Equal(t, "first", m[0])
	assert.Equal(t, "second", m[1])
	assert.Equal(t, "third", m[2])
}

func Test_Keys_Values(t *testing.T) {
	s := seq2.Of("first", "second", "third")
	k := seq.Slice(seq2.Keys(s))
	v := seq.Slice(seq2.Values(s))
	assert.Equal(t, slice.Of(0, 1, 2), k)
	assert.Equal(t, slice.Of("first", "second", "third"), v)
}

func Test_Group(t *testing.T) {
	s := seq2.Convert(seq2.Of("first", "second", "third"), func(i int, s string) (bool, string) { return i%2 == 0, s })
	m := seq2.Group(s)

	assert.Equal(t, slice.Of("first", "third"), sort.Asc(m[true]))
	assert.Equal(t, slice.Of("second"), sort.Asc(m[false]))
}

func Test_Filter(t *testing.T) {
	s := seq2.Filter(seq2.Of("first", "second", "third"), func(i int, _ string) bool { return i%2 == 0 })
	k := seq.Slice(seq2.Keys(s))
	v := seq.Slice(seq2.Values(s))

	assert.Equal(t, slice.Of(0, 2), k)
	assert.Equal(t, slice.Of("first", "third"), v)
}
