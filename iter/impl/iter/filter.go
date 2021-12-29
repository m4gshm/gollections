package iter

import (
	"github.com/m4gshm/container/typ"
)

type Fit[T any] struct {
	Iter typ.Iterator[T]
	By   typ.Predicate[T]

	current T
	err     error
}

var _ typ.Iterator[interface{}] = (*Fit[interface{}])(nil)

func (s *Fit[T]) HasNext() bool {
	v, ok, err := nextFiltered(s.Iter, s.By)
	if err != nil {
		s.err = err
	} else {
		s.current = v
	}
	return ok
}

func (s *Fit[T]) Get() T {
	if err := s.err; err != nil {
		panic(err)
	}
	return s.current
}

func (s *Fit[T]) Err() error {
	return s.err
}

func nextFiltered[T any](iter typ.Iterator[T], filter typ.Predicate[T]) (T, bool, error) {
	for iter.HasNext() {
		if err := iter.Err(); err != nil {
			var no T
			return no, true, err
		}
		if v := iter.Get(); filter(v) {
			return v, true, nil
		}
	}
	var v T
	return v, false, nil
}
