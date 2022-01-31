//package vector provides the Vector (ordered) implementation
package vector

func Of[T any](elements ...T) *Vector[T] {
	return Convert(elements)
}

func New[T any](elements []T) *Vector[T] {
	return Convert(elements)
}
