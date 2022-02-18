package it

import "github.com/m4gshm/gollections/c"

type Fit[T any, IT c.Iterator[T]] struct {
	Iter IT
	By   c.Predicate[T]

	current T
	err     error
}

var _ c.Iterator[any] = (*Fit[any, c.Iterator[any]])(nil)

func (s *Fit[T, IT]) GetNext() (T, bool) {
	return nextFiltered(s.Iter, s.By)
}

type FitKV[K, V any, IT c.KVIterator[K, V]] struct {
	Iter IT
	By   c.BiPredicate[K, V]

	currentK K
	currentV V
}

var _ c.KVIterator[any, any] = (*FitKV[any, any, c.KVIterator[any, any]])(nil)

func (s *FitKV[K, V, IT]) GetNext() (K, V, bool) {
	return nextFilteredKV(s.Iter, s.By)
}

func nextFiltered[T any, IT c.Iterator[T], F c.Predicate[T]](iter IT, filter F) (T, bool) {
	for v, ok := iter.GetNext(); ok; v, ok = iter.GetNext() {
		if filter(v) {
			return v, true
		}
	}
	var v T
	return v, false
}

func nextFilteredKV[K any, V any, IT c.KVIterator[K, V], F c.BiPredicate[K, V]](iter IT, filter F) (K, V, bool) {
	for key, val, ok := iter.GetNext(); ok; key, val, ok = iter.GetNext() {
		if filter(key, val) {
			return key, val, true
		}

	}
	var key K
	var val V
	return key, val, false
}
