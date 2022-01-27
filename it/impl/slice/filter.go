package slice

import (
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/typ"
)

type Fit[T any] struct {
	Elements []T
	By       typ.Predicate[T]

	current T
	i       int
	err     error
}

var _ typ.Iterator[any] = (*Fit[any])(nil)

func (s *Fit[T]) HasNext() bool {
	v, ok := nextArrayElem(s.Elements, s.By, &s.i)
	if ok {
		s.current = v
	} else {
		s.err = it.Exhausted
	}
	return ok
}

func (s *Fit[T]) Get() (T, error) {
	return s.current, s.err
}

func (s *Fit[T]) Next() T {
	return it.Next[T](s)
}

// func nextFiltered[T any](iter typ.Iterator[T], fit typ.Predicate[T]) (T, bool) {
// 	for iter.HasNext() {
// 		if v := iter.Get(); fit(v) {
// 			return v, true
// 		}
// 	}
// 	var v T
// 	return v, false
// }