package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/mutable/map_"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	omap "github.com/m4gshm/gollections/collection/mutable/ordered/map_"
	"github.com/m4gshm/gollections/k"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func Test_Map_Iterate(t *testing.T) {
	ordered := omap.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))
	assert.Equal(t, 4, len(ordered.Map()))

	expectedK := slice.Of(1, 2, 4, 3)
	expectedV := slice.Of("1", "2", "4", "3")

	keys := make([]int, 0)
	values := make([]string, 0)
	next := ordered.Loop()
	for key, val, ok := next(); ok; key, val, ok = next() {
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
	m.SetMap(map_.Of(k.V("d", "D")))

	out := m.Map()
	assert.Equal(t, 0, len(out))

	e := m.IsEmpty()
	assert.True(t, e)

	head, _, _, ok := m.First()
	assert.False(t, ok)

	head = m.Head()
	_, _, ok = head.Next()
	assert.False(t, ok)

	m.Track(nil)
	m.TrackEach(nil)

	m.Reduce(nil)
	m.Convert(nil).Track(nil)
	m.ConvertKey(nil).FiltKey(nil)
	m.ConvertKey(nil).Track(nil)
	m.ConvertValue(nil).Track(nil)
	m.Filter(nil).Convert(nil).Track(nil)

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
	m.SetMap(map_.Of(k.V("d", "D")))

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

	m.Track(func(k, v string) error { return nil })
	m.TrackEach(func(k, v string) {})

	m.Reduce(func(k1, v1, k2, v2 string) (string, string) { return k1 + k2, v1 + v2 })
	m.Convert(func(_, _ string) (string, string) { return k, v }).Track(func(_, _ string) error { return nil })
	m.ConvertKey(func(s string) string { return s }).Track(func(_, _ string) error { return nil })
	m.ConvertValue(func(s string) string { return s }).Track(func(_, _ string) error { return nil })
	m.Filter(func(_, _ string) bool { return true }).Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(_, _ string) error { return nil })

	m.Keys().For(func(_ string) error { return nil })
	m.Keys().ForEach(func(_ string) {})
	m.Keys().Convert(func(s string) string { return s }).Slice()
	m.Keys().Convert(func(s string) string { return s }).For(func(_ string) error { return nil })
	m.Keys().Filter(func(_ string) bool { return true }).Slice()
	m.Keys().Filter(func(_ string) bool { return true }).ForEach(func(_ string) {})

	m.Values().For(func(_ string) error { return nil })
	m.Values().ForEach(func(_ string) {})
	m.Values().Convert(func(s string) string { return s }).Slice()
	m.Values().Convert(func(s string) string { return s }).For(func(_ string) error { return nil })
	m.Values().Filter(func(_ string) bool { return true }).Slice()
	m.Values().Filter(func(_ string) bool { return true }).ForEach(func(_ string) {})
}

func Test_Map_new(t *testing.T) {
	var m = new(ordered.Map[string, string])

	m.Set("a", "A")
	assert.True(t, m.SetNew("b", "B"))
	assert.False(t, m.SetNew("b", "B"))

	assert.True(t, m.Contains("b"))

	m.SetMap(nil)
	m.SetMap(map_.Of(k.V("d", "D")))

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

	m.Track(func(k, v string) error { return nil })
	m.TrackEach(func(k, v string) {})

	m.Reduce(func(k1, v1, k2, v2 string) (string, string) { return k1 + k2, v1 + v2 })
	m.Convert(func(_, _ string) (string, string) { return k, v }).Track(func(_, _ string) error { return nil })
	m.ConvertKey(func(s string) string { return s }).Track(func(_, _ string) error { return nil })
	m.ConvertValue(func(s string) string { return s }).Track(func(_, _ string) error { return nil })
	m.Filter(func(_, _ string) bool { return true }).Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(_, _ string) error { return nil })

	m.Keys().For(func(_ string) error { return nil })
	m.Keys().ForEach(func(_ string) {})
	m.Keys().Convert(func(s string) string { return s }).Slice()
	m.Keys().Convert(func(s string) string { return s }).For(func(_ string) error { return nil })
	m.Keys().Filter(func(_ string) bool { return true }).Slice()
	m.Keys().Filter(func(_ string) bool { return true }).ForEach(func(_ string) {})

	m.Values().For(func(_ string) error { return nil })
	m.Values().ForEach(func(_ string) {})
	m.Values().Convert(func(s string) string { return s }).Slice()
	m.Values().Convert(func(s string) string { return s }).For(func(_ string) error { return nil })
	m.Values().Filter(func(_ string) bool { return true }).Slice()
	m.Values().Filter(func(_ string) bool { return true }).ForEach(func(_ string) {})
}

func Test_Map_Sort(t *testing.T) {
	var m = new(ordered.Map[int, string])

	m.Set(5, "5")
	m.Set(4, "4")
	m.Set(-8, "-8")
	m.Set(10, "10")

	o := m.Sort(op.Compare)

	expected := ordered.NewMap(k.V(-8, "-8"), k.V(4, "4"), k.V(5, "5"), k.V(10, "10"))

	assert.Equal(t, expected, o)
	assert.Same(t, m, o)
}
