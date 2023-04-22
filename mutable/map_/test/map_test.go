package test

import (
	"sort"
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/mutable/map_"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Iterate(t *testing.T) {
	unordered := map_.Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(unordered.Map()))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)
	for it, key, val, ok := unordered.First(); ok; key, val, ok = it.Next() {
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = unordered.Keys().Slice()
	sort.Ints(keys)
	values = unordered.Values().Slice()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_IterateOverRange(t *testing.T) {
	unordered := map_.Of(K.V(1, "1"), K.V(1, "1"), K.V(2, "2"), K.V(4, "4"), K.V(3, "3"), K.V(1, "1"))
	assert.Equal(t, 4, len(unordered.Map()))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)
	for key, val := range *unordered {
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = unordered.Keys().Slice()
	sort.Ints(keys)
	values = unordered.Values().Slice()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_Add(t *testing.T) {
	d := map_.New[int, string](4)
	s := d.SetNew(1, "1")
	assert.Equal(t, s, true)

	d.Set(2, "2")
	d.Set(4, "4")
	d.Set(3, "3")

	s = d.SetNew(1, "11")
	assert.Equal(t, s, false)

	keys := d.Keys().Slice()
	sort.Ints(keys)
	values := d.Values().Slice()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_Nil(t *testing.T) {
	var m *mutable.Map[string, string]

	m.Set("a", "A")
	assert.False(t, m.SetNew("b", "B"))

	m.SetMap(nil)
	m.SetMap(map_.Of(K.V("d", "D")))

	out := m.Map()
	assert.Equal(t, 0, len(out))

	m.Delete("b")
	var nilKeys []string
	m.Delete(nilKeys...)
	m.DeleteOne("a")

	e := m.IsEmpty()
	assert.True(t, e)

	head, _, _, ok := m.First()
	assert.False(t, ok)

	head = m.Head()
	_, _, ok = head.Next()
	assert.False(t, ok)

	m.Reduce(nil)
	m.Convert(nil).Track(nil)
	m.ConvertKey(nil).Next()
	m.ConvertValue(nil).Next()
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
	var m mutable.Map[string, string]

	m.Set("a", "A")
	assert.True(t, m.SetNew("b", "B"))
	assert.False(t, m.SetNew("b", "B"))

	m.SetMap(nil)
	m.SetMap(map_.Of(K.V("d", "D")))

	out := m.Map()
	assert.Equal(t, 3, len(out))

	m.Delete("b")
	var nilKeys []string
	m.Delete(nilKeys...)
	m.DeleteOne("a")

	e := m.IsEmpty()
	assert.False(t, e)

	l := m.Len()
	assert.Equal(t, 1, l)

	head, k, v, ok := m.First()
	assert.True(t, ok)
	assert.Equal(t, "d", k)
	assert.Equal(t, "D", v)

	head = m.Head()
	_, _, ok = head.Next()
	assert.True(t, ok)

	m.Reduce(func(k1, v1, k2, v2 string) (string, string) { return k1 + k2, v1 + v2 })
	m.Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(position, element string) error { return nil })
	m.ConvertKey(as.Is[string]).Next()
	m.ConvertValue(as.Is[string]).Next()
	m.Filter(func(s1, s2 string) bool { return true }).Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(position, element string) error { return nil })
	m.Filter(func(s1, s2 string) bool { return true }).Convert(func(s1, s2 string) (string, string) { return s1, s2 }).TrackEach(func(position, element string) {})

	m.Keys().For(func(element string) error { return nil })
	m.Keys().ForEach(func(element string) {})
	m.Values().For(func(element string) error { return nil })
	m.Values().ForEach(func(element string) {})
	m.Values().Convert(as.Is[string]).For(func(element string) error { return nil })
	m.Values().Filter(func(s string) bool { return true }).ForEach(func(element string) {})
}

func Test_Map_new(t *testing.T) {
	var m = new(mutable.Map[string, string])

	m.Set("a", "A")
	assert.True(t, m.SetNew("b", "B"))
	assert.False(t, m.SetNew("b", "B"))

	m.SetMap(nil)
	m.SetMap(map_.Of(K.V("d", "D")))

	out := m.Map()
	assert.Equal(t, 3, len(out))

	m.Delete("b")
	var nilKeys []string
	m.Delete(nilKeys...)
	m.DeleteOne("a")

	e := m.IsEmpty()
	assert.False(t, e)

	l := m.Len()
	assert.Equal(t, 1, l)

	head, k, v, ok := m.First()
	assert.True(t, ok)
	assert.Equal(t, "d", k)
	assert.Equal(t, "D", v)

	head = m.Head()
	_, _, ok = head.Next()
	assert.True(t, ok)

	m.Reduce(func(k1, v1, k2, v2 string) (string, string) { return k1 + k2, v1 + v2 })
	m.Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(position, element string) error { return nil })
	m.ConvertKey(as.Is[string]).Next()
	m.ConvertValue(as.Is[string]).Next()
	m.Filter(func(s1, s2 string) bool { return true }).Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(position, element string) error { return nil })
	m.Filter(func(s1, s2 string) bool { return true }).Convert(func(s1, s2 string) (string, string) { return s1, s2 }).TrackEach(func(position, element string) {})

	m.Keys().For(func(element string) error { return nil })
	m.Keys().ForEach(func(element string) {})
	m.Values().For(func(element string) error { return nil })
	m.Values().ForEach(func(element string) {})
	m.Values().Convert(as.Is[string]).For(func(element string) error { return nil })
	m.Values().Filter(func(s string) bool { return true }).ForEach(func(element string) {})
}
