// Package ordered provides mutable ordered collection implementations
package ordered

import "github.com/m4gshm/gollections/slice/clone"

// NewMap instantiates Map populated by the 'elements' map key/values
func NewMap[K comparable, V any](order []K, elements map[K]V) *Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for _, key := range order {
		uniques[key] = elements[key]
	}
	return WrapMap(clone.Of(order), uniques)
}

// NewSet instantiates Set and copies elements to it
func NewSet[T comparable](elements []T) *Set[T] {
	var (
		l       = len(elements)
		uniques = make(map[T]int, l)
		order   = make([]T, 0, l)
	)
	pos := 0
	for _, e := range elements {
		if _, ok := uniques[e]; !ok {
			order = append(order, e)
			uniques[e] = pos
			pos++
		}
	}
	return WrapSet(order, uniques)
}
