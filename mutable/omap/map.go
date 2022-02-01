package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func AsMap[k comparable, v any](elements []*typ.KV[k, v]) *Map[k, v] {
	var (
		uniques = make(map[k]v, 0)
		order   = make([]k, 0, 0)
	)
	for _, kv := range elements {
		key := kv.Key()
		val := kv.Value()
		if _, ok := uniques[key]; !ok {
			order = append(order, key)
			uniques[key] = val
		}
	}
	return WrapMap(order, uniques)
}

func ToMap[k comparable, v any](elements map[k]v) *Map[k, v] {
	var (
		uniques = make(map[k]v, len(elements))
		order   = make([]k, len(elements))
	)
	for key, val := range elements {
		order = append(order, key)
		uniques[key] = val
	}
	return WrapMap(order, uniques)
}

func WrapMap[k comparable, v any](order []k, uniques map[k]v) *Map[k, v] {
	return &Map[k, v]{elements: order, uniques: uniques}
}

//Map provides access to elements by key.
type Map[k comparable, v any] struct {
	changeMark int32
	elements   []k
	uniques    map[k]v
	err        error
}

var (
	_ mutable.Settable[any, any] = (*Map[any, any])(nil)
	_ typ.Map[any, any]          = (*Map[any, any])(nil)
	_ fmt.Stringer               = (*Map[any, any])(nil)
)

func (s *Map[k, v]) Begin() typ.KVIterator[k, v] {
	return s.Iter()
}

func (s *Map[k, v]) Iter() *it.OrderedKV[k, v] {
	return it.NewOrderedKV(s.elements, s.uniques)
}

func (s *Map[k, v]) Collect() map[k]v {
	e := s.uniques
	out := make(map[k]v, len(e))
	for key, val := range e {
		out[key] = val
	}
	return out
}

func (s *Map[k, v]) Len() int {
	return len(s.elements)
}

func (s *Map[k, v]) For(walker func(*typ.KV[k, v]) error) error {
	return map_.ForOrdered(s.elements, s.uniques, walker)
}

func (s *Map[k, v]) ForEach(walker func(*typ.KV[k, v])) {
	map_.ForEachOrdered(s.elements, s.uniques, walker)
}

func (s *Map[k, v]) Track(tracker func(k, v) error) error {
	return map_.TrackOrdered(s.elements, s.uniques, tracker)
}

func (s *Map[k, v]) TrackEach(tracker func(k, v)) {
	map_.TrackEachOrdered(s.elements, s.uniques, tracker)
}

func (s *Map[k, v]) Contains(key k) bool {
	_, ok := s.uniques[key]
	return ok
}

func (s *Map[k, v]) Get(key k) (v, bool) {
	val, ok := s.uniques[key]
	return val, ok
}

func (s *Map[k, v]) Set(key k, value v) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	u := s.uniques
	if _, ok := u[key]; !ok {
		markOnStart := s.changeMark
		e := s.elements
		u[key] = value
		s.elements = append(e, key)
		return mutable.Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *Map[k, v]) Keys() typ.Collection[k, []k, typ.Iterator[k]] {
	return vector.New(s.elements)
}

func (s *Map[k, v]) Values() typ.Collection[v, []v, typ.Iterator[v]] {
	return ordered.WrapVal(s.elements, s.uniques)
}

func (s *Map[k, v]) String() string {
	return map_.ToStringOrdered(s.elements, s.uniques)
}

func (s *Map[k, v]) FilterKey(fit typ.Predicate[k]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), func(key k, val v) bool { return fit(key) }), collect.Map[k, v])
}

func (s *Map[k, v]) MapKey(by typ.Converter[k, k]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), func(key k, val v) (k, v) { return by(key), val }), collect.Map[k, v])
}

func (s *Map[k, v]) FilterValue(fit typ.Predicate[v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), func(key k, val v) bool { return fit(val) }), collect.Map[k, v])
}

func (s *Map[k, v]) MapValue(by typ.Converter[v, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), func(key k, val v) (k, v) { return key, by(val) }), collect.Map[k, v])
}

func (s *Map[k, v]) Filter(filter typ.BiPredicate[k, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), filter), collect.Map[k, v])
}

func (s *Map[k, v]) Map(by typ.BiConverter[k, v, k, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), by), collect.Map[k, v])
}

func (s *Map[k, v]) Reduce(by op.Quaternary[k, v]) (k, v) {
	return it.ReduceKV(s.Iter(), by)
}
