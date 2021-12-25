package iter

import "github.com/m4gshm/container/typ"

func New[T any](elements []T) *Iter[T] {
	return &Iter[T]{elements: elements}
}

func NewReseteable[T any](elements []T) *ResIter[T] {
	return &ResIter[T]{&Iter[T]{elements: elements}}
}

type Iter[T any] struct {
	elements []T
	current  T
	i        int
}

var _ typ.Iterator[interface{}] = (*Iter[interface{}])(nil)

func (s *Iter[T]) HasNext() bool {
	var v T
	e := s.elements
	i := s.i
	if i < len(e) {
		v = e[i]
		s.current = v
		s.i = i + 1
		return true
	}
	return false
}

func (s *Iter[T]) Get() T {
	return s.current
}

type ResIter[T any] struct {
	*Iter[T]
}

var _ typ.Resetable = (*ResIter[interface{}])(nil)

func (s *ResIter[T]) Reset() {
	var v T
	s.current = v
	s.i = 0
}
