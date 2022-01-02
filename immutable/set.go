package immutable

import (
	"fmt"

	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
)

func NewOrderedSet[T comparable](values []T) *OrderedSet[T] {
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
var _ fmt.GoStringer = (*OrderedSet[any])(nil)

func (s *OrderedSet[T]) Begin() typ.Iterator[T] {
	return &OrderIter[T]{s.newIter()}
}

func (s *OrderedSet[T]) newIter() *iter.Iter[*T] {
	return iter.New(s.elements)
}

func (s *OrderedSet[T]) Values() []T {
	e := s.elements
	out := make([]T, len(e))
	for i, v := range e {
		out[i] = *v
	}
	return out
}

func (s *OrderedSet[T]) ForEach(walker func(T)) {
	for _, e := range s.elements {
		walker(*e)
	}
}

func (s *OrderedSet[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T] {
	return iter.NewPipe[T](&OrderIter[T]{iter.Filter(s.newIter(), func(ref *T) bool { return filter(*ref) })})
}

func (s *OrderedSet[T]) Map(by typ.Converter[T, T]) typ.Pipe[T] {
	return iter.NewPipe[T](&OrderIter[T]{iter.Map(s.newIter(), func(ref *T) *T {
		conv := by(*ref)
		return &conv
	})})
}

func (s *OrderedSet[T]) Reduce(by op.Binary[T]) T {
	return iter.Reduce(&OrderIter[T]{s.newIter()}, by)
}

func (s *OrderedSet[T]) Len() int {
	return len(s.elements)
}

func (s *OrderedSet[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *OrderedSet[T]) String() string {
	return s.GoString()
}

func (s *OrderedSet[T]) GoString() string {
	return slice.ToStringRefs(s.elements)
}
