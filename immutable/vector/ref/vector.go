package ref

import (
	"fmt"

	"github.com/m4gshm/container/immutable"
	"github.com/m4gshm/container/immutable/vector"
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
)

func Convert[T any](elements []*T) *Vector[T] {
	return Wrap(slice.Copy(elements))
}

func Wrap[T any](elements []*T) *Vector[T] {
	return &Vector[T]{vector.Wrap(elements)}
}

type Vector[T any] struct {
	*vector.Vector[*T]
}

var _ immutable.Vector[any, typ.Iterator[any]] = (*Vector[any])(nil)
var _ fmt.Stringer = (*Vector[any])(nil)

func (s *Vector[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *Vector[T]) Iter() *iter.RefIter[T] {
	return &iter.RefIter[T]{Iterator: s.Vector.Iter()}
}

func (s *Vector[T]) Elements() []T {
	refs := s.Vector.Elements()
	elements := make([]T, len(refs))
	for i, r := range refs {
		elements[i] = *r
	}
	return elements
}

func (s *Vector[T]) ForEach(walker func(T)) {
	s.Vector.ForEach(func(ref *T) { walker(*ref) })
}

func (s *Vector[T]) Get(index int) (T, bool) {
	if r, ok := s.Vector.Get(index); ok {
		return *r, true
	}
	var no T
	return no, false
}

func (s *Vector[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, typ.Iterator[T]] {
	return iter.NewPipe[T](iter.Filter(s.Iter(), filter))
}

func (s *Vector[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, typ.Iterator[T]] {
	return iter.NewPipe[T](iter.Map(s.Iter(), by))
}

func (s *Vector[T]) Reduce(by op.Binary[T]) T {
	return iter.Reduce(s.Iter(), by)
}

func (s *Vector[T]) String() string {
	return slice.ToString(s.Elements())
}
