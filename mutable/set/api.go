package set

func Of[T comparable](elements ...T) *Set[T] {
	return Convert(elements)
}

func Empty[T comparable]() *Set[T] {
	return New[T](0)
}

func New[T comparable](capacity int) *Set[T] {
	return Wrap(make(map[T]struct{}, capacity))
}
