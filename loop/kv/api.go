package kv

// Group collects sets of values grouped by keys obtained by passing a key/value iterator
func Group[K comparable, V any](next func() (K, V, bool)) map[K][]V {
	e := map[K][]V{}
	for k, v, ok := next(); ok; k, v, ok = next() {
		e[k] = append(e[k], v)
	}
	return e
}
