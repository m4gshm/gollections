package iterator

type SliceIter[T any] struct {
	values  []T
	current T
	i       int
}

var _ Iterator[interface{}] = (*SliceIter[interface{}])(nil)

func (s *SliceIter[T]) Next() bool {
	if s.i < len(s.values) {
		s.current = s.values[s.i]
		s.i++
		return true
	}
	return false
}

func (s *SliceIter[T]) Get() T {
	return s.current
}
