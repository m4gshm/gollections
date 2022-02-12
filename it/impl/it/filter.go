package it

import "github.com/m4gshm/gollections/c"

type Fit[T any] struct {
	Iter c.Iterator[T]
	By   c.Predicate[T]

	current T
	err     error
}

var _ c.Iterator[any] = (*Fit[any])(nil)

func (s *Fit[T]) HasNext() bool {
	v, ok := nextFiltered(s.Iter, s.By)
	s.current = v
	return ok
}

func (s *Fit[T]) Get() T {
	return s.current
}

type FitKV[k, v any] struct {
	Iter c.KVIterator[k, v]
	By   c.BiPredicate[k, v]

	currentK k
	currentV v
}

var _ c.KVIterator[any, any] = (*FitKV[any, any])(nil)

func (s *FitKV[k, v]) HasNext() bool {
	key, val, ok := nextFilteredKV(s.Iter, s.By)
	s.currentK = key
	s.currentV = val
	return ok
}

func (s *FitKV[k, v]) Get() (k, v) {
	return s.currentK, s.currentV
}

func nextFiltered[T any](iter c.Iterator[T], filter c.Predicate[T]) (T, bool) {
	for iter.HasNext() {
		if v := iter.Get(); filter(v) {
			return v, true
		}
	}
	var v T
	return v, false
}

func nextFilteredKV[k any, v any](iter c.KVIterator[k, v], filter c.BiPredicate[k, v]) (k, v, bool) {
	for iter.HasNext() {
		if k, v := iter.Get(); filter(k, v) {
			return k, v, true
		}
	}
	var key k
	var val v
	return key, val, false
}
