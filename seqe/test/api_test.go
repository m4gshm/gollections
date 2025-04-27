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

func errIfContains[T comparable](errVal T) func([]T) ([]T, error) {
	return func(val []T) ([]T, error) {
		if slice.Contains(val, errVal) {
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

	r, _ = seqe.Sum(s)
	assert.Equal(t, 1+3+5+7+9+11, r)
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

	sum, err = seqe.Reduce(s, op.Sum)
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

	r, err = seqe.Reducee(s, func(i1, i2 int) (int, error) { return i1 + i2, nil })
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
	sequence := seq.Of(1, 2, 3, 4, 5, 6)
	result, ok, err := seqe.First(seq.ToSeq2(sequence, noErr), more.Than(5))

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	_, ok, err = seqe.First(seq.ToSeq2(sequence, errOn(1)), more.Than(5))

	assert.False(t, ok)
	assert.ErrorContains(t, err, "abort")

	ok, err = seqe.HasAny(seq.ToSeq2(sequence, noErr), more.Than(5))
	assert.True(t, ok)
	assert.NoError(t, err)
}

func Test_Firstt(t *testing.T) {
	sequence := seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr)
	result, ok, err := seqe.Firstt(sequence, func(i int) (bool, error) {
		return more.Than(5)(i), nil
	})

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	result, ok, err = seqe.Firstt(sequence, func(_ int) (bool, error) { return true, errors.New("abort") })

	assert.True(t, ok)
	assert.Equal(t, 1, result)
	assert.ErrorContains(t, err, "abort")

	result, ok, err = seqe.Firstt(sequence, func(_ int) (bool, error) { return false, errors.New("abort") })

	assert.False(t, ok)
	assert.Equal(t, 0, result)
	assert.ErrorContains(t, err, "abort")

	sequence = seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), errOn(1))
	result, ok, err = seqe.Firstt(sequence, func(i int) (bool, error) { return more.Than(5)(i), nil })

	assert.False(t, ok)
	assert.Equal(t, 0, result)
	assert.ErrorContains(t, err, "abort")

	_, ok, _ = seqe.Firstt(sequence, nil)
	assert.False(t, ok)

	_, ok, _ = seqe.Firstt[seq.SeqE[int]](nil, func(_ int) (bool, error) { return false, errors.New("abort") })
	assert.False(t, ok)
}

var even = func(v int) bool { return v%2 == 0 }

func Test_Flat(t *testing.T) {
	md := seq.ToSeq2(seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...), noErr)
	f := seqe.Flat(md, as.Is)
	s, err := seqe.Slice(f)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s)
	assert.NoError(t, err)

	var out []int
	for v, err := range f {
		if err != nil {
			panic(err)
		}
		out = append(out, v)
		if v == 5 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, out)

	md = seq.ToSeq2(seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...), errIfContains(5))
	f = seqe.Flat(md, as.Is)
	s, err = seqe.Slice(f)
	assert.Equal(t, []int{1, 2, 3, 4}, s)
	assert.ErrorContains(t, err, "abort")
}

func Test_FlatSeq(t *testing.T) {
	md := seq.ToSeq2(seq.Of([][]int{{1, 2, 3}, {4}, {5}, {6}}...), noErr)
	f := seqe.FlatSeq(md, slices.Values)
	s, err := seqe.Slice(f)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s)
	assert.NoError(t, err)

	var out []int
	for v, err := range f {
		if err != nil {
			panic(err)
		}
		out = append(out, v)
		if v == 5 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, out)

	md = seq.ToSeq2(seq.Of([][]int{{1, 2, 3}, {4}, {5}, {6}}...), errIfContains(5))
	f = seqe.FlatSeq(md, slices.Values)
	s, err = seqe.Slice(f)
	assert.Equal(t, []int{1, 2, 3, 4}, s)
	assert.ErrorContains(t, err, "abort")

	out = nil
	for i, err := range seqe.FlatSeq(md, func(i []int) seq.Seq[int] {
		if len(i) == 1 {
			return nil
		}
		return seq.Of(i[0:1]...)
	}) {
		if err == nil {
			out = append(out, i)
		}
	}
	assert.Equal(t, []int{1}, out)
}

func Test_Flatt(t *testing.T) {
	var (
		input     iter.Seq2[[]string, error]
		flattener func([]string) ([]int, error)
	)

	for _, err := range seqe.Flatt(input, flattener) {
		if err != nil {
			panic(err)
		}
	}

	s := seq.ToSeq2(seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5"}, {"6"}}...), noErr)
	f := func(strInteger []string) ([]int, error) { return slice.Conv(strInteger, strconv.Atoi) }
	out, err := seqe.Slice(seqe.Flatt(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, out)
	assert.ErrorContains(t, err, "parsing \"_5\"")

	s = seq.ToSeq2(seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5"}, {"6"}}...), errIfContains("_5"))
	out, err = seqe.Slice(seqe.Flatt(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, out)
	assert.ErrorContains(t, err, "abort")

	out = nil
	for v, err := range seqe.Flatt(s, f) {
		_ = err
		out = append(out, v)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 0, 6}, out)
}

func Test_FlattSeq(t *testing.T) {
	var (
		input     iter.Seq2[[]string, error]
		flattener func([]string) seq.SeqE[int]
	)
	for i, err := range seqe.FlattSeq(input, flattener) {
		if err != nil {
			panic(err)
		}
		_ = i
	}

	s := seq.ToSeq2(seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5"}, {"6"}}...), noErr)
	f := func(strInteger []string) seq.SeqE[int] { return seq.Conv(seq.Of(strInteger...), strconv.Atoi) }
	i, err := seqe.Slice(seqe.FlattSeq(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, i)
	assert.ErrorContains(t, err, "parsing \"_5\"")

	s = seq.ToSeq2(seq.Of([][]string{{"1", "2", "3"}, {"4"}, {"_5"}, {"6"}}...), func(s []string) ([]string, error) {
		return s, op.IfElse(slice.Contains(s, "_5"), errors.New("abort"), nil)
	})
	i, err = seqe.Slice(seqe.FlattSeq(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, i)
	assert.ErrorContains(t, err, "abort")

	var out []int
	for v, err := range seqe.FlattSeq(s, f) {
		_ = err
		out = append(out, v)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 0, 6}, out)
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
	assert.Equal(t, slice.Of(4), r)

	l = seqe.Filt(s, nil)
	r, err = seqe.Slice(l)
	assert.NoError(t, err)
	assert.Empty(t, r)

	l = seqe.Filt[seq.SeqE[int]](nil, filter)
	r, err = seqe.Slice(l)
	assert.NoError(t, err)
	assert.Empty(t, r)

	s = seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), errOn(4))
	l = seqe.Filt(s, filter)
	r, err = seqe.Slice(l)
	assert.Error(t, err)
	assert.Nil(t, r)
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

type Rows[T any] struct {
	row    []T
	cursor int
}

func (r *Rows[T]) Reset()       { r.cursor = 0 }
func (r *Rows[T]) Next() bool   { return r.cursor < len(r.row) }
func (r *Rows[T]) Scan(dest *T) error { *dest = r.row[r.cursor]; r.cursor++; return nil }

func Test_OfNextPush(t *testing.T) {
	rows := &Rows[int]{slice.Of(1, 2, 3), 0}

	result, err := seqe.Slice(seqe.OfNextPush(rows.Next, rows.Scan))
	assert.NoError(t, err)
	assert.Equal(t, slice.Of(1, 2, 3), result)

	rows.Reset()

	result, err = seqe.Slice(seqe.OfSourceNextPush(rows, (*Rows[int]).Next, func(r *Rows[int], out *int) error { return r.Scan(out) }))
	assert.NoError(t, err)
	assert.Equal(t, slice.Of(1, 2, 3), result)
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

	from = seq.ToSeq2(seq.Of(1, 2, 3, 5, 7, 8, 9, 11), errOn((7)))

	s = nil
	for e, err := range seqe.Convert(from, strconv.Itoa) {
		//ignore err
		_ = err
		s = append(s, e)
	}

	assert.Equal(t, slice.Of("1", "2", "3", "5", "", "8", "9", "11"), s)
}

func Test_Conv(t *testing.T) {
	from := seq.ToSeq2(seq.Of("1", "2", "3", "5", "_7", "8", "9", "11"), noErr)
	out := []int{}

	for v, err := range seqe.Conv(from, strconv.Atoi) {
		if err == nil {
			out = append(out, v)
		}
	}

	assert.Equal(t, slice.Of(1, 2, 3, 5, 8, 9, 11), out)

	from = seq.ToSeq2(seq.Of("1", "2", "3", "5", "_7", "8", "9", "11"), errOn("_7"))
	out = nil
	for v, err := range seqe.Conv(from, strconv.Atoi) {
		if err != nil {
			break
		}
		out = append(out, v)
	}
	assert.Equal(t, slice.Of(1, 2, 3, 5), out)
}

func Test_ConvertOK(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), noErr)
	converter := func(i int) (string, bool) { return strconv.Itoa(i), even(i) }
	r := seqe.ConvertOK(s, converter)
	out, err := seqe.Slice(r)
	assert.Equal(t, []string{"4", "8"}, out)
	assert.NoError(t, err)

	r = seqe.ConvertOK(seq.SeqE[int](nil), converter)
	out, err = seqe.Slice(r)
	assert.NoError(t, err)
	assert.Empty(t, out)
	r = seqe.ConvertOK(s, (func(i int) (string, bool))(nil))
	out, err = seqe.Slice(r)
	assert.NoError(t, err)
	assert.Empty(t, out)

	s = seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11, 12), errOn(9))
	r = seqe.ConvertOK(s, converter)
	out, err = seqe.Slice(r)
	assert.Equal(t, []string{"4", "8"}, out)
	assert.ErrorContains(t, err, "abort")

	out = nil
	for s, err := range r {
		_ = err
		out = append(out, s)
	}
	assert.Equal(t, []string{"4", "8", "", "12"}, out)
}

func Test_ConvOK(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), noErr)
	converter := func(i int) (string, bool, error) { return strconv.Itoa(i), even(i), nil }
	r := seqe.ConvOK(s, converter)
	o, err := seqe.Slice(r)
	assert.Equal(t, []string{"4", "8"}, o)
	assert.NoError(t, err)

	r = seqe.ConvOK(seq.SeqE[int](nil), converter)
	o, _ = seqe.Slice(r)
	assert.Empty(t, o)

	r = seqe.ConvOK(s, (func(i int) (string, bool, error))(nil))
	o, _ = seqe.Slice(r)
	assert.Empty(t, o)

	r = seqe.ConvOK(s, func(i int) (string, bool, error) {
		return strconv.Itoa(i), even(i), op.IfElse(i == 9, errors.New("abort"), nil)
	})
	o, err = seqe.Slice(r)

	assert.Equal(t, []string{"4", "8"}, o)
	assert.ErrorContains(t, err, "abort")

	s = seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), errOn(5))
	r = seqe.ConvOK(s, converter)
	o, err = seqe.Slice(r)
	assert.Error(t, err)
	assert.Equal(t, []string{"4"}, o)
}

func Test_Group(t *testing.T) {
	even := func(v int) bool { return v%2 == 0 }

	groups, err := seqe.Group(seq.ToSeq2(seq.Of(1, 1, 2, 4, 3, 5), errOn(5)), even, as.Is)
	assert.Equal(t, slice.Of(2, 4), groups[true])
	assert.Equal(t, slice.Of(1, 1, 3), groups[false])
	assert.Error(t, err)

	groups, err = seqe.Group(seq.ToSeq2(seq.Of(1, 1, 2, 4, 3, 5), noErr), even, as.Is)
	assert.Equal(t, slice.Of(2, 4), groups[true])
	assert.Equal(t, slice.Of(1, 1, 3, 5), groups[false])
	assert.NoError(t, err)
}
