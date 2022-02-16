package sync

import (
	"sync"

	"github.com/m4gshm/gollections/c"
)

//Map is typed wrapper of sync.Map
type Map[k comparable, v any] struct {
	m sync.Map
}

var (
	_ c.Settable[int, any]  = (*Map[int, any])(nil)
	_ c.Deleteable[int]     = (*Map[int, any])(nil)
	_ c.Removable[int, any] = (*Map[int, any])(nil)
	_ c.TrackEach[any, int] = (*Map[int, any])(nil)
	_ c.Access[int, any]    = (*Map[int, any])(nil)
)

func (m *Map[k, v]) TrackEach(traker func(key k, value v)) {
	m.m.Range(func(key, value any) bool {
		traker(key.(k), value.(v))
		return true
	})
}

func (m *Map[k, v]) Set(key k, value v) bool {
	m.m.Store(key, value)
	return true
}

func (m *Map[k, v]) Get(key k) (v, bool) {
	value, ok := m.m.Load(key)
	return value.(v), ok
}

func (m *Map[k, v]) Delete(keys ...k) bool {
	for _, key := range keys {
		m.m.Delete(key)
	}
	return true
}

func (m *Map[k, v]) Remove(key k) (v, bool) {
	rawVal, ok := m.m.LoadAndDelete(key)
	return rawVal.(v), ok
}
