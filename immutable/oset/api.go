//Package oset provides the ordered set container implementation
package oset

import "github.com/m4gshm/gollections/immutable/ordered"

func Of[T comparable](elements ...T) *ordered.Set[T] {
	return ordered.NewSet(elements)
}

func New[T comparable](elements []T) *ordered.Set[T] {
	return ordered.NewSet(elements)
}
