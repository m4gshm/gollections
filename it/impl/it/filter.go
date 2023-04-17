package it

import (
	"github.com/m4gshm/gollections/c"
)

// Fit is the Iterator wrapper that provides filtering of elements by a Predicate.
type Fit[T, IT any] struct {
	iter IT
	next func() (T, bool)
	by   func(T) bool
}

var (
	_ c.Iterator[any] = (*Fit[any, any])(nil)
	_ c.Iterator[any] = Fit[any, any]{}
)

func (s Fit[T, IT]) Next() (T, bool) {
	return nextFiltered(s.next, s.by)
}

// FitKV is the KVIterator wrapper that provides filtering of key/value elements by a Predicate.
type FitKV[K, V any, IT c.KVIterator[K, V]] struct {
	iter IT
	by   func(K, V) bool
}

var (
	_ c.KVIterator[any, any] = (*FitKV[any, any, c.KVIterator[any, any]])(nil)
	_ c.KVIterator[any, any] = FitKV[any, any, c.KVIterator[any, any]]{}
)

func (s FitKV[K, V, IT]) Next() (K, V, bool) {
	return nextFilteredKV(s.iter, s.by)
}

func nextFiltered[T any](next func() (T, bool), filter func(T) bool) (T, bool) {
	for v, ok := next(); ok; v, ok = next() {
		if filter(v) {
			return v, true
		}
	}
	var v T
	return v, false
}

func nextFilteredKV[K any, V any, IT c.KVIterator[K, V]](iter IT, filter func(K, V) bool) (K, V, bool) {
	for key, val, ok := iter.Next(); ok; key, val, ok = iter.Next() {
		if filter(key, val) {
			return key, val, true
		}
	}
	var key K
	var val V
	return key, val, false
}
