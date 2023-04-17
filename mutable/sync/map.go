package sync

import (
	"sync"

	"github.com/m4gshm/gollections/c"
)

// NewMap sync Map constructor
func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{}
}

// Map is typed wrapper of sync.Map
type Map[K comparable, V any] struct {
	m sync.Map
}

var (
	_ c.Settable[int, any]      = (*Map[int, any])(nil)
	_ c.Deleteable[int]         = (*Map[int, any])(nil)
	_ c.Removable[int, any]     = (*Map[int, any])(nil)
	_ c.TrackEachLoop[any, int] = (*Map[int, any])(nil)
	_ c.Access[int, any]        = (*Map[int, any])(nil)
)

func (m *Map[K, V]) TrackEach(traker func(key K, value V)) {
	m.m.Range(func(key, value any) bool {
		traker(key.(K), value.(V))
		return true
	})
}

func (m *Map[K, V]) Set(key K, value V) {
	m.m.Store(key, value)
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	value, ok := m.m.Load(key)
	return value.(V), ok
}

func (m *Map[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		m.m.Delete(key)
	}
}

func (m *Map[K, V]) DeleteOne(key K) {
	m.m.Delete(key)
}

func (m *Map[K, V]) Remove(key K) (V, bool) {
	rawVal, ok := m.m.LoadAndDelete(key)
	return rawVal.(V), ok
}
