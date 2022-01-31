package vector

func Of[T any](elements ...T) *Vector[T] {
	return Convert(elements)
}

func Empty[T any]() *Vector[T] {
	return New[T](0)
}

func New[T any](capacity int) *Vector[T] {
	return Create[T](capacity)
}
