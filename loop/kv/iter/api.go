package iter

// Convert creates an Iterator that applies a transformer to iterable key\values.
func Convert[K, V any, k2, v2 any](next func() (K, V, bool), by func(K, V) (k2, v2)) ConvertKVIter[K, V, k2, v2, func(K, V) (k2, v2)] {
	return ConvertKVIter[K, V, k2, v2, func(K, V) (k2, v2)]{next: next, by: by}
}

// Filter creates an Iterator that checks elements by a filter and returns successful ones
func Filter[K, V any](next func() (K, V, bool), filter func(K, V) bool) FitKV[K, V] {
	return FitKV[K, V]{next: next, by: filter}
}

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[K comparable, E, V any](next func() (K, E, bool), valResolv func(bool, K, V, E) V) map[K]V {
	e := map[K]V{}
	for k, elem, ok := next(); ok; k, elem, ok = next() {
		exists, ok := e[k]
		e[k] = valResolv(ok, k, exists, elem)
	}
	return e
}
