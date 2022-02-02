package slice

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

type Fit[T any] struct {
	Elements []T
	By       c.Predicate[T]

	current T
	i       int
	err     error
}

var _ c.Iterator[any] = (*Fit[any])(nil)

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

// func nextFiltered[T any](iter c.Iterator[T], fit c.Predicate[T]) (T, bool) {
// 	for iter.HasNext() {
// 		if v := iter.Get(); fit(v) {
// 			return v, true
// 		}
// 	}
// 	var v T
// 	return v, false
// }
