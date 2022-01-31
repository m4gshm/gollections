//package set provides the unordered set container implementation
package set

func Of[T comparable](elements ...T) *Set[T] {
	return Convert(elements)
}

func New[T comparable](elements []T) *Set[T] {
	return Convert(elements)
}
