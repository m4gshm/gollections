package iterator

type Slice[T any] struct {
	values  []T
	current T
	i       int
}

var _ Iterator[interface{}] = (*Slice[interface{}])(nil)

func (s *Slice[T]) Next() bool {
	if s.i < len(s.values) {
		s.current = s.values[s.i]
		s.i++
		return true
	}
	return false
}

func (s *Slice[T]) Get() T {
	return s.current
}
