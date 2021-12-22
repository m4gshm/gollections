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

func (s *FilterIter[T]) HasNext() bool {
	v, ok := filterNext(s.iter, s.filters)
	s.current = v
	return ok
}

func (s *FilterIter[T]) Get() T {
	return s.current
}

func filterNext[T any](iter Iterator[T], filters []check.Predicate[T]) (T, bool) {
	for iter.HasNext() {
		v := iter.Get()
		if ok := check.IsFit(v, filters...); ok {
			return v, true
		}
	}
	var v T
	return v, false
}
