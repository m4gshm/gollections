package immutable

import (
	"fmt"

	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
)

func newSet[T comparable](values []T) *OrderedSet[T] {
	var (
		uniques = make(map[T]struct{}, 0)
		order   = make([]*T, 0, 0)
	)
	for _, v := range values {
		if _, ok := uniques[v]; !ok {
			vv := v
			order = append(order, &vv)
			uniques[vv] = struct{}{}
		}
	}
	return &OrderedSet[T]{elements: order, uniques: uniques}
}

type OrderedSet[T comparable] struct {
	elements []*T
	uniques  map[T]struct{}
}

var _ Set[any] = (*OrderedSet[any])(nil)
var _ fmt.Stringer = (*OrderedSet[any])(nil)

func (s *OrderedSet[T]) Begin() typ.Iterator[T] {
	return &OrderIter[T]{Iter: iter.New(&s.elements)}
}

func (s *OrderedSet[T]) Values() []T {
	e := s.elements
	out := make([]T, len(e))
	for i, v := range e {
		out[i] = *v
	}
	return out
}

func (s *OrderedSet[T]) Len() int {
	return len(s.elements)
}

func (s *OrderedSet[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *OrderedSet[T]) String() string {
	return slice.ToString(s.elements)
}
