package iterator

import (
	"github.com/m4gshm/container/check"
)

type FilterIter[T any] struct {
	iter    Iterator[T]
	filters []check.Predicate[T]
	current T
}

var _ Iterator[interface{}] = (*FilterIter[interface{}])(nil)

func (s *FilterIter[T]) Next() (T, bool) {
	var (
		v  T
		ok = true
	)
	for ok {
		if v, ok = s.iter.Next(); ok {
			if fit:= check.IsFit(v, s.filters...); fit {
				return v, true
			}
			
		}
	}
	return v, ok
}

func (s *FilterIter[T]) HasNext() bool {
	v, ok := s.Next()
	s.current = v
	return ok
}

func (s *FilterIter[T]) Get() T {
	return s.current
}

func filter[T any](iter Iterator[T], filters []check.Predicate[T]) (T, bool) {
	var (
		v  T
		ok = true
	)
	for ok {
		if v, ok = iter.Next(); ok {
			if fit := check.IsFit(v, filters...); fit {
				return v, true
			}
		}
	}
	return v, ok
}
