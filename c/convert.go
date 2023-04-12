package c

import "github.com/m4gshm/gollections/predicate"

// Converter convert From -> To.
type Converter[From, To any] func(From) To

// BiConverter convert pairs of From -> To.
type BiConverter[From1, From2, To1, To2 any] func(From1, From2) (To1, To2)

// Flatter extracts slice of To.
type Flatter[From, To any] Converter[From, []To]

// FitKey adapts a key appliable predicate to a key\value one
func FitKey[K, V any](filter predicate.Predicate[K]) predicate.BiPredicate[K, V] {
	return func(key K, val V) bool { return filter(key) }
}

// FitValue adapts a value appliable predicate to a key\value one
func FitValue[K, V any](filter predicate.Predicate[V]) predicate.BiPredicate[K, V] {
	return func(key K, val V) bool { return filter(val) }
}
