package iter

import "github.com/m4gshm/container/typ"

func SliceWrapper[T any](elements []T) *Slice[T] {
	return &Slice[T]{Elements: elements}
}

type Slice[T any] struct {
	Elements []T
	current  T
	i        int
}

var _ typ.Iterator[interface{}] = (*Slice[interface{}])(nil)
var _ typ.Resetable = (*Slice[interface{}])(nil)

func (s *Slice[T]) HasNext() bool {
	var v T
	e := s.Elements
	i := s.i
	if i < len(e) {
		v = e[i]
		s.current = v
		s.i = i + 1
		return true
	}
	return false
}

func (s *Slice[T]) Get() T {
	return s.current
}

func (s *Slice[T]) Reset() {
	var v T
	s.current = v
	s.i = 0
}
