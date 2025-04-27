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

func Test_Of(t *testing.T) {
	sequence := seq.Of(0, 1, 2, 3, 4)
	var out []int
	for v := range sequence {
		out = append(out, v)
	}
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), out)
	out = nil
	for v := range sequence {
		if v == 1 {
			break
		}
		out = append(out, v)
	}
	assert.Equal(t, slice.Of(0), out)

	out = nil
	for v := range sequence {
		_ = v
		break
	}
	assert.Nil(t, out)
}

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")

	getAt := func(i int) string { return indexed[i] }
	sequence := seq.OfIndexed(len(indexed), getAt)
	assert.Equal(t, indexed, seq.Slice(sequence))

	var out []string
	var iter = false
	for v := range sequence {
		iter = true
		if v == "3" {
			break
		}
		out = append(out, v)
	}
	assert.True(t, iter)
	assert.Equal(t, slice.Of("0", "1", "2"), out)

	sequence = seq.OfIndexed(len(indexed), (func(i int) string)(nil))

	assert.Nil(t, seq.Slice(sequence))
}

func Test_Series(t *testing.T) {
	generator := func(prev int) (int, bool) { return prev + 1, prev < 3 }
	sequence := seq.Series(-1, generator)
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), seq.Slice(sequence))

	var out []int
	for v := range sequence {
		out = append(out, v)
		break
	}
	assert.Equal(t, slice.Of(-1), out)

	out = nil
	for v := range sequence {
		out = append(out, v)
		if v == 2 {
			break
		}
	}
	assert.Equal(t, slice.Of(-1, 0, 1, 2), out)

	assert.Nil(t, seq.Slice(seq.Series(-1, (func(prev int) (int, bool))(nil))))
}

func Test_Append(t *testing.T) {
	in := slice.Of(1)
	out := seq.Append[seq.Seq[int]](nil, in)
	assert.Equal(t, in, out)

	out = seq.Append(seq.Of(2), in)
	assert.Equal(t, []int{1, 2}, out)
}

func Test_AccumSum(t *testing.T) {
	s := seq.Of(1, 3, 5, 7, 9, 11)
	r := seq.Accum(100, s, op.Sum[int])
	assert.Equal(t, 100+1+3+5+7+9+11, r)
	r = seq.Accum(100, s, (func(a int, b int) int)(nil))
	assert.Equal(t, 100, r)

	r = seq.Sum(s)
	assert.Equal(t, 1+3+5+7+9+11, r)
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
	reducer := func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errors.New("stop")
		}
		return i1 + i2, nil
	}
	r, ok, err := seq.ReduceeOK(s, reducer)
	assert.True(t, ok)
	assert.Equal(t, 1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")

	_, ok, err = seq.ReduceeOK[seq.Seq[int]](nil, reducer)
	assert.False(t, ok)
	assert.NoError(t, err)

	r, err = seq.Reducee(s, reducer)
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
	sequence := seq.Of(1, 2, 3, 4, 5, 6)
	result, ok := seq.First(sequence, more.Than(5))

	assert.True(t, ok)
	assert.Equal(t, 6, result)

	assert.True(t, seq.HasAny(sequence, more.Than(5)))

	_, ok = seq.First[seq.Seq[int]](nil, more.Than(5))
	assert.False(t, ok)

	_, ok = seq.First(sequence, nil)
	assert.False(t, ok)

	ok = seq.HasAny(sequence, more.Than(5))
	assert.True(t, ok)
}

func Test_Firstt(t *testing.T) {
	sequence := seq.Of(1, 2, 3, 4, 5, 6)
	result, ok, err := seq.Firstt(sequence, func(i int) (bool, error) {
		return more.Than(5)(i), nil
	})

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	result, ok, err = seq.Firstt(sequence, func(_ int) (bool, error) { return true, errors.New("abort") })

	assert.True(t, ok)
	assert.Equal(t, 1, result)
	assert.ErrorContains(t, err, "abort")

	result, ok, err = seq.Firstt(sequence, func(_ int) (bool, error) { return false, errors.New("abort") })

	assert.False(t, ok)
	assert.Equal(t, 0, result)
	assert.ErrorContains(t, err, "abort")

	_, ok, _ = seq.Firstt(sequence, nil)
	assert.False(t, ok)

	_, ok, _ = seq.Firstt[seq.Seq[int]](nil, func(_ int) (bool, error) { return false, errors.New("abort") })
	assert.False(t, ok)

}

var even = func(v int) bool { return v%2 == 0 }

func Test_Flat(t *testing.T) {
	md := seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := seq.Flat(md, as.Is)
	s := seq.Slice(f)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s)

	var out []int
	for v := range f {
		out = append(out, v)
		if v == 5 {
			break
		}
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5}, out)
}

func Test_FlatSeq(t *testing.T) {
	md := seq.Of([][]int{{1, 2, 3}, {4}, {5}, {6}}...)
	f := seq.FlatSeq(md, slices.Values)
	s := seq.Slice(f)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s)

	var out []int
	for v := range f {
		out = append(out, v)
		if v == 5 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, out)

	out = nil
	for i := range seq.FlatSeq(md, func(i []int) seq.Seq[int] {
		if len(i) == 1 {
			return nil
		}
		return seq.Of(i[0:1]...)
	}) {
		out = append(out, i)
	}
	assert.Equal(t, []int{1}, out)

	out = nil
	for v := range f {
		out = append(out, v)
		if v == 5 {
			break
		}
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5}, out)
}

func Test_Flatt(t *testing.T) {
	var (
		input     iter.Seq[[]string]
		flattener func([]string) ([]int, error)
	)
	for _, err := range seq.Flatt(input, flattener) {
		if err != nil {
			panic(err)
		}
	}

	s := seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5"}, {"6"}}...)
	f := func(strInteger []string) ([]int, error) { return slice.Conv(strInteger, strconv.Atoi) }
	out, err := seqe.Slice(seq.Flatt(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, out)
	assert.ErrorContains(t, err, "parsing \"_5\"")

	out = nil
	for v, err := range seq.Flatt(s, f) {
		if err != nil {
			panic(err)
		}
		out = append(out, v)
		if v == 4 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3, 4}, out)

	out = nil
	for v, err := range seq.Flatt(s, f) {
		_ = err
		out = append(out, v)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 0, 6}, out)
}

func Test_FlattSeq(t *testing.T) {
	var (
		input     iter.Seq[[]string]
		flattener func([]string) seq.SeqE[int]
	)
	for i, err := range seq.FlattSeq(input, flattener) {
		if err != nil {
			panic(err)
		}
		_ = i
	}

	s := seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5"}, {"6"}}...)
	f := func(strInteger []string) seq.SeqE[int] { return seq.Conv(seq.Of(strInteger...), strconv.Atoi) }
	i, err := seqe.Slice(seq.FlattSeq(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, i)
	assert.ErrorContains(t, err, "parsing \"_5\"")

	var out []int
	for v, err := range seq.FlattSeq(s, f) {
		if err != nil {
			panic(err)
		}
		out = append(out, v)
		if v == 4 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3, 4}, i)

	out = nil
	for v, err := range seq.FlattSeq(s, f) {
		_ = err
		out = append(out, v)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 0, 6}, out)
}

func Test_Filter(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := seq.Filter(s, even)
	assert.Equal(t, slice.Of(4, 8), seq.Slice(r))
}

func Test_Filt(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	filter := func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errors.New("abort"), nil) }
	l := seq.Filt(s, filter)
	r, err := seqe.Slice(l)

	assert.Error(t, err)
	assert.Equal(t, slice.Of(4), r)

	l = seq.Filt(s, nil)
	r, err = seqe.Slice(l)

	assert.NoError(t, err)
	assert.Empty(t, r)
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

type Rows[T any] struct {
	row    []T
	cursor int
}

func (r *Rows[T]) Reset()       { r.cursor = 0 }
func (r *Rows[T]) Next() bool   { return r.cursor < len(r.row) }
func (r *Rows[T]) Scan(dest *T) { *dest = r.row[r.cursor]; r.cursor++ }

func Test_OfNextPush(t *testing.T) {
	rows := &Rows[int]{slice.Of(1, 2, 3), 0}

	result := seq.Slice(seq.OfNextPush(rows.Next, rows.Scan))
	assert.Equal(t, slice.Of(1, 2, 3), result)

	rows.Reset()

	result = seq.Slice(seq.OfSourceNextPush(rows, (*Rows[int]).Next, func(r *Rows[int], out *int) { r.Scan(out) }))
	assert.Equal(t, slice.Of(1, 2, 3), result)
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
	keys := seq.Slice(seq2.Keys(s2))
	vals := seq.Slice(seq2.Values(s2))

	assert.Equal(t, slice.Of(2, 2, 1), keys)
	assert.Equal(t, slice.Of(1, 2, 3), vals)

	keys = nil
	vals = nil
	for k, v := range s2 {
		if v == 3 {
			break
		}
		keys = append(keys, k)
		vals = append(vals, v)
	}
	assert.Equal(t, slice.Of(2, 2), keys)
	assert.Equal(t, slice.Of(1, 2), vals)
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

	filter := func(e int) bool { return e%2 == 0 }
	for e := range seq.Filter(from, filter) {
		s = append(s, e)
		if e == 11 {
			break
		}
	}

	assert.Equal(t, slice.Of(2, 8), sort.Asc(s))
	assert.Equal(t, slice.Of(2, 8), sort.Asc(seq.Slice(seq.Filter(from, filter))))

	assert.Empty(t, seq.Slice(seq.Filter(from, nil)))
	assert.Empty(t, seq.Slice(seq.Filter[seq.Seq[int]](nil, filter)))
}

func Test_AllConverted(t *testing.T) {
	from := seq.Of(1, 2, 3, 5, 7, 8, 9, 11)
	s := []string{}

	for e := range seq.Convert(from, strconv.Itoa) {
		s = append(s, e)
	}

	expected := slice.Of("1", "2", "3", "5", "7", "8", "9", "11")
	assert.Equal(t, expected, s)
	assert.Equal(t, expected, seq.Slice(seq.Convert(from, strconv.Itoa)))

	assert.Empty(t, seq.Slice(seq.Convert[seq.Seq[int], int, int](from, nil)))
	assert.Empty(t, seq.Slice(seq.Convert[seq.Seq[int]](nil, strconv.Itoa)))
}

func Test_Conv(t *testing.T) {
	from := seq.Of("1", "2", "3", "5", "_7", "8", "9", "11")
	i := []int{}

	for v, err := range seq.Conv(from, strconv.Atoi) {
		if err == nil {
			i = append(i, v)
		}
		if v == 11 {
			break
		}
	}

	assert.Equal(t, slice.Of(1, 2, 3, 5, 8, 9, 11), i)
}

func Test_ConvertOK(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	converter := func(i int) (string, bool) { return strconv.Itoa(i), even(i) }
	r := seq.ConvertOK(s, converter)
	assert.Equal(t, []string{"4", "8"}, seq.Slice(r))
	r = seq.ConvertOK(iter.Seq[int](nil), converter)
	assert.Empty(t, seq.Slice(r))
	r = seq.ConvertOK(s, (func(i int) (string, bool))(nil))
	assert.Empty(t, seq.Slice(r))
}

func Test_ConvOK(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	converter := func(i int) (string, bool, error) { return strconv.Itoa(i), even(i), nil }
	r := seq.ConvOK(s, converter)
	o, err := seqe.Slice(r)
	assert.Equal(t, []string{"4", "8"}, o)
	assert.NoError(t, err)

	r = seq.ConvOK(iter.Seq[int](nil), converter)
	o, _ = seqe.Slice(r)
	assert.Empty(t, o)

	r = seq.ConvOK(s, (func(i int) (string, bool, error))(nil))
	o, _ = seqe.Slice(r)
	assert.Empty(t, o)
}

func Test_Range(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), seq.Slice(seq.Range(-1, 4)))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), seq.Slice(seq.Range(3, -2)))
	assert.Nil(t, seq.Slice(seq.Range(1, 1)))

	var out []int
	for i := range seq.Range(-1, 3) {
		if i == 2 {
			break
		}
		out = append(out, i)
	}
	assert.Equal(t, slice.Of(-1, 0, 1), out)
}

func Test_RangeClosed(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), seq.Slice(seq.RangeClosed(-1, 3)))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), seq.Slice(seq.RangeClosed(3, -1)))
	assert.Equal(t, slice.Of(1), seq.Slice(seq.RangeClosed(1, 1)))

	var out []int
	for i := range seq.RangeClosed(-1, 3) {
		if i == 2 {
			break
		}
		out = append(out, i)
	}
	assert.Equal(t, slice.Of(-1, 0, 1), out)
}

func Test_ToSeq2(t *testing.T) {
	s, err := seqe.Slice(seq.ToSeq2[seq.Seq[int], int, int, error](seq.Of(1), nil))

	assert.NoError(t, err)
	assert.Empty(t, s)

	s, err = seqe.Slice(seq.ToSeq2[seq.Seq[int]](nil, func(_ int) (int, error) { return 0, errors.New("abort") }))

	assert.NoError(t, err)
	assert.Empty(t, s)
}

func Test_Group(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = seq.Group(seq.Of(1, 1, 2, 4, 3, 1), even, as.Is)
	)
	assert.Equal(t, slice.Of(2, 4), groups[true])
	assert.Equal(t, slice.Of(1, 1, 3, 1), groups[false])
}
