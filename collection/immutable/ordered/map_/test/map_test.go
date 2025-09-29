package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/collection/immutable/ordered/map_"
	"github.com/m4gshm/gollections/k"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
)

func Test_Map_Of(t *testing.T) {
	m := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))
	iterCheck(t, m)
}

func Test_Map_FromSeq(t *testing.T) {
	m := map_.FromSeq2(seq.ToKV(seq.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1")), c.KV[int, string].Key, c.KV[int, string].Value))
	iterCheck(t, m)
}

func iterCheck(t *testing.T, m ordered.Map[int, string]) {
	assert.Equal(t, 4, m.Len())
	assert.Equal(t, 4, len(m.Map()))

	expectedK := slice.Of(1, 2, 4, 3)
	expectedV := slice.Of("1", "2", "4", "3")

	keys := []int{}
	values := []string{}
	for key, val := range m.All {
		keys = append(keys, key)
		values = append(values, val)
	}
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = m.Keys().Slice()
	values = m.Values().Slice()
	assert.Equal(t, slice.Of(1, 2, 4, 3), keys)
	assert.Equal(t, slice.Of("1", "2", "4", "3"), values)
}

func Test_Map_Iterate_Keys(t *testing.T) {
	ordered := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Map()))

	expectedK := slice.Of(1, 2, 4, 3)

	keys := []int{}
	for key := range ordered.Keys().All {
		keys = append(keys, key)
	}
	assert.Equal(t, expectedK, keys)

}

func Test_Map_Iterate_Values(t *testing.T) {
	ordered := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Map()))

	expectedV := slice.Of("1", "2", "4", "3")

	values := []string{}
	for val := range ordered.Values().All {
		values = append(values, val)
	}

	assert.Equal(t, expectedV, values)
}

func Test_Map_Zero(t *testing.T) {
	var m ordered.Map[string, string]

	m.Contains("")

	out := m.Map()
	assert.Equal(t, 0, len(out))

	e := m.IsEmpty()
	assert.True(t, e)

	_, _, ok := m.Head()
	assert.False(t, ok)

	_, ok = m.Get("")
	assert.False(t, ok)

	m.TrackEach(nil)

	m.Filter(nil)
	m.FilterKey(nil)
	m.FilterValue(nil)

	m.Values().ForEach(nil)
	m.ConvertValue(nil).TrackEach(nil)
	m.ConvertValue(nil).Filter(nil).FilterKey(nil)
	m.ConvertValue(nil).Filter(nil).FilterValue(nil)

	m.Keys().ForEach(nil)
	m.ConvertKey(nil).TrackEach(nil)
	m.ConvertKey(nil).Filter(nil).FilterKey(nil)
	m.ConvertKey(nil).Filter(nil).FilterValue(nil)
	m.Convert(nil)

	m.Sort(nil).TrackEach(nil)

	m.StableSort(nil).TrackEach(nil)
}

func Test_Map_Sort(t *testing.T) {
	var m = map_.Of(k.V(5, "5"), k.V(4, "4"), k.V(-8, "-8"), k.V(10, "10"))
	o := m.Sort(op.Compare)

	expected := ordered.NewMap(k.V(-8, "-8"), k.V(4, "4"), k.V(5, "5"), k.V(10, "10"))

	assert.Equal(t, expected, o)
	assert.NotSame(t, &m, &o)
}
