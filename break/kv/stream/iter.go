// Package stream provides a stream implementation and helper functions
package stream

import (
	"github.com/m4gshm/gollections/break/kv"
	"github.com/m4gshm/gollections/break/kv/loop"
	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakFilter "github.com/m4gshm/gollections/break/map_/filter"
	"github.com/m4gshm/gollections/map_/filter"
)

// New is the main stream constructor
func New[K comparable, V any, M map[K]V | map[K][]V](next func() (K, V, bool, error), collector MapCollector[K, V, M]) Iter[K, V, M] {
	return Iter[K, V, M]{next: next, collector: collector}
}

// Iter is the key/value Iterator based stream implementation.
type Iter[K comparable, V any, M map[K]V | map[K][]V] struct {
	next      func() (K, V, bool, error)
	collector MapCollector[K, V, M]
}

var (
	_ kv.Iterator[string, any]              = (*Iter[string, any, map[string]any])(nil)
	_ Stream[string, any, map[string]any]   = (*Iter[string, any, map[string]any])(nil)
	_ Stream[string, any, map[string][]any] = (*Iter[string, any, map[string][]any])(nil)

	_ kv.Iterator[string, any]              = Iter[string, any, map[string]any]{}
	_ Stream[string, any, map[string]any]   = Iter[string, any, map[string]any]{}
	_ Stream[string, any, map[string][]any] = Iter[string, any, map[string][]any]{}
)

var _ kv.IterFor[string, any, Iter[string, any, map[string]any]] = Iter[string, any, map[string]any]{}

// Next implements kv.KVIterator
func (i Iter[K, V, M]) Next() (K, V, bool, error) {
	return i.next()
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (i Iter[K, V, M]) FilterKey(predicate func(K) bool) Iter[K, V, M] {
	return New(loop.Filter(i.next, filter.Key[V](predicate)), i.collector)
}

// FiltKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (i Iter[K, V, M]) FiltKey(predicate func(K) (bool, error)) Iter[K, V, M] {
	return New(loop.Filt(i.next, breakFilter.Key[V](predicate)), i.collector)
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (i Iter[K, V, M]) FilterValue(predicate func(V) bool) Iter[K, V, M] {
	return New(loop.Filter(i.next, filter.Value[K](predicate)), i.collector)
}

// FiltValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (i Iter[K, V, M]) FiltValue(predicate func(V) (bool, error)) Iter[K, V, M] {
	return New(loop.Filt(i.next, breakFilter.Value[K](predicate)), i.collector)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (i Iter[K, V, M]) Filter(predicate func(K, V) bool) Iter[K, V, M] {
	return New(loop.Filter(i.next, predicate), i.collector)
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (i Iter[K, V, M]) Filt(predicate func(K, V) (bool, error)) Iter[K, V, M] {
	return New(loop.Filt(i.next, predicate), i.collector)
}

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (i Iter[K, V, M]) Track(tracker func(K, V) error) error {
	return breakLoop.Track(i.next, tracker)
}

// Reduce reduces the key/value pairs into an one pair using the 'merge' function
func (i Iter[K, V, M]) Reduce(by func(K, K, V, V) (K, V, error)) (K, V, error) {
	return loop.Reducee(i.next, by)
}

// HasAny finds the first key/value pari that satisfies the 'predicate' function condition and returns true if successful
func (i Iter[K, V, M]) HasAny(predicate func(K, V) (bool, error)) (bool, error) {
	next := i.next
	return loop.HasAnyy(next, predicate)
}

// Iter creates an iterator and returns as interface
func (i Iter[K, V, M]) Loop() func() (K, V, bool, error) {
	return i.Next
}

// Map collects the key/value pairs to a map
func (i Iter[K, V, M]) Map() (M, error) {
	return i.collector(i.next)
}

// Start is used with for loop construct like 'for i, k, v, ok, err := i.Start(); ok || err != nil ; k, v, ok, err = i.Next() { if err != nil { return err }}'
func (i Iter[K, V, M]) Start() (Iter[K, V, M], K, V, bool, error) {
	k, v, ok, err := i.next()
	return i, k, v, ok, err
}

// MapCollector is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result
type MapCollector[K comparable, V any, M map[K]V | map[K][]V] func(next func() (K, V, bool, error)) (M, error)
