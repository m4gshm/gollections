package test

import (
	"errors"
	"iter"
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_AccumSum(t *testing.T) {
	s := seq.Of(1, 3, 5, 7, 9, 11)
	r := seq.Accum(100, s, op.Sum[int])
	assert.Equal(t, 100+1+3+5+7+9+11, r)
}

func Test_AccummSum(t *testing.T) {
	s := seq.Of(1, 3, 5, 7, 9, 11)
	r, err := seq.Accumm(100, s, func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errors.New("stop")
		}
		return i1 + i2, nil
	})
	assert.Equal(t, 100+1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceSum(t *testing.T) {
	sum, ok := seq.ReduceOK(seq.Of(1, 2, 3, 4, 5, 6), op.Sum)

	assert.True(t, ok)
	assert.Equal(t, 21, sum)
}

func Test_ReduceeSum(t *testing.T) {
	s := seq.Of(1, 3, 5, 7, 9, 11)
	r, ok, err := seq.ReduceeOK(s, func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errors.New("stop")
		}
		return i1 + i2, nil
	})
	assert.True(t, ok)
	assert.Equal(t, 1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceeSumFirstErr(t *testing.T) {
	s := seq.Of(1, 3, 5, 7, 9, 11)
	r, ok, err := seq.ReduceeOK(s, func(_, _ int) (int, error) {
		return 0, errors.New("stop")
	})
	assert.True(t, ok)
	assert.Equal(t, 0, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceEmpty(t *testing.T) {
	s := seq.Of[int]()
	sum, ok := seq.ReduceOK(s, op.Sum)

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
}

func Test_ReduceNil(t *testing.T) {
	var s iter.Seq[int]
	sum, ok := seq.ReduceOK(s, op.Sum)

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
}

func Test_First(t *testing.T) {

	result, ok := seq.First(seq.Of(1, 2, 3, 4, 5, 6), more.Than(5)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 6, result)

}

var even = func(v int) bool { return v%2 == 0 }

func Test_Flat(t *testing.T) {
	md := seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := seq.Flat(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, seq.Slice(f))
}

func Test_Filter(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := seq.Filter(s, even)
	assert.Equal(t, slice.Of(4, 8), seq.Slice(r))
}

func Test_Filt(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	l := seq.Filt(s, func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errors.New("abort"), nil) })
	r, err := seq2.Slice(l)
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4, 8), r)
}

func Test_Filt2(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	l := seq.Filt(s, func(i int) (bool, error) {
		ok := i <= 7
		return ok && even(i), op.IfElse(ok, nil, errors.New("abort"))
	})
	r, err := seq2.Slice(l)
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4), r)
}

func Test_Contains(t *testing.T) {
	assert.True(t, seq.Contains(seq.Of(1, 2, 3), 3))
	assert.False(t, seq.Contains(seq.Of(1, 2, 3), 0))
}

func Test_KeyValue(t *testing.T) {
	s := seq.Of(1, 2, 3)
	s2 := seq.KeyValue(s, as.Is, strconv.Itoa)
	k := seq.Slice(seq2.Keys(s2))
	v := seq.Slice(seq2.Values(s2))

	assert.Equal(t, slice.Of(1, 2, 3), k)
	assert.Equal(t, slice.Of("1", "2", "3"), v)
}
