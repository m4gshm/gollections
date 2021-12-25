package iter

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/typ"
)

type Fit[T any] struct {
	Iter typ.Iterator[T]
	By   check.Predicate[T]

	current T
}

var _ typ.Iterator[interface{}] = (*Fit[interface{}])(nil)

func (s *Fit[T]) HasNext() bool {
	v, ok := nextFiltered(s.Iter, s.By)
	s.current = v
	return ok
}

func (s *Fit[T]) Get() T {
	return s.current
}

func nextFiltered[T any](iter typ.Iterator[T], filter check.Predicate[T]) (T, bool) {
	for iter.HasNext() {
		if v := iter.Get(); filter(v) {
			return v, true
		}
	}
	var v T
	return v, false
}
