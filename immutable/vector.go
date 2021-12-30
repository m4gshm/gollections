package immutable

import (
	"fmt"

	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
)


func New[T any](values []T) *Slice[T] {
	elements := make([]T, len(values))
	copy(elements, values)
	return Wrap(elements)
}

func Wrap[T any](elements []T) *Slice[T] {
	return &Slice[T]{elements: elements}
}

type Slice[T any] struct {
	elements   []T
	changeMark int32
}

var _ Vector[any] = (*Slice[any])(nil)
var _ fmt.Stringer = (*Slice[any])(nil)

func (s *Slice[T]) Begin() typ.Iterator[T] {
	return s.newIter()
}

func (s *Slice[T]) newIter() *iter.Iter[T] {
	return iter.New(&s.elements)
}

func (s *Slice[T]) Values() []T {
	elements := make([]T, len(s.elements))
	copy(elements, elements)
	return elements
}

func (s *Slice[T]) ForEach(walker func(T)) {
	for _, e := range s.elements {
		walker(e)
	}
}

func (s *Slice[T]) Len() int {
	return len(s.elements)
}

func (s *Slice[T]) Get(index int) (T, bool) {
	elements := s.elements
	l := len(elements)
	if l > 0 && (index >= 0 || index < l) {
		return s.elements[index], true
	}
	var no T
	return no, false
}

func (s *Slice[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T] {
	return iter.NewPipe[T](iter.Filter(s.newIter(), filter))
}

func (s *Slice[T]) Map(by typ.Converter[T, T]) typ.Pipe[T] {
	return iter.NewPipe[T](iter.Map(s.newIter(), by))
}

func (s *Slice[T]) Reduce(by op.Binary[T]) T {
	return iter.Reduce(s.newIter(), by)
}

func (s *Slice[T]) String() string {
	return slice.ToString(s.elements)
}
