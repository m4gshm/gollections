package test

import (
	"errors"
	"slices"
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seqe"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/stretchr/testify/assert"
)

func noErr[T any](t T) (T, error) { return t, nil }
func errOn(max int) func(int) (int, error) {
	return func(i int) (int, error) {
		if i == max {
			return i, errors.New("abort")
		}
		return i, nil
	}
}

func Test_AccumSum(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 5, 7, 9, 11), noErr)
	r, err := seqe.Accum(100, s, op.Sum[int])
	assert.Equal(t, 100+1+3+5+7+9+11, r)
	assert.NoError(t, err)
}

func Test_AccummSum(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 5, 7, 9, 11), noErr)
	r, err := seqe.Accumm(100, s, func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errors.New("stop")
		}
		return i1 + i2, nil
	})
	assert.Equal(t, 100+1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceSum(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr)
	sum, ok, err := seqe.ReduceOK(s, op.Sum)

	assert.True(t, ok)
	assert.Equal(t, 21, sum)
	assert.NoError(t, err)

	s = seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), errOn(4))
	sum, ok, err = seqe.ReduceOK(s, op.Sum)

	assert.True(t, ok)
	assert.Equal(t, 6, sum)
	assert.ErrorContains(t, err, "abort")

	s = seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), errOn(1))
	sum, ok, err = seqe.ReduceOK(s, op.Sum)

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
	assert.ErrorContains(t, err, "abort")
}

func Test_ReduceeSum(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 5, 7, 9, 11), noErr)
	r, ok, err := seqe.ReduceeOK(s, func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errors.New("stop")
		}
		return i1 + i2, nil
	})
	assert.True(t, ok)
	assert.Equal(t, 1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")

	s = seq.ToSeq2(seq.Of(1, 3, 5, 7, 9, 11), errOn(5))
	r, ok, err = seqe.ReduceeOK(s, func(i1, i2 int) (int, error) { return i1 + i2, nil })

	assert.True(t, ok)
	assert.Equal(t, 1+3, r)
	assert.ErrorContains(t, err, "abort")

	s = seq.ToSeq2(seq.Of(1, 3, 5, 7, 9, 11), errOn(1))
	r, ok, err = seqe.ReduceeOK(s, func(i1, i2 int) (int, error) { return i1 + i2, nil })

	assert.False(t, ok)
	assert.Equal(t, 0, r)
	assert.ErrorContains(t, err, "abort")
}

func Test_ReduceeSumFirstErr(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 5, 7, 9, 11), noErr)
	r, ok, err := seqe.ReduceeOK(s, func(_, _ int) (int, error) {
		return 0, errors.New("stop")
	})
	assert.True(t, ok)
	assert.Equal(t, 0, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceEmpty(t *testing.T) {
	s := seq.ToSeq2(seq.Of[int](), func(i int) (int, error) { return i, errors.New("unexpected") })
	sum, ok, err := seqe.ReduceOK(s, op.Sum)

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
	assert.NoError(t, err)
}

func Test_ReduceNil(t *testing.T) {
	var s seqe.SeqE[int]
	sum, ok, err := seqe.ReduceOK(s, op.Sum)

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
	assert.NoError(t, err)
}

func Test_First(t *testing.T) {
	result, ok, err := seqe.First(seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr), more.Than(5)) //7, true

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

}

func Test_Firstt(t *testing.T) {
	result, ok, err := seq.Firstt(seq.Of(1, 2, 3, 4, 5, 6), func(i int) (bool, error) {
		return more.Than(5)(i), nil
	})

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	result, ok, err = seq.Firstt(seq.Of(1, 2, 3, 4, 5, 6), func(_ int) (bool, error) {
		return true, errors.New("abort")
	})

	assert.False(t, ok)
	assert.Equal(t, 0, result)
	assert.ErrorContains(t, err, "abort")
}

var even = func(v int) bool { return v%2 == 0 }

func Test_Flat(t *testing.T) {
	md := seq.ToSeq2(seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...), noErr)
	f := seqe.Flat(md, as.Is)
	s, err := seqe.Slice(f)
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, s)
	assert.NoError(t, err)
}

func Test_FlatSeq(t *testing.T) {
	md := seq.ToSeq2(seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...), noErr)
	f := seqe.FlatSeq(md, slices.Values)
	s, err := seqe.Slice(f)
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, s)
	assert.NoError(t, err)
}

func Test_Filter(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := seq.Filter(s, even)
	assert.Equal(t, slice.Of(4, 8), seq.Slice(r))
}

func Test_Filt(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	l := seq.Filt(s, func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errors.New("abort"), nil) })
	r, err := seqe.Slice(l)
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4, 8), r)
}

func Test_Filt2(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	l := seq.Filt(s, func(i int) (bool, error) {
		ok := i <= 7
		return ok && even(i), op.IfElse(ok, nil, errors.New("abort"))
	})
	r, err := seqe.Slice(l)
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4), r)
}

func Test_Contains(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 2, 3), noErr)
	ok, err := seqe.Contains(s, 3)
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = seqe.Contains(s, 0)
	assert.False(t, ok)
	assert.NoError(t, err)

	s = seq.ToSeq2(seq.Of(1, 2, 3), errOn(1))
	ok, err = seqe.Contains(s, 3)
	assert.False(t, ok)
	assert.ErrorContains(t, err, "abort")
}

// func Test_KeyValue(t *testing.T) {
// 	s := seq.Of(1, 2, 3)
// 	s2 := seqe.KeyValue(s, as.Is, strconv.Itoa)
// 	k := seq.Slice(seq2.Keys(s2))
// 	v := seq.Slice(seq2.Values(s2))

// 	assert.Equal(t, slice.Of(1, 2, 3), k)
// 	assert.Equal(t, slice.Of("1", "2", "3"), v)
// }

// func Test_KeyValues(t *testing.T) {
// 	s := seq.Of([]int{1, 2}, []int{3})
// 	s2 := seq.KeyValues(s, slice.Len, as.Is)
// 	k := seq.Slice(seq2.Keys(s2))
// 	v := seq.Slice(seq2.Values(s2))

//		assert.Equal(t, slice.Of(2, 2, 1), k)
//		assert.Equal(t, slice.Of(1, 2, 3), v)
//	}
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
	var in seq.SeqE[int]
	var out []int

	iter := false
	for e, err := range seqe.Convert(in, as.Is) {
		assert.NoError(t, err)
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
	from := seq.ToSeq2(seq.Of(1, 2, 3, 5, 7, 8, 9, 11), noErr)
	s := []string{}

	for e, err := range seqe.Convert(from, strconv.Itoa) {
		assert.NoError(t, err)
		s = append(s, e)
	}

	assert.Equal(t, slice.Of("1", "2", "3", "5", "7", "8", "9", "11"), s)
}

func Test_AllConv(t *testing.T) {
	from := seq.ToSeq2(seq.Of("1", "2", "3", "5", "_7", "8", "9", "11"), noErr)
	i := []int{}

	for v, err := range seqe.Conv(from, strconv.Atoi) {
		if err == nil {
			i = append(i, v)
		}
	}

	assert.Equal(t, slice.Of(1, 2, 3, 5, 8, 9, 11), i)
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := seq.ConvertOK(s, func(i int) (string, bool) { return strconv.Itoa(i), even(i) })
	assert.Equal(t, []string{"4", "8"}, seq.Slice(r))
}

func Test_ConvFilteredInplace(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := seq.ConvOK(s, func(i int) (string, bool, error) { return strconv.Itoa(i), even(i), nil })
	o, _ := seqe.Slice(r)
	assert.Equal(t, []string{"4", "8"}, o)
}
