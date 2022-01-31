//package oset provides the ordered set container implementation
package oset

func Of[T comparable](elements ...T) *OrderedSet[T] {
	return Convert(elements)
}

func New[T comparable](elements []T) *OrderedSet[T] {
	return Convert(elements)
}
