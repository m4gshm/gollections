package test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/collection/immutable/map_"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/k"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func Test_Map_Iterate(t *testing.T) {
	dict := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))

	assert.Equal(t, 4, dict.Len())
	assert.Equal(t, 4, len(dict.Map()))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)

	it := dict.Head()
	for key, val, ok := it.Next(); ok; key, val, ok = it.Next() {
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = dict.Keys().Slice()
	values = dict.Values().Slice()
	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_Iterate_Keys(t *testing.T) {
	dict := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))
	assert.Equal(t, 4, len(dict.Map()))

	expectedK := slice.Of(1, 2, 3, 4)

	keys := []int{}
	for it, key, ok := dict.K().First(); ok; key, ok = it.Next() {
		keys = append(keys, key)
	}

	sort.Ints(keys)
	assert.Equal(t, expectedK, keys)

}

func Test_Map_Iterate_Values(t *testing.T) {
	ordered := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Map()))

	expectedV := slice.Of("1", "2", "3", "4")

	values := []string{}
	for it, val, ok := ordered.V().First(); ok; val, ok = it.Next() {
		values = append(values, val)
	}

	sort.Strings(values)
	assert.Equal(t, expectedV, values)
}

func Test_Map_Zero(t *testing.T) {
	var m immutable.Map[string, string]

	m.Contains("")

	out := m.Map()
	assert.Equal(t, 0, len(out))

	e := m.IsEmpty()
	assert.True(t, e)

	head, _, _, ok := m.First()
	assert.False(t, ok)
	_, _, ok = head.Next()
	assert.False(t, ok)

	head = m.Head()
	_, _, ok = head.Next()
	assert.False(t, ok)

	_, ok = m.Get("")
	assert.False(t, ok)

	m.For(nil)
	m.ForEach(nil)
	m.Track(nil)
	m.TrackEach(nil)

	m.Filter(nil)
	m.FilterKey(nil)
	m.FilterValue(nil)

	m.Values().For(nil)
	m.Values().ForEach(nil)
	m.ConvertValue(nil).Track(nil)
	m.ConvertValue(nil).Filter(nil).FilterKey(nil)
	m.ConvertValue(nil).Filter(nil).FilterValue(nil)

	m.Keys().For(nil)
	m.Keys().ForEach(nil)
	m.ConvertKey(nil).Track(nil)
	m.ConvertKey(nil).Filter(nil).FilterKey(nil)
	m.ConvertKey(nil).Filter(nil).FilterValue(nil)
	m.Convert(nil)

	m.Sort(nil).For(nil)
	m.Sort(nil).ForEach(nil)

	m.StableSort(nil).For(nil)
	m.StableSort(nil).ForEach(nil)
}

func Test_Map_Sort(t *testing.T) {
	var m = map_.Of(k.V(5, "5"), k.V(4, "4"), k.V(-8, "-8"), k.V(10, "10"))
	o := m.Sort(op.Compare)

	expected := ordered.NewMap(k.V(-8, "-8"), k.V(4, "4"), k.V(5, "5"), k.V(10, "10"))

	assert.Equal(t, expected, o)
	assert.NotSame(t, m, o)
}
