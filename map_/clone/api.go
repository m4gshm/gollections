package clone

import "github.com/m4gshm/gollections/map_"

// Of - synonym of the map_.Clone
func Of[M ~map[K]V, K comparable, V any](elements M) M {
	return map_.Clone(elements)
}

// Deep - synonym of the map_.DeepClone
func Deep[M ~map[K]V, K comparable, V any](elements M, valCopier func(V) V) M {
	return map_.DeepClone(elements, valCopier)
}
