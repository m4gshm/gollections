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
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/seqe"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/stretchr/testify/assert"
)

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")
	result := seq.OfIndexed(len(indexed), func(i int) string { return indexed[i] })
	assert.Equal(t, indexed, seq.Slice(result))
}

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

	result, ok := seq.First(seq.Of(1, 2, 3, 4, 5, 6), more.Than(5))

	assert.True(t, ok)
	assert.Equal(t, 6, result)

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
	md := seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := seq.Flat(md, as.Is)

	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, seq.Slice(f))
}

func Test_FlatSeq(t *testing.T) {
	md := seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := seq.FlatSeq(md, slices.Values)

	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, seq.Slice(f))
}

func Test_Flatt(t *testing.T) {
	var (
		input     iter.Seq[[]string]
		flattener func([]string) seq.SeqE[int]
		out       seq.SeqE[int]
	)
	out = seq.Flatt(input, flattener)
	for i, err := range out {
		if err != nil {
			panic(err)
		}
		_ = i
	}

	s := seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5", "6"}}...)
	f := func(strInteger []string) seq.SeqE[int] { return seq.Conv(seq.Of(strInteger...), strconv.Atoi) }
	i, err := seqe.Slice(seq.Flatt(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, i)
	assert.ErrorContains(t, err, "parsing \"_5\"")
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

func Test_KeyValues(t *testing.T) {
	s := seq.Of([]int{1, 2}, []int{3})
	s2 := seq.KeyValues(s, slice.Len, as.Is)
	k := seq.Slice(seq2.Keys(s2))
	v := seq.Slice(seq2.Values(s2))

	assert.Equal(t, slice.Of(2, 2, 1), k)
	assert.Equal(t, slice.Of(1, 2, 3), v)
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
	from := seq.Of(1, 2, 3, 5, 7, 8, 9, 11)
	s := []string{}

	for e := range seq.Convert(from, strconv.Itoa) {
		s = append(s, e)
	}

	assert.Equal(t, slice.Of("1", "2", "3", "5", "7", "8", "9", "11"), s)
}

func Test_AllConv(t *testing.T) {
	from := seq.Of("1", "2", "3", "5", "_7", "8", "9", "11")
	i := []int{}

	for v, err := range seq.Conv(from, strconv.Atoi) {
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
