package slice

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/typ"
)

type Fit[T any] struct {
	Elements []T
	By   check.Predicate[T]

	current  T
	i        int
}

var _ typ.Iterator[interface{}] = (*Fit[interface{}])(nil)

func (s *Fit[T]) HasNext() bool {
	v, ok := nextArrayElem(s.Elements, s.By, &s.i)
	s.current = v
	return ok
}

func (s *Fit[T]) Get() T {
	return s.current
}

func nextFiltered[T any](iter typ.Iterator[T], fit check.Predicate[T]) (T, bool) {
	for iter.HasNext() {
		if v := iter.Get(); fit(v) {
			return v, true
		}
	}
	var v T
	return v, false
}
