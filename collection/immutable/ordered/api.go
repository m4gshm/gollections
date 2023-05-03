// Package ordered provides immutable ordered collection implementations
package ordered

import "github.com/m4gshm/gollections/slice/clone"

// NewMap instantiates Map populated by the 'elements' map key/values
func NewMap[K comparable, V any](order []K, elements map[K]V) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for _, key := range order {
		uniques[key] = elements[key]
	}
	return WrapMap(clone.Of(order), uniques)
}

// NewSet instantiates Set and copies elements to it
func NewSet[T comparable](elements []T) Set[T] {
	var (
		l       = len(elements)
		uniques = make(map[T]struct{}, l)
		order   = make([]T, 0, l)
	)
	for _, e := range elements {
		order = add(e, uniques, order)
	}
	return WrapSet(order, uniques)
}
