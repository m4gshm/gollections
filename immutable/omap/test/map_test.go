package test

import (
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/immutable/omap"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	ordered := omap.Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Map()))

	expectedK := slice.Of(1, 2, 4, 3)
	expectedV := slice.Of("1", "2", "4", "3")

	keys := []int{}
	values := []string{}
	for it, key, val, ok := ordered.First(); ok; key, val, ok = it.Next() {
		keys = append(keys, key)
		values = append(values, val)
	}
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	assert.Equal(t, slice.Of(1, 2, 4, 3), ordered.Keys().Slice())
	assert.Equal(t, slice.Of("1", "2", "4", "3"), ordered.Values().Slice())
}

func Test_Map_Iterate_Keys(t *testing.T) {
	ordered := omap.Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Map()))

	expectedK := slice.Of(1, 2, 4, 3)

	keys := []int{}
	for it, key, ok := ordered.K().First(); ok; key, ok = it.Next() {
		keys = append(keys, key)
	}
	assert.Equal(t, expectedK, keys)

}

func Test_Map_Iterate_Values(t *testing.T) {
	ordered := omap.Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Map()))

	expectedV := slice.Of("1", "2", "4", "3")

	values := []string{}
	for it, val, ok := ordered.V().First(); ok; val, ok = it.Next() {
		values = append(values, val)
	}

	assert.Equal(t, expectedV, values)
}

func Test_Map_Nil(t *testing.T) {
	var m *ordered.Map[string, string]

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

	_ = m.For(nil)
	m.ForEach(nil)
	_ = m.Track(nil)
	m.TrackEach(nil)

	m.Filter(nil)
	m.FilterKey(nil)
	m.FilterValue(nil)

	_ = m.Values().For(nil)
	m.Values().ForEach(nil)
	_ = m.ConvertValue(nil).Track(nil)
	m.ConvertValue(nil).TrackEach(nil)
	m.ConvertValue(nil).Filter(nil).FilterKey(nil)
	m.ConvertValue(nil).Filter(nil).FilterValue(nil)

	_ = m.Keys().For(nil)
	m.Keys().ForEach(nil)
	_ = m.ConvertKey(nil).Track(nil)
	m.ConvertKey(nil).TrackEach(nil)
	m.ConvertKey(nil).Filter(nil).FilterKey(nil)
	m.ConvertKey(nil).Filter(nil).FilterValue(nil)
	m.Convert(nil)

	_ = m.Sort(nil).For(nil)
	m.Sort(nil).ForEach(nil)

	_ = m.StableSort(nil).For(nil)
	m.StableSort(nil).ForEach(nil)
}

func Test_Map_Zero(t *testing.T) {
	var m ordered.Map[string, string]

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

	_ = m.For(nil)
	m.ForEach(nil)
	m.Track(nil)
	m.TrackEach(nil)

	m.Filter(nil)
	m.FilterKey(nil)
	m.FilterValue(nil)

	m.Values().For(nil)
	m.Values().ForEach(nil)
	m.ConvertValue(nil).Track(nil)
	m.ConvertValue(nil).TrackEach(nil)
	m.ConvertValue(nil).Filter(nil).FilterKey(nil)
	m.ConvertValue(nil).Filter(nil).FilterValue(nil)

	m.Keys().For(nil)
	m.Keys().ForEach(nil)
	m.ConvertKey(nil).Track(nil)
	m.ConvertKey(nil).TrackEach(nil)
	m.ConvertKey(nil).Filter(nil).FilterKey(nil)
	m.ConvertKey(nil).Filter(nil).FilterValue(nil)
	m.Convert(nil)

	m.Sort(nil).For(nil)
	m.Sort(nil).ForEach(nil)

	m.StableSort(nil).For(nil)
	m.StableSort(nil).ForEach(nil)
}
