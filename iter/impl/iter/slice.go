package iter

import "github.com/m4gshm/container/typ"

func New[T any](elements []T) *Iter[T] {
	return &Iter[T]{elements: elements}
}

func NewReseteable[T any](elements []T) *ResIter[T] {
	return &ResIter[T]{New(elements)}
}

type Iter[T any] struct {
	elements []T
	next     int
}

var _ typ.Iterator[interface{}] = (*Iter[interface{}])(nil)

func (s *Iter[T]) HasNext() bool {
	e := s.elements
	l := len(e)
	if l == 0 {
		return false
	}
	next := s.next
	if next < l {
		return true
	}
	return false
}

func (s *Iter[T]) Get() T {
	current := s.next
	s.next++
	return s.elements[current]
}

type ResIter[T any] struct {
	*Iter[T]
}

var _ typ.Resetable = (*ResIter[interface{}])(nil)

func (s *ResIter[T]) Reset() {
	s.next = 0
}
