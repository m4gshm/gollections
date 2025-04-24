package test

import (
	"iter"
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/ordered/map_"
	"github.com/m4gshm/gollections/k"
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
	generator := func(prev int) (int, bool) { return prev + 1, prev < 3 }
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

	assert.Nil(t, seq.Slice(seq2.Values(seq2.Series(-1, (func(prev int) (int, bool))(nil)))))
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
	var in iter.Seq2[int, int] = nil
	var out []int = nil

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

func Test_Slice_ToMapResolvOrder(t *testing.T) {
	var (
		even          = func(v int) bool { return v%2 == 0 }
		order, groups = seq2.MapResolvOrder(seq.ToSeq2(seq.Of(2, 1, 1, 2, 4, 3, 1), func(i int) (bool, int) {
			return even(i), i
		}), func(exists bool, key bool, valResov []int, val int) []int {
			return append(valResov, val)
		})
	)
	assert.Equal(t, []int{1, 1, 3, 1}, groups[false])
	assert.Equal(t, []int{2, 2, 4}, groups[true])
	assert.Equal(t, []bool{true, false}, order)
}
