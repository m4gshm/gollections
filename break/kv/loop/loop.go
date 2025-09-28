package loop

// Loop is a function that returns the next key\value or ok==false if there are no more elements.
//
// Deprecated: replaced by [github.com/m4gshm/gollections/seq.Seq2]
type Loop[K, V any] func() (key K, value V, ok bool, err error)

// Track applies the 'consumer' function to position/element pairs retrieved by the 'next' function until the consumer returns the c.Break to stop.
func (next Loop[K, V]) Track(consumer func(K, V) error) error {
	return Track(next, consumer)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func (next Loop[K, V]) First(predicate func(K, V) bool) (K, V, bool, error) {
	return First(next, predicate)
}

// Reduce reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero values of 'K', 'V' types are returned.
func (next Loop[K, V]) Reduce(merge func(K, K, V, V) (K, V)) (K, V, error) {
	return Reduce(next, merge)
}

// ReduceOK reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func (next Loop[K, V]) ReduceOK(merge func(K, K, V, V) (K, V)) (K, V, bool, error) {
	return ReduceOK(next, merge)
}

// Reducee reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero values of 'K', 'V' types are returned.
func (next Loop[K, V]) Reducee(merge func(K, K, V, V) (K, V, error)) (K, V, error) {
	return Reducee(next, merge)
}

// ReduceeOK reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func (next Loop[K, V]) ReduceeOK(merge func(K, K, V, V) (K, V, error)) (K, V, bool, error) {
	return ReduceeOK(next, merge)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (next Loop[K, V]) HasAny(predicate func(K, V) bool) (bool, error) {
	return HasAny(next, predicate)
}

// Filt creates a seq that checks elements by the 'filter' function and returns successful ones.
func (next Loop[K, V]) Filt(filter func(K, V) (bool, error)) Loop[K, V] {
	return Filt(next, filter)
}

// Filter creates a seq that checks elements by the 'filter' function and returns successful ones.
func (next Loop[K, V]) Filter(filter func(K, V) bool) Loop[K, V] {
	return Filter(next, filter)
}

// Crank rertieves next key\value elements from the 'next' function, returns the function, element, successfully flag.
func (next Loop[K, V]) Crank() (Loop[K, V], K, V, bool, error) {
	return Crank(next)
}
