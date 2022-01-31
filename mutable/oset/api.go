package oset

// "github.com/m4gshm/gollections/mutable"

func Of[T comparable](elements ...T) *OrderedSet[T] {
	return Convert(elements)
}

func Empty[T comparable]() *OrderedSet[T] {
	return New[T](0)
}

func New[T comparable](capacity int) *OrderedSet[T] {
	return NewOrderedSet[T](capacity)
}
