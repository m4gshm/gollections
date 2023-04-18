package test

import (
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/mutable/map_"
	"github.com/m4gshm/gollections/mutable/omap"
	"github.com/m4gshm/gollections/mutable/ordered"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	ordered := omap.Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Map()))

	expectedK := slice.Of(1, 2, 4, 3)
	expectedV := slice.Of("1", "2", "4", "3")

	keys := make([]int, 0)
	values := make([]string, 0)
	it := ordered.Begin()
	for key, val, ok := it.Next(); ok; key, val, ok = it.Next() {
		keys = append(keys, key)
		values = append(values, val)
	}
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	assert.Equal(t, slice.Of(1, 2, 4, 3), ordered.Keys().Slice())
	assert.Equal(t, slice.Of("1", "2", "4", "3"), ordered.Values().Slice())
}

func Test_Map_Add(t *testing.T) {
	d := omap.New[int, string](4)
	s := d.SetNew(1, "1")
	assert.Equal(t, s, true)
	d.Set(2, "2")
	d.Set(4, "4")
	d.Set(3, "3")
	s = d.SetNew(1, "11")
	assert.Equal(t, s, false)

	assert.Equal(t, slice.Of(1, 2, 4, 3), d.Keys().Slice())
	assert.Equal(t, slice.Of("1", "2", "4", "3"), d.Values().Slice())
}

func Test_Map_Nil(t *testing.T) {
	var m *ordered.Map[string, string]

	m.Set("a", "A")
	assert.False(t, m.SetNew("b", "B"))

	assert.False(t, m.Contains("b"))

	m.SetMap(nil)
	m.SetMap(map_.Of(K.V("d", "D")))

	out := m.Map()
	assert.Equal(t, 0, len(out))

	e := m.IsEmpty()
	assert.True(t, e)

	head, _, _, ok := m.First()
	assert.False(t, ok)

	head = m.Head()
	_, _, ok = head.Next()
	assert.False(t, ok)

	m.For(nil)
	m.ForEach(nil)
	m.Track(nil)
	m.TrackEach(nil)

	m.Reduce(nil)
	m.Convert(nil).Track(nil)
	m.ConvertKey(nil).Next()
	m.ConvertKey(nil).Track(nil)
	m.ConvertKey(nil).TrackEach(nil)
	m.ConvertValue(nil).Next()
	m.ConvertValue(nil).Track(nil)
	m.ConvertValue(nil).TrackEach(nil)
	m.Filter(nil).Convert(nil).Track(nil)
	m.Filter(nil).Convert(nil).TrackEach(nil)

	m.Keys().For(nil)
	m.Keys().ForEach(nil)
	m.Values().For(nil)
	m.Values().ForEach(nil)
	m.Values().Convert(nil).For(nil)
	m.Values().Filter(nil).ForEach(nil)

}

func Test_Map_Zero(t *testing.T) {
	var m ordered.Map[string, string]

	m.Set("a", "A")
	assert.True(t, m.SetNew("b", "B"))
	assert.False(t, m.SetNew("b", "B"))

	assert.True(t, m.Contains("b"))

	m.SetMap(nil)
	m.SetMap(map_.Of(K.V("d", "D")))

	out := m.Map()
	assert.Equal(t, 3, len(out))

	e := m.IsEmpty()
	assert.False(t, e)

	head, k, v, ok := m.First()
	assert.True(t, ok)
	assert.Equal(t, "a", k)
	assert.Equal(t, "A", v)

	head = m.Head()
	_, _, ok = head.Next()
	assert.True(t, ok)

	m.For(nil)
	m.ForEach(nil)
	m.Track(nil)
	m.TrackEach(nil)

	m.Reduce(nil)
	m.Convert(nil).Track(nil)
	m.ConvertKey(nil).Next()
	m.ConvertKey(nil).Track(nil)
	m.ConvertKey(nil).TrackEach(nil)
	m.ConvertValue(nil).Next()
	m.ConvertValue(nil).Track(nil)
	m.ConvertValue(nil).TrackEach(nil)
	m.Filter(nil).Convert(nil).Track(nil)
	m.Filter(nil).Convert(nil).TrackEach(nil)

	m.Keys().For(nil)
	m.Keys().ForEach(nil)
	m.Values().For(nil)
	m.Values().ForEach(nil)
	m.Values().Convert(nil).For(nil)
	m.Values().Filter(nil).ForEach(nil)
}

