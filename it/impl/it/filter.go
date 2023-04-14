package it

import (
	"github.com/m4gshm/gollections/c"
)

// Fit is the Iterator wrapper that provides filtering of elements by a Predicate.
type Fit[T any, IT c.Iterator[T]] struct {
	iter IT
	by   func(T) bool
}

var (
	_ c.Iterator[any] = (*Fit[any, c.Iterator[any]])(nil)
	_ c.Iterator[any] = Fit[any, c.Iterator[any]]{}
)

func (s Fit[T, IT]) Next() (T, bool) {
	return nextFiltered(s.iter, s.by)
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

func nextFiltered[T any, IT c.Iterator[T]](iter IT, filter func(T) bool) (T, bool) {
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
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
