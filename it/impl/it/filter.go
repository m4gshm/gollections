package it

import (
	"github.com/m4gshm/gollections/typ"
)

type Fit[T any] struct {
	Iter typ.Iterator[T]
	By   typ.Predicate[T]

	current T
	err     error
}

var _ typ.Iterator[any] = (*Fit[any])(nil)

func (s *Fit[T]) HasNext() bool {
	v, ok, err := nextFiltered(s.Iter, s.By)
	if err != nil {
		s.err = err
	} else {
		s.current = v
	}
	return ok
}

func (s *Fit[T]) Get() (T, error) {
	return s.current, s.err
}

func (s *Fit[T]) Next() T {
	return Next[T](s)
}

type FitKV[k, v any] struct {
	Iter typ.KVIterator[k, v]
	By   typ.BiPredicate[k, v]

	currentK k
	currentV v
	err      error
}

var _ typ.KVIterator[any, any] = (*FitKV[any, any])(nil)

func (s *FitKV[k, v]) HasNext() bool {
	key, val, ok, err := nextFilteredKV(s.Iter, s.By)
	if err != nil {
		s.err = err
	} else {
		s.currentK = key
		s.currentV = val
	}
	return ok
}

func (s *FitKV[k, v]) Get() (k, v, error) {
	return s.currentK, s.currentV, s.err
}

func (s *FitKV[k, v]) Next() (k, v) {
	return NextKV[k, v](s)
}

func nextFiltered[T any](iter typ.Iterator[T], filter typ.Predicate[T]) (T, bool, error) {
	for iter.HasNext() {
		if v, err := iter.Get(); err != nil {
			var no T
			return no, true, err
		} else if filter(v) {
			return v, true, nil
		}
	}
	var v T
	return v, false, nil
}

func nextFilteredKV[k any, v any](iter typ.KVIterator[k, v], filter typ.BiPredicate[k, v]) (k, v, bool, error) {
	for iter.HasNext() {
		if k, v, err := iter.Get(); err != nil {

			return k, v, true, err
		} else if filter(k, v) {
			return k, v, true, nil
		}
	}
	var key k
	var val v
	return key, val, false, nil
}
