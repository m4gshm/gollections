// Package loop provides helpers for loop operation over key/value pairs and iterator implementations
package loop

// Group collects sets of values grouped by keys obtained by passing a key/value iterator
func Group[K comparable, V any](next func() (K, V, bool)) map[K][]V {
	e := map[K][]V{}
	for k, v, ok := next(); ok; k, v, ok = next() {
		e[k] = append(e[k], v)
	}
	return e
}

// Reduce reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function
func Reduce[K, V any](next func() (K, V, bool), merge func(K, V, K, V) (K, V)) (rk K, rv V) {
	if k, v, ok := next(); ok {
		rk, rv = k, v
	} else {
		return rk, rv
	}
	for k, v, ok := next(); ok; k, v, ok = next() {
		rk, rv = merge(rk, rv, k, v)
	}
	return rk, rv
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func HasAny[K, V any](next func() (K, V, bool), predicate func(K, V) bool) bool {
	for k, v, ok := next(); ok; k, v, ok = next() {
		if predicate(k, v) {
			return true
		}
	}
	return false
}

// Convert creates an Iterator that applies a transformer to iterable key\values.
func Convert[K, V any, k2, v2 any](next func() (K, V, bool), by func(K, V) (k2, v2)) ConvertKVIter[K, V, k2, v2, func(K, V) (k2, v2)] {
	return ConvertKVIter[K, V, k2, v2, func(K, V) (k2, v2)]{next: next, by: by}
}

// Filter creates an Iterator that checks elements by a filter and returns successful ones
func Filter[K, V any](next func() (K, V, bool), filter func(K, V) bool) FitKV[K, V] {
	return FitKV[K, V]{next: next, by: filter}
}

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[K comparable, V any](next func() (K, V, bool), resolver func(bool, K, V, V) V) map[K]V {
	m := map[K]V{}
	for k, v, ok := next(); ok; k, v, ok = next() {
		exists, ok := m[k]
		m[k] = resolver(ok, k, exists, v)
	}
	return m
}
