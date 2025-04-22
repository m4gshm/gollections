package test

import (
	"errors"
	"iter"
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
func errOn[T comparable](errVal T) func(T) (T, error) {
	return func(val T) (T, error) {
		if val == errVal {
			return val, errors.New("abort")
		}
		return val, nil
	}
}

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")
	result := seqe.OfIndexed(len(indexed), func(i int) (string, error) { return indexed[i], nil })
	out, err := seqe.Slice(result)

	assert.Equal(t, indexed, out)
	assert.NoError(t, err)

	result = seqe.OfIndexed(len(indexed), func(i int) (string, error) { return indexed[i], op.IfElse(i == 3, errors.New("abort"), nil) })
	out, err = seqe.Slice(result)

	assert.Equal(t, slice.Of("0", "1", "2"), out)
	assert.ErrorContains(t, err, "abort")
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
	result, ok, err := seqe.Firstt(seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr), func(i int) (bool, error) {
		return more.Than(5)(i), nil
	})

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	result, ok, err = seqe.Firstt(seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr), func(_ int) (bool, error) {
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

func Test_Flatt(t *testing.T) {
	var (
		input     iter.Seq2[[]string, error]
		flattener func([]string) ([]int, error)
		out       seq.SeqE[int]
	)
	out = seqe.Flatt(input, flattener)
	for i, err := range out {
		if err != nil {
			panic(err)
		}
		_ = i
	}

	s := seq.ToSeq2(seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5", "6"}}...), noErr)
	f := func(strInteger []string) ([]int, error) { return slice.Conv(strInteger, strconv.Atoi) }
	i, err := seqe.Slice(seqe.Flatt(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, i)
	assert.ErrorContains(t, err, "parsing \"_5\"")
}

func Test_FlattSeq(t *testing.T) {
	var (
		input     iter.Seq2[[]string, error]
		flattener func([]string) seq.SeqE[int]
		out       seq.SeqE[int]
	)
	out = seqe.FlattSeq(input, flattener)
	for i, err := range out {
		if err != nil {
			panic(err)
		}
		_ = i
	}

	s := seq.ToSeq2(seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5", "6"}}...), noErr)
	f := func(strInteger []string) seq.SeqE[int] { return seq.Conv(seq.Of(strInteger...), strconv.Atoi) }
	i, err := seqe.Slice(seqe.FlattSeq(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, i)
	assert.ErrorContains(t, err, "parsing \"_5\"")
}

func Test_Filter(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), noErr)
	f := seqe.Filter(s, even)
	r, err := seqe.Slice(f)
	assert.Equal(t, slice.Of(4, 8), r)
	assert.NoError(t, err)
}

func Test_Filt(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), noErr)
	filter := func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errors.New("abort"), nil) }
	l := seqe.Filt(s, filter)
	r, err := seqe.Slice(l)
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4, 8), r)

	l = seqe.Filt(s, nil)
	r, err = seqe.Slice(l)
	assert.NoError(t, err)
	assert.Empty(t, r)

	l = seqe.Filt[seq.SeqE[int]](nil, filter)
	r, err = seqe.Slice(l)
	assert.NoError(t, err)
	assert.Empty(t, r)
}

func Test_Filt2(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), noErr)
	l := seqe.Filt(s, func(i int) (bool, error) {
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
	s := seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), noErr)
	r := seqe.ConvertOK(s, func(i int) (string, bool) { return strconv.Itoa(i), even(i) })
	out, err := seqe.Slice(r)
	assert.Equal(t, []string{"4", "8"}, out)
	assert.NoError(t, err)
}

func Test_ConvFilteredInplace(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), noErr)
	r := seqe.ConvOK(s, func(i int) (string, bool, error) { return strconv.Itoa(i), even(i), nil })
	o, err := seqe.Slice(r)
	assert.Equal(t, []string{"4", "8"}, o)
	assert.NoError(t, err)

	r = seqe.ConvOK(s, func(i int) (string, bool, error) {
		return strconv.Itoa(i), even(i), op.IfElse(i == 9, errors.New("abort"), nil)
	})
	o, err = seqe.Slice(r)

	assert.Equal(t, []string{"4", "8"}, o)
	assert.ErrorContains(t, err, "abort")
}
