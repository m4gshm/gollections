package sync

import (
	"sync"

	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/typ"
)

type Map[k comparable, v any] struct {
	m sync.Map
}

var _ mutable.Settable[any, any] = (*Map[any, any])(nil)
var _ mutable.Deleteable[any] = (*Map[any, any])(nil)
var _ mutable.Removable[any, any] = (*Map[any, any])(nil)
var _ typ.TrackEach[any, any] = (*Map[any, any])(nil)
var _ typ.Access[any, any] = (*Map[any, any])(nil)

func (m *Map[k, v]) TrackEach(traker func(key k, value v)) {
	m.m.Range(func(key, value any) bool {
		traker(key.(k), value.(v))
		return true
	})
}

func (m *Map[k, v]) Set(key k, value v) (bool, error) {
	m.m.Store(key, value)
	return true, nil
}

func (m *Map[k, v]) Get(key k) (v, bool) {
	value, ok := m.m.Load(key)
	return value.(v), ok
}

func (m *Map[k, v]) Delete(keys ...k) (bool, error) {
	for _, key := range keys {
		m.m.Delete(key)
	}
	return true, nil
}

func (m *Map[k, v]) Remove(key k) (v, bool, error) {
	rawVal, ok := m.m.LoadAndDelete(key)
	return rawVal.(v), ok, nil
}