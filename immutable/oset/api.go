//package oset provides the ordered set container implementation
package oset

import (
	"github.com/m4gshm/gollections/immutable"
)

func Of[T comparable](elements ...T) immutable.Set[T] {
	return Convert(elements)
}

func New[T comparable](elements []T) immutable.Set[T] {
	return Convert(elements)
}
