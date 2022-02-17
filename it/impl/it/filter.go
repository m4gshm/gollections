package it

import "github.com/m4gshm/gollections/c"

type Fit[T any, IT c.Iterator[T]] struct {
	Iter IT
	By   c.Predicate[T]

	current T
	err     error
}

var _ c.Iterator[any] = (*Fit[any, c.Iterator[any]])(nil)

func (s *Fit[T, IT]) HasNext() bool {
	V, ok := nextFiltered[T](s.Iter, s.By)
	s.current = V
	return ok
}

func (s *Fit[T, IT]) Next() T {
	return s.current
}

type FitKV[K, V any, IT c.KVIterator[K, V]] struct {
	Iter IT
	By   c.BiPredicate[K, V]

	currentK K
	currentV V
}

var _ c.KVIterator[any, any] = (*FitKV[any, any, c.KVIterator[any, any]])(nil)

func (s *FitKV[K, V, IT]) HasNext() bool {
	key, val, ok := nextFilteredKV[K, V](s.Iter, s.By)
	s.currentK = key
	s.currentV = val
	return ok
}

func (s *FitKV[K, V, IT]) Next() (K, V) {
	return s.currentK, s.currentV
}

func nextFiltered[T any, IT c.Iterator[T], F c.Predicate[T]](iter IT, filter F) (T, bool) {
	for iter.HasNext() {
		if V := iter.Next(); filter(V) {
			return V, true
		}
	}
	var V T
	return V, false
}

func nextFilteredKV[K any, V any, IT c.KVIterator[K, V], F c.BiPredicate[K, V]](iter IT, filter F) (K, V, bool) {
	for iter.HasNext() {
		if K, V := iter.Next(); filter(K, V) {
			return K, V, true
		}
	}
	var key K
	var val V
	return key, val, false
}
