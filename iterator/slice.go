package iterator

func New[T any](elements ...T) *Slice[T] {
	return &Slice[T]{elements: elements}
}

type Slice[T any] struct {
	Iterator[T]
	elements []T
	current  T
	i        int
}

var _ Iterator[interface{}] = (*Slice[interface{}])(nil)

func (s *Slice[T]) Next() (T, bool) {
	var v T
	e := s.elements
	i := s.i
	if i < len(e) {
		v = e[i]
		s.current = v
		s.i = i + 1
		return v, true
	}
	return v, false
}

func (s *Slice[T]) HasNext() bool {
	v, ok := s.Next()
	s.current = v
	return ok
}

func (s *Slice[T]) Get() T {
	return s.current
}
