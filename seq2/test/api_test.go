package test

import (
	"errors"
	"iter"
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/ordered/map_"
	"github.com/m4gshm/gollections/k"
	kvpredicate "github.com/m4gshm/gollections/kv/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/predicate/not"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/stretchr/testify/assert"
)

func Test_Of(t *testing.T) {
	sequence := seq2.Of(0, 1, 2, 3, 4)
	var out []int
	var ind []int
	for i, v := range sequence {
		out = append(out, v)
		ind = append(ind, i)
	}
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), out)
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), ind)
	out = nil
	for _, v := range sequence {
		if v == 1 {
			break
		}
		out = append(out, v)
	}
	assert.Equal(t, slice.Of(0), out)

	out = nil
	var iter = false
	for _, v := range sequence {
		iter = true
		_ = v
		break
	}
	assert.True(t, iter)
	assert.Nil(t, out)
}
func Test_Union(t *testing.T) {
	sequence := seq2.Union(seq2.Of(0, 1), nil, seq2.Of[int](), seq2.Of(2, 3, 4))
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), seq.Slice(seq2.Values(sequence)))
	assert.Equal(t, slice.Of(0, 1, 0, 1, 2), seq.Slice(seq2.Keys(sequence)))

	r := []int{}
	for _, v := range sequence {
		if v == 4 {
			break
		}
		r = append(r, v)
	}
	assert.Equal(t, slice.Of(0, 1, 2, 3), r)
}

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")
	getAt := func(i int) string { return indexed[i] }
	sequence := seq2.OfIndexed(len(indexed), getAt)
	assert.Equal(t, indexed, seq.Slice(seq2.Values(sequence)))

	var out []string
	var iter = false
	for _, v := range sequence {
		iter = true
		if v == "3" {
			break
		}
		out = append(out, v)
	}
	assert.True(t, iter)
	assert.Equal(t, slice.Of("0", "1", "2"), out)

	sequence = seq2.OfIndexed(len(indexed), (func(i int) string)(nil))
	assert.Nil(t, seq.Slice(seq2.Values(sequence)))
}

func Test_Series(t *testing.T) {
	generator := func(_, prev int) (int, bool) { return prev + 1, prev < 3 }
	sequence := seq2.Series(-1, generator)
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), seq.Slice(seq2.Values(sequence)))
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), seq.Slice(seq2.Keys(sequence)))

	var out []int
	for _, v := range sequence {
		out = append(out, v)
		break
	}
	assert.Equal(t, slice.Of(-1), out)

	out = nil
	for _, v := range sequence {
		out = append(out, v)
		if v == 2 {
			break
		}
	}
	assert.Equal(t, slice.Of(-1, 0, 1, 2), out)

	assert.Nil(t, seq.Slice(seq2.Values(seq2.Series(-1, (func(i, prev int) (int, bool))(nil)))))
}

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

func pairSum(prev *string, i int, val string) string {
	r := strconv.Itoa(i) + val
	if prev == nil {
		return r
	}
	return *prev + r
}

func Test_ReduceSum(t *testing.T) {

	sum, ok := seq2.ReduceOK(seq2.Of("A", "B", "C"), pairSum)

	assert.True(t, ok)
	assert.Equal(t, "0A1B2C", sum)
}

func Test_ReduceeSum(t *testing.T) {
	s := seq2.Of(1, 3, 5, 7, 9, 11)
	reducer := func(prev *int, _, v int) (int, error) {
		p := 0
		if prev != nil {
			p = *prev
		}
		if v == 11 {
			return p, errors.New("stop")
		}
		return v + p, nil
	}
	r, ok, err := seq2.ReduceeOK(s, reducer)
	assert.True(t, ok)
	assert.Equal(t, 1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")

	_, ok, err = seq2.ReduceeOK[seq.Seq2[int, int]](nil, reducer)
	assert.False(t, ok)
	assert.NoError(t, err)

	r, err = seq2.Reducee(s, reducer)
	assert.Equal(t, 1+3+5+7+9, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceeSumFirstErr(t *testing.T) {
	s := seq2.Of(1, 3, 5, 7, 9, 11)
	r, ok, err := seq2.ReduceeOK(s, func(_ *int, _, _ int) (int, error) {
		return 0, errors.New("stop")
	})
	assert.True(t, ok)
	assert.Equal(t, 0, r)
	assert.ErrorContains(t, err, "stop")
}

func Test_ReduceEmpty(t *testing.T) {
	s := seq2.Of[string]()
	sum, ok := seq2.ReduceOK(s, pairSum)

	assert.False(t, ok)
	assert.Equal(t, "", sum)
}

func Test_Head(t *testing.T) {
	sequence := seq2.Of(1, 2, 3, 4, 5, 6)
	_, result, ok := seq2.Head(sequence)

	assert.True(t, ok)
	assert.Equal(t, 1, result)

	_, result, ok = seq2.Head[seq.Seq2[int, int]](nil)
	assert.Zero(t, result)
	assert.False(t, ok)
}

func Test_While(t *testing.T) {
	sequence := seq2.Of(1, 2, 3, 4, 5, 6)
	part := seq2.While(sequence, kvpredicate.Value[int](not.Eq(5)))

	assert.Equal(t, slice.Of(1, 2, 3, 4), seq.Slice(seq2.Values(part)))

	part = seq2.While(sequence, kvpredicate.Value[int](not.Eq(7)))
	assert.Equal(t, slice.Of(1, 2, 3, 4, 5, 6), seq.Slice(seq2.Values(part)))

	part = seq2.While(sequence, kvpredicate.Value[int](eq.To(0)))
	assert.Nil(t, seq.Slice(seq2.Values(part)))

	part = seq2.While[seq.Seq2[int, int]](nil, kvpredicate.Value[int](eq.To(0)))
	assert.Nil(t, seq.Slice(seq2.Values(part)))

	r := []int{}
	for _, i := range seq2.While(sequence, kvpredicate.Value[int](not.Eq(7))) {
		if i == 5 {
			break
		}
		r = append(r, i)
	}
	assert.Equal(t, slice.Of(1, 2, 3, 4), r)
}

func Test_SkipWhile(t *testing.T) {
	sequence := seq2.Of(1, 2, 3, 4, 5, 6)
	part := seq2.SkipWhile(sequence, kvpredicate.Value[int](less.Than(4)))

	assert.Equal(t, slice.Of(4, 5, 6), seq.Slice(seq2.Values(part)))

	part = seq2.SkipWhile(sequence, kvpredicate.Value[int](not.Eq(7)))
	assert.Nil(t, seq.Slice(seq2.Values(part)))

	part = seq2.SkipWhile(sequence, kvpredicate.Value[int](less.Than(0)))
	assert.Equal(t, slice.Of(1, 2, 3, 4, 5, 6), seq.Slice(seq2.Values(part)))

	r := []int{}
	for _, i := range seq2.SkipWhile(sequence, kvpredicate.Value[int](less.Than(4))) {
		if i == 6 {
			break
		}
		r = append(r, i)
	}
	assert.Equal(t, slice.Of(4, 5), r)
}

func Test_Top(t *testing.T) {
	sequence := seq2.Of(1, 2, 3, 4, 5, 6)
	top := seq2.Values(seq2.Top(4, sequence))
	result := seq.Slice(top)
	result2 := seq.Slice(top)

	assert.Equal(t, slice.Of(1, 2, 3, 4), result)
	assert.Equal(t, result2, result)

	result = seq.Slice(seq2.Values(seq2.Top(0, sequence)))
	assert.Nil(t, result)

	result = seq.Slice(seq2.Values(seq2.Top[seq.Seq2[int, int]](10, nil)))
	assert.Nil(t, result)
	result = nil
	for _, v := range seq2.Top(4, seq2.Of(1, 2, 3, 4, 5, 6)) {
		if v != 3 {
			result = append(result, v)
		} else {
			break
		}
	}
	assert.Equal(t, slice.Of(1, 2), result)
}

func Test_Skip(t *testing.T) {
	sequence := seq2.Of(1, 2, 3, 4, 5, 6)
	skip := seq2.Values(seq2.Skip(4, sequence))
	result := seq.Slice(skip)
	result2 := seq.Slice(skip)

	assert.Equal(t, slice.Of(5, 6), result)
	assert.Equal(t, result2, result)

	result = seq.Slice(seq2.Values(seq2.Skip(0, sequence)))
	assert.Equal(t, seq.Slice(seq2.Values(sequence)), result)

	result = seq.Slice(seq2.Values(seq2.Skip[seq.Seq2[int, int]](10, nil)))
	assert.Nil(t, result)
	result = nil
	for _, v := range seq2.Skip(2, seq2.Of(1, 2, 3, 4, 5, 6)) {
		if v != 5 {
			result = append(result, v)
		} else {
			break
		}
	}
	assert.Equal(t, slice.Of(3, 4), result)
}

func Test_SkipTop(t *testing.T) {
	sequence := seq2.Of(1, 2, 3, 4, 5, 6)
	middle := seq2.Top(2, seq2.Skip(2, sequence))
	result := seq.Slice(seq2.Values(middle))
	i := seq.Slice(seq2.Keys(middle))

	assert.Equal(t, slice.Of(3, 4), result)
	assert.Equal(t, slice.Of(2, 3), i)
}

func Test_First(t *testing.T) {
	sequence := seq2.Of(1, 2, 3, 4, 5, 6)
	condition := func(_ int, v int) bool { return more.Than(5)(v) }
	i, result, ok := seq2.First(sequence, condition)

	assert.True(t, ok)
	assert.Equal(t, 5, i)
	assert.Equal(t, 6, result)
	_, _, ok = seq2.First[seq.Seq2[int, int]](nil, condition)
	assert.False(t, ok)

	_, _, ok = seq2.First(sequence, nil)
	assert.False(t, ok)
}

func Test_Filter(t *testing.T) {
	s := seq2.Filter(seq2.Of("first", "second", "third"), func(i int, _ string) bool { return i%2 == 0 })
	k := seq.Slice(seq2.Keys(s))
	v := seq.Slice(seq2.Values(s))

	assert.Equal(t, slice.Of(0, 2), k)
	assert.Equal(t, slice.Of("first", "third"), v)
}

var testMap = map_.Of(k.V(1, "10"), k.V(2, "20"), k.V(3, "30"), k.V(5, "50"), k.V(7, "70"), k.V(8, "80"), k.V(9, "90"), k.V(11, "110"))

func Test_SeqOfNil(t *testing.T) {
	var in, out []int

	iter := false
	for _, e := range seq2.Of(in...) {
		iter = true
		out = append(out, e)
	}

	assert.Nil(t, out)
	assert.False(t, iter)
}

func Test_OfMap(t *testing.T) {
	in := map[int]string{}

	iter := false
	for _ = range seq2.OfMap(in) {
		iter = true
	}
	assert.False(t, iter)

	in[0] = "1"
	for _ = range seq2.OfMap(in) {
		iter = true
	}
	assert.True(t, iter)

	ignoreBreak := false
	for key := range seq2.OfMap(in) {
		if key == 0 {
			break
		}
		ignoreBreak = true
	}
	assert.False(t, ignoreBreak)
}

func Test_ConvertNilSeq(t *testing.T) {
	var in iter.Seq2[int, int]
	var out []int

	iter := false
	for _, e := range seq2.Convert(in, func(i, e int) (int, int) { return i, e }) {
		iter = true
		out = append(out, e)
	}

	assert.Nil(t, out)
	assert.False(t, iter)
}

func Test_AllFiltered(t *testing.T) {
	s := []string{}

	for _, v := range seq2.Filter(testMap.All, func(k int, _ string) bool { return k%2 == 0 }) {
		s = append(s, v)
	}

	assert.Equal(t, slice.Of("20", "80"), sort.Asc(s))
}

func Test_AllConverted(t *testing.T) {
	i := []int{}

	for _, e := range seq2.Convert(testMap.All, func(k int, v string) (int, int) { c, _ := strconv.Atoi(v); return k, c }) {
		i = append(i, e)
	}

	assert.Equal(t, slice.Of(10, 20, 30, 50, 70, 80, 90, 110), i)
}

func Test_AllConv(t *testing.T) {
	i := []int{}

	for kv, err := range seq2.Conv(testMap.All, func(k int, v string) (int, int, error) {
		if k == 5 {
			return 0, 0, errors.New("stop")
		}
		c, err := strconv.Atoi(v)
		return k, c, err
	}) {
		if err != nil {
			break
		}
		i = append(i, kv.V)
	}

	assert.Equal(t, slice.Of(10, 20, 30), i)
}

func Test_Slice_ToMapResolvOrder(t *testing.T) {
	var (
		even          = func(v int) bool { return v%2 == 0 }
		order, groups = seq2.MapResolvOrder(seq.ToSeq2(seq.Of(2, 1, 1, 2, 4, 3, 1), func(i int) (bool, int) {
			return even(i), i
		}), func(_ bool, _ bool, valResov []int, val int) []int {
			return append(valResov, val)
		})
	)
	assert.Equal(t, []int{1, 1, 3, 1}, groups[false])
	assert.Equal(t, []int{2, 2, 4}, groups[true])
	assert.Equal(t, []bool{true, false}, order)
}

func Test_Range(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), seq.Slice(seq2.Values(seq2.Range(-1, 4))))
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), seq.Slice(seq2.Keys(seq2.Range(-1, 4))))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), seq.Slice(seq2.Values(seq2.Range(3, -2))))
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), seq.Slice(seq2.Keys(seq2.Range(3, -2))))
	assert.Nil(t, seq.Slice(seq2.Values(seq2.Range(1, 1))))

	var out, ind []int
	for i, v := range seq2.Range(-1, 3) {
		if v == 2 {
			break
		}
		out = append(out, v)
		ind = append(ind, i)
	}
	assert.Equal(t, slice.Of(-1, 0, 1), out)
	assert.Equal(t, slice.Of(0, 1, 2), ind)
}

func Test_RangeClosed(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), seq.Slice(seq2.Values(seq2.RangeClosed(-1, 3))))
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), seq.Slice(seq2.Keys(seq2.RangeClosed(-1, 3))))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), seq.Slice(seq2.Values(seq2.RangeClosed(3, -1))))
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), seq.Slice(seq2.Keys(seq2.RangeClosed(3, -1))))
	assert.Equal(t, slice.Of(1), seq.Slice(seq2.Values(seq2.RangeClosed(1, 1))))

	var out, ind []int
	for i, v := range seq2.RangeClosed(-1, 3) {
		if v == 2 {
			break
		}
		out = append(out, v)
		ind = append(ind, i)
	}
	assert.Equal(t, slice.Of(-1, 0, 1), out)
	assert.Equal(t, slice.Of(0, 1, 2), ind)
}

func Test_ToSeq(t *testing.T) {
	s := seq.Slice(seq2.ToSeq(seq2.Of("A", "B", "C"), func(i int, v string) string { return strconv.Itoa(i) + v }))
	assert.Equal(t, slice.Of("0A", "1B", "2C"), s)
}

func Test_TrackEach(t *testing.T) {
	var out, ind []int
	seq2.TrackEach(seq2.RangeClosed(-1, 3), func(i int, v int) {
		out = append(out, v)
		ind = append(ind, i)
	})
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), out)
	assert.Equal(t, slice.Of(0, 1, 2, 3, 4), ind)
}
