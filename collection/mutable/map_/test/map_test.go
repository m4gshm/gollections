package test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/map_"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/k"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
)

func Test_Map_Of(t *testing.T) {
	m := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))
	iterCheck(t, m)
}

func Test_Map_From(t *testing.T) {
	m := map_.From(loop.KeyValue(loop.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1")), c.KV[int, string].Key, c.KV[int, string].Value))
	iterCheck(t, m)
}

func Test_Map_FromSeq(t *testing.T) {
	m := map_.FromSeq2(seq.KeyValue(seq.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1")), c.KV[int, string].Key, c.KV[int, string].Value))
	iterCheck(t, m)
}

func iterCheck(t *testing.T, m *mutable.Map[int, string]) {
	assert.Equal(t, 4, m.Len())
	assert.Equal(t, 4, len(m.Map()))

	expectedK := slice.Of(1, 2, 3, 4)
	expectedV := slice.Of("1", "2", "3", "4")

	keys := make([]int, 0)
	values := make([]string, 0)
	for key, val := range m.All {
		keys = append(keys, key)
		values = append(values, val)
	}

	sort.Ints(keys)
	sort.Strings(values)
	assert.Equal(t, expectedK, keys)
	assert.Equal(t, expectedV, values)

	keys = m.Keys().Slice()
	sort.Ints(keys)
	values = m.Values().Slice()
	sort.Strings(values)
	assert.Equal(t, slice.Of(1, 2, 3, 4), keys)
	assert.Equal(t, slice.Of("1", "2", "3", "4"), values)
}

func Test_Map_IterateOverRange(t *testing.T) {
	unordered := map_.Of(k.V(1, "1"), k.V(1, "1"), k.V(2, "2"), k.V(4, "4"), k.V(3, "3"), k.V(1, "1"))
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
	m.SetMap(map_.Of(k.V("d", "D")))

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
	m.ConvertKey(nil).Filter(nil)
	m.ConvertValue(nil).Filter(nil)
	m.Filter(nil).Convert(nil).Track(nil)

	m.Keys().For(nil)
	m.Keys().ForEach(nil)
	m.Values().For(nil)
	m.Values().ForEach(nil)
	// m.Values().Convert(nil).For(nil)
	m.Values().Filter(nil).ForEach(nil)
}

func Test_Map_Zero(t *testing.T) {
	var m mutable.Map[string, string]

	m.Set("a", "A")
	assert.True(t, m.SetNew("b", "B"))
	assert.False(t, m.SetNew("b", "B"))

	m.SetMap(nil)
	m.SetMap(map_.Of(k.V("d", "D")))

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
	m.Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(_, _ string) error { return nil })
	m.Filter(func(_, _ string) bool { return true }).Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(_, _ string) error { return nil })

	m.Keys().For(func(_ string) error { return nil })
	m.Keys().ForEach(func(_ string) {})
	m.Values().For(func(_ string) error { return nil })
	m.Values().ForEach(func(_ string) {})
	m.Values().Convert(as.Is[string]).ForEach(func(_ string) {})
	m.Values().Filter(func(_ string) bool { return true }).ForEach(func(_ string) {})
}

func Test_Map_new(t *testing.T) {
	var m = new(mutable.Map[string, string])

	m.Set("a", "A")
	assert.True(t, m.SetNew("b", "B"))
	assert.False(t, m.SetNew("b", "B"))

	m.SetMap(nil)
	m.SetMap(map_.Of(k.V("d", "D")))

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
	m.Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(_, _ string) error { return nil })

	m.Filter(func(_, _ string) bool { return true }).Convert(func(s1, s2 string) (string, string) { return s1, s2 }).Track(func(_, _ string) error { return nil })

	m.Keys().For(func(_ string) error { return nil })
	m.Keys().ForEach(func(_ string) {})
	m.Values().For(func(_ string) error { return nil })
	m.Values().ForEach(func(_ string) {})
	// m.Values().Convert(as.Is[string]).For(func(_ string) error { return nil })
	m.Values().Filter(func(_ string) bool { return true }).ForEach(func(_ string) {})
}

func Test_Map_Sort(t *testing.T) {
	var m = new(mutable.Map[int, string])

	m.Set(5, "5")
	m.Set(4, "4")
	m.Set(-8, "-8")
	m.Set(10, "10")

	o := m.Sort(op.Compare)

	assert.Equal(t, ordered.NewMap(k.V(-8, "-8"), k.V(4, "4"), k.V(5, "5"), k.V(10, "10")), o)
	assert.NotSame(t, m, o)
}
