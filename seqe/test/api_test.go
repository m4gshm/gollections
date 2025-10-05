package test

import (
	"errors"
	"iter"
	"slices"
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/predicate/not"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seqe"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

var errStop = errors.New("stop")
var even = func(v int) bool { return v%2 == 0 }

func noErr[T any](t T) (T, error) { return t, nil }
func errOn[T comparable](errVal T) func(T) (T, error) {
	return func(val T) (T, error) {
		if val == errVal {
			return val, errStop
		}
		return val, nil
	}
}

func errIfContains[T comparable](errVal T) func([]T) ([]T, error) {
	return func(val []T) ([]T, error) {
		if slice.Contains(val, errVal) {
			return val, errStop
		}
		return val, nil
	}
}

func Test_Union(t *testing.T) {
	sequence := seqe.Union(seq.ToSeq2(seq.Of(0, 1), noErr), nil, seq.ToSeq2(seq.Of[int](), errOn(0)), seq.ToSeq2(seq.Of(2, 3, 4), errOn(3)))
	result, err := seqe.Slice(sequence)
	assert.Equal(t, slice.Of(0, 1, 2), result)
	assert.Error(t, err)

	sequence2 := seq.Conv(seq.Of(0, 1), noErr).Union(nil, seq.ToSeq2(seq.Of[int](), errOn(0)).Union(seq.ToSeq2(seq.Of(2, 3, 4), errOn(3))))
	result, err = seqe.Slice(sequence2)
	assert.Equal(t, slice.Of(0, 1, 2), result)
	assert.Error(t, err)
}

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")
	result := seqe.OfIndexed(len(indexed), func(i int) (string, error) { return indexed[i], nil })
	out, err := seqe.Slice(result)

	assert.Equal(t, indexed, out)
	assert.NoError(t, err)

	result = seqe.OfIndexed(len(indexed), func(i int) (string, error) {
		return indexed[i], op.IfElse(i == 3, errStop, nil)
	})
	out, err = seqe.Slice(result)

	assert.Equal(t, slice.Of("0", "1", "2"), out)
	assert.ErrorContains(t, err, "stop")
}

func Test_Append(t *testing.T) {
	in := slice.Of(1)
	out, err := seqe.Append[seq.SeqE[int]](nil, in)
	assert.Equal(t, in, out)
	assert.NoError(t, err)

	out, err = seqe.Append(seq.ToSeq2(seq.Of(2), noErr), out)
	assert.Equal(t, []int{1, 2}, out)
	assert.NoError(t, err)

	out, err = seq.Of(3).Conv(noErr).Append(out)
	assert.Equal(t, []int{1, 2, 3}, out)
	assert.NoError(t, err)

	out, err = seqe.Append(seq.Of(4, 5, 6).Conv(errOn(5)), out)
	assert.Equal(t, []int{1, 2, 3, 4}, out)
	assert.Error(t, err)

	out, err = seq.Of(7, 8, 9).Conv(errOn(8)).Append(out)
	assert.Equal(t, []int{1, 2, 3, 4, 7}, out)
	assert.Error(t, err)
}

func Test_AccumSum(t *testing.T) {
	s := seq.Conv(seq.Of(1, 3, 5, 7, 9, 11), noErr)
	r, err := seqe.Accum(100, s, op.Sum[int])
	assert.Equal(t, 100+1+3+5+7+9+11, r)
	assert.NoError(t, err)

	r, _ = seqe.Sum(s)
	assert.Equal(t, 1+3+5+7+9+11, r)

	r, _ = s.Accum(100, op.Sum[int])
	assert.Equal(t, 100+1+3+5+7+9+11, r)
}

func Test_AccummSum(t *testing.T) {
	s := seq.Conv(seq.Of(1, 3, 5, 7, 9, 11), noErr)
	adder := func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errStop
		}
		return i1 + i2, nil
	}
	r, err := seqe.Accumm(100, s, adder)
	assert.Equal(t, 100+1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")

	r, err = s.Accumm(100, adder)
	assert.Equal(t, 100+1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceSum(t *testing.T) {
	ns := seq.Of(1, 2, 3, 4, 5, 6)
	s2 := seq.ToSeq2(ns, noErr)
	sum, ok, err := seqe.ReduceOK(s2, op.Sum)

	assert.True(t, ok)
	assert.Equal(t, 21, sum)
	assert.NoError(t, err)

	sum, err = seqe.Reduce(s2, op.Sum)
	assert.Equal(t, 21, sum)
	assert.NoError(t, err)

	sum, err = seq.Conv(ns, noErr).Reduce(op.Sum)
	assert.Equal(t, 21, sum)
	assert.NoError(t, err)

	s2 = seq.ToSeq2(ns, errOn(4))
	sum, ok, err = seqe.ReduceOK(s2, op.Sum)

	assert.True(t, ok)
	assert.Equal(t, 6, sum)
	assert.ErrorContains(t, err, "stop")

	sum, ok, err = ns.Conv(errOn(4)).ReduceOK(op.Sum)

	assert.True(t, ok)
	assert.Equal(t, 6, sum)
	assert.ErrorContains(t, err, "stop")

	s2 = seq.ToSeq2(ns, errOn(1))
	sum, ok, err = seqe.ReduceOK(s2, op.Sum)

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
	assert.ErrorContains(t, err, "stop")

	se := ns.Conv(errOn(1))
	sum, ok, err = se.ReduceOK(op.Sum)

	assert.False(t, ok)
	assert.Equal(t, 0, sum)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceeSum(t *testing.T) {
	s := seq.Of(1, 3, 5, 7, 9, 11)
	s2 := seq.ToSeq2(s, noErr)
	adderErr := func(i1, i2 int) (int, error) {
		if i2 == 11 {
			return i1, errStop
		}
		return i1 + i2, nil
	}
	r, ok, err := seqe.ReduceeOK(s2, adderErr)
	assert.True(t, ok)
	assert.Equal(t, 1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")

	r, ok, err = s.Conv(noErr).ReduceeOK(adderErr)
	assert.True(t, ok)
	assert.Equal(t, 1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")

	s2 = seq.ToSeq2(s, errOn(5))
	adder := func(i1, i2 int) (int, error) { return i1 + i2, nil }
	r, ok, err = seqe.ReduceeOK(s2, adder)

	assert.True(t, ok)
	assert.Equal(t, 1+3, r)
	assert.ErrorContains(t, err, "stop")

	r, err = seqe.Reducee(s2, adder)
	assert.Equal(t, 1+3, r)
	assert.ErrorContains(t, err, "stop")

	r, err = s.Conv(errOn(5)).Reducee(adder)
	assert.Equal(t, 1+3, r)
	assert.ErrorContains(t, err, "stop")

	s2 = seq.ToSeq2(s, errOn(1))
	r, ok, err = seqe.ReduceeOK(s2, adder)

	assert.False(t, ok)
	assert.Equal(t, 0, r)
	assert.ErrorContains(t, err, "stop")

	r, ok, err = seq.Conv(s, errOn(1)).ReduceeOK(adder)
	assert.False(t, ok)
	assert.Equal(t, 0, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceeSumFirstErr(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 5, 7, 9, 11), noErr)
	r, ok, err := seqe.ReduceeOK(s, func(_, _ int) (int, error) {
		return 0, errStop
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

func Test_Head(t *testing.T) {
	sequence := seq.Conv(seq.Of(1, 2, 3, 4, 5, 6), noErr)
	result, ok, err := seqe.Head(sequence)

	assert.True(t, ok)
	assert.Equal(t, 1, result)
	assert.NoError(t, err)

	result, ok, err = sequence.Head()

	assert.True(t, ok)
	assert.Equal(t, 1, result)
	assert.NoError(t, err)

	result, ok, err = seqe.Head[seq.SeqE[int]](nil)
	assert.Zero(t, result)
	assert.False(t, ok)
	assert.NoError(t, err)
}

func Test_While(t *testing.T) {
	sequence := seq.SeqE[int](seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr))
	part := seqe.While(sequence, not.Eq(5))

	s, err := seqe.Slice(part)
	assert.Equal(t, slice.Of(1, 2, 3, 4), s)
	assert.NoError(t, err)

	part = sequence.While(not.Eq(5))

	s, err = seqe.Slice(part)
	assert.Equal(t, slice.Of(1, 2, 3, 4), s)
	assert.NoError(t, err)

	part = seqe.While(sequence, not.Eq(7))
	s, err = seqe.Slice(part)
	assert.Equal(t, slice.Of(1, 2, 3, 4, 5, 6), s)
	assert.NoError(t, err)

	part = seqe.While(sequence, eq.To(0))
	s, err = seqe.Slice(part)
	assert.Nil(t, s)
	assert.NoError(t, err)

	part = seqe.While[seq.SeqE[int]](nil, eq.To(0))
	s, err = seqe.Slice(part)
	assert.Nil(t, s)
	assert.NoError(t, err)

	sequence = seq.SeqE[int](seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), errOn(5)))
	r := []int{}
	for i, err := range sequence.While(not.Eq(7)) {
		if err != nil {
			break
		}
		r = append(r, i)
	}
	assert.Equal(t, slice.Of(1, 2, 3, 4), r)
}

func Test_SkipWhile(t *testing.T) {
	sequence := seq.SeqE[int](seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr))
	part := seqe.SkipWhile(sequence, less.Than(4))

	s, err := seqe.Slice(part)
	assert.Equal(t, slice.Of(4, 5, 6), s)
	assert.NoError(t, err)

	part = sequence.SkipWhile(less.Than(4))

	s, err = seqe.Slice(part)
	assert.Equal(t, slice.Of(4, 5, 6), s)
	assert.NoError(t, err)

	part = seqe.SkipWhile(sequence, not.Eq(7))
	s, err = seqe.Slice(part)
	assert.Nil(t, s)
	assert.NoError(t, err)

	part = seqe.SkipWhile(sequence, less.Than(0))
	s, err = seqe.Slice(part)
	assert.Equal(t, slice.Of(1, 2, 3, 4, 5, 6), s)
	assert.NoError(t, err)

	r := []int{}
	sequence = seq.SeqE[int](seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), errOn(6)))
	for i, err := range sequence.SkipWhile(less.Than(4)) {
		if err != nil {
			break
		}
		r = append(r, i)
	}
	assert.Equal(t, slice.Of(4, 5), r)
}

func Test_Top(t *testing.T) {
	sequence := seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr)
	top := seqe.Top(4, sequence)
	result, err := seqe.Slice(top)

	assert.Equal(t, slice.Of(1, 2, 3, 4), result)
	assert.NoError(t, err)

	result, err = seqe.Slice(seqe.Top(0, sequence))
	assert.Nil(t, result)
	assert.NoError(t, err)

	result, err = seqe.Slice(seqe.Top[seq.SeqE[int]](10, nil))
	assert.Nil(t, result)
	assert.NoError(t, err)

	result = nil
	for v, err := range seqe.Top(4, seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), errOn(3))) {
		if v != 3 {
			result = append(result, v)
			assert.NoError(t, err, "unexpected error on %i", v)
		} else {
			assert.Error(t, err, "expected error on %i", v)
			break
		}
	}
	assert.Equal(t, slice.Of(1, 2), result)

	result = nil
	for v := range seq.SeqE[int](seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), errOn(2))).Top(4) {
		if v != 2 {
			result = append(result, v)
		}
	}
	assert.Equal(t, slice.Of(1, 3, 4), result)
}

func Test_Skip(t *testing.T) {
	s := seq.Of(1, 2, 3, 4, 5, 6)
	s2 := seq.ToSeq2(s, noErr)
	skip := seqe.Skip(4, s2)
	result, err := seqe.Slice(skip)

	assert.Equal(t, slice.Of(5, 6), result)
	assert.NoError(t, err)

	result, err = seqe.Slice(seqe.Skip(0, s2))
	assert.Equal(t, slice.Of(1, 2, 3, 4, 5, 6), result)
	assert.NoError(t, err)

	result, err = seqe.Slice(seqe.Skip[seq.SeqE[int]](10, nil))
	assert.Nil(t, result)
	assert.NoError(t, err)

	result = nil
	for v, err := range seqe.Skip(2, seq.ToSeq2(s, errOn(5))) {
		if v != 5 {
			result = append(result, v)
			assert.NoError(t, err, "unexpected error on %i", v)
		} else {
			assert.Error(t, err, "expected error on %i", v)
			break
		}
	}
	assert.Equal(t, slice.Of(3, 4), result)

	result = nil
	for v := range seqe.Skip(2, seq.ToSeq2(s, errOn(5))) {
		if v != 5 {
			result = append(result, v)
		}
	}
	assert.Equal(t, slice.Of(3, 4, 6), result)

	result = nil
	for v := range seq.Conv(s, errOn(5)).Skip(2) {
		if v != 5 {
			result = append(result, v)
		}
	}
	assert.Equal(t, slice.Of(3, 4, 6), result)
}

func Test_SkipTop(t *testing.T) {
	sequence := seq.ToSeq2(seq.Of(1, 2, 3, 4, 5, 6), noErr)
	result, err := seqe.Slice(seqe.Top(2, seqe.Skip(2, sequence)))

	assert.Equal(t, slice.Of(3, 4), result)
	assert.NoError(t, err)
}

func Test_First(t *testing.T) {
	sequence := seq.Of(1, 2, 3, 4, 5, 6)
	result, ok, err := seqe.First(seq.Conv(sequence, noErr), more.Than(5))

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	result, ok, err = sequence.Conv(noErr).First(more.Than(5))

	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	_, ok, err = seqe.First(seq.Conv(sequence, errOn(1)), more.Than(5))
	assert.False(t, ok)
	assert.ErrorContains(t, err, "stop")

	_, ok, err = sequence.Conv(errOn(1)).First(more.Than(5))
	assert.False(t, ok)
	assert.ErrorContains(t, err, "stop")
}

func Test_HasAny(t *testing.T) {
	sequence := seq.Of(1, 2, 3, 4, 5, 6)
	ok, err := seqe.HasAny(seq.Conv(sequence, noErr), more.Than(5))
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = sequence.Conv(noErr).HasAny(more.Than(5))
	assert.True(t, ok)
	assert.NoError(t, err)
}

func Test_Firstt(t *testing.T) {
	firstErr := func(_ int) (bool, error) { return true, errStop }
	mor5NoErr := func(i int) (bool, error) { return more.Than(5)(i), nil }
	justErr := func(_ int) (bool, error) { return false, errStop }

	s := seq.Of(1, 2, 3, 4, 5, 6)
	s2 := seq.ToSeq2(s, noErr)
	result, ok, err := seqe.Firstt(s2, mor5NoErr)
	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	result, ok, err = seq.Conv(s, noErr).Firstt(mor5NoErr)
	assert.True(t, ok)
	assert.Equal(t, 6, result)
	assert.NoError(t, err)

	result, ok, err = seqe.Firstt(s2, firstErr)
	assert.True(t, ok)
	assert.Equal(t, 1, result)
	assert.ErrorContains(t, err, "stop")

	result, ok, err = seqe.Firstt(s2, justErr)

	assert.False(t, ok)
	assert.Equal(t, 0, result)
	assert.ErrorContains(t, err, "stop")

	s2 = seq.ToSeq2(s, errOn(1))
	result, ok, err = seqe.Firstt(s2, mor5NoErr)

	assert.False(t, ok)
	assert.Equal(t, 0, result)
	assert.ErrorContains(t, err, "stop")

	_, ok, _ = seqe.Firstt(s2, nil)
	assert.False(t, ok)

	_, ok, _ = seqe.Firstt[seq.SeqE[int]](nil, justErr)
	assert.False(t, ok)
}

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
	assert.ErrorContains(t, err, "stop")
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
	assert.ErrorContains(t, err, "stop")

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
	assert.ErrorContains(t, err, "stop")

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
		return s, op.IfElse(slice.Contains(s, "_5"), errStop, nil)
	})
	i, err = seqe.Slice(seqe.FlattSeq(s, f))

	assert.Equal(t, []int{1, 2, 3, 4}, i)
	assert.ErrorContains(t, err, "stop")

	var out []int
	for v, err := range seqe.FlattSeq(s, f) {
		_ = err
		out = append(out, v)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 0, 6}, out)
}

func Test_Filter(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)
	s2 := seq.ToSeq2(s, noErr)
	r, err := seqe.Filter(s2, even).Slice()
	assert.Equal(t, slice.Of(4, 8), r)
	assert.NoError(t, err)

	r, err = seq.Conv(s, noErr).Filter(even).Slice()
	assert.Equal(t, slice.Of(4, 8), r)
	assert.NoError(t, err)

}

func Test_Filt(t *testing.T) {
	s := seq.Of(1, 3, 4, 5, 7, 8, 9, 11)

	s2 := seq.ToSeq2(s, noErr)
	filter := func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errStop, nil) }
	r, err := seqe.Filt(s2, filter).Slice()
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4), r)

	se := seq.Conv(s, noErr)
	r, err = se.Filt(filter).Slice()
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4), r)

	r, err = seqe.Filt(s2, nil).Slice()
	assert.NoError(t, err)
	assert.Empty(t, r)

	r, err = se.Filt(nil).Slice()
	assert.NoError(t, err)
	assert.Empty(t, r)

	r, err = seqe.Filt[seq.SeqE[int]](nil, filter).Slice()
	assert.NoError(t, err)
	assert.Empty(t, r)

	s2 = seq.ToSeq2(s, errOn(4))
	r, err = seqe.Filt(s2, filter).Slice()
	assert.Error(t, err)
	assert.Nil(t, r)

	r, err = seq.Conv(s, errOn(4)).Filt(filter).Slice()
	assert.Error(t, err)
	assert.Nil(t, r)
}

func Test_Filt2(t *testing.T) {
	s := seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), noErr)
	l := seqe.Filt(s, func(i int) (bool, error) {
		ok := i <= 7
		return ok && even(i), op.IfElse(ok, nil, errStop)
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
	assert.ErrorContains(t, err, "stop")
}

type Rows[T any] struct {
	row    []T
	cursor int
}

func (r *Rows[T]) Reset()             { r.cursor = 0 }
func (r *Rows[T]) Next() bool         { return r.cursor < len(r.row) }
func (r *Rows[T]) Scan(dest *T) error { *dest = r.row[r.cursor]; r.cursor++; return nil }

func Test_OfNextPush(t *testing.T) {
	rows := &Rows[int]{slice.Of(1, 2, 3), 0}

	result, err := seqe.Slice(seqe.OfNext(rows.Next, rows.Scan))
	assert.NoError(t, err)
	assert.Equal(t, slice.Of(1, 2, 3), result)

	rows.Reset()

	result, err = seqe.Slice(seqe.OfSourceNext(rows, (*Rows[int]).Next, func(r *Rows[int], out *int) error { return r.Scan(out) }))
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

func Test_Convert(t *testing.T) {
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

func Test_ConvertNilSafe(t *testing.T) {
	type entity struct{ val *string }
	var (
		first    = "first"
		third    = "third"
		fifth    = "fifth"
		source   = seq.ToSeq2(seq.Of([]*entity{{&first}, {}, {&third}, nil, {&fifth}}...), noErr)
		result   = seqe.ConvertNilSafe(source, func(e *entity) *string { return e.val })
		expected = []*string{&first, &third, &fifth}
	)
	s, err := result.Slice()

	assert.NoError(t, err)
	assert.Equal(t, expected, s)
}

func Test_NotLit(t *testing.T) {
	var (
		first  = "first"
		third  = "third"
		fifth  = "fifth"
		source = seq.ToSeq2(seq.Of(&first, nil, &third, nil, &fifth), noErr)
		result = seqe.NotNil(source)
	)
	s, err := result.Slice()
	assert.Equal(t, slice.Of(&first, &third, &fifth), s)
	assert.NoError(t, err)
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
	assert.ErrorContains(t, err, "stop")

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
		return strconv.Itoa(i), even(i), op.IfElse(i == 9, errStop, nil)
	})
	o, err = seqe.Slice(r)

	assert.Equal(t, []string{"4", "8"}, o)
	assert.ErrorContains(t, err, "stop")

	s = seq.ToSeq2(seq.Of(1, 3, 4, 5, 7, 8, 9, 11), errOn(5))
	r = seqe.ConvOK(s, converter)
	o, err = seqe.Slice(r)
	assert.Error(t, err)
	assert.Equal(t, []string{"4"}, o)
}

func Test_Group(t *testing.T) {
	groups, err := seqe.Group(seq.ToSeq2(seq.Of(1, 1, 2, 4, 3, 5), errOn(5)), even, as.Is)
	assert.Equal(t, slice.Of(2, 4), groups[true])
	assert.Equal(t, slice.Of(1, 1, 3), groups[false])
	assert.Error(t, err)

	groups, err = seqe.Group(seq.ToSeq2(seq.Of(1, 1, 2, 4, 3, 5), noErr), even, as.Is)
	assert.Equal(t, slice.Of(2, 4), groups[true])
	assert.Equal(t, slice.Of(1, 1, 3, 5), groups[false])
	assert.NoError(t, err)
}

func Test_TrackEach(t *testing.T) {
	var out []int
	s2 := seq.ToSeq2(seq.RangeClosed(-1, 3), errOn(2))
	seqe.ForEach(s2, func(v int) { out = append(out, v) })
	assert.Equal(t, slice.Of(-1, 0, 1), out)

	out = nil
	seq.Conv(seq.RangeClosed(-1, 3), errOn(2)).ForEach(func(v int) { out = append(out, v) })
	assert.Equal(t, slice.Of(-1, 0, 1), out)
}
