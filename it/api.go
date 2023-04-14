package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/ptr"
)

// Of instantiates Iterator of predefined elements
func Of[T any](elements ...T) c.Iterator[T] {
	return ptr.Of(it.NewHead(elements))
}

// OfLoop creates an IteratorBreakable instance that loops over elements of a source
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfLoop[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) c.IteratorBreakable[T] {
	return ptr.Of(it.NewLoop(source, hasNext, getNext))
}

// Wrap instantiates Iterator using slice as the elements source
func Wrap[TS ~[]T, T any](elements TS) c.Iterator[T] {
	return ptr.Of(it.NewHead(elements))
}

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[From, To any, IT c.Iterator[From]](elements IT, converter func(From) To) c.Iterator[To] {
	return it.Convert(elements, converter)
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[From, To any, IT c.Iterator[From]](elements IT, filter func(From) bool, converter func(From) To) c.Iterator[To] {
	return it.FilterAndConvert(elements, filter, converter)
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, IT c.Iterator[From]](elements IT, flatt func(From) []To) c.Iterator[To] {
	return ptr.Of(it.Flatt(elements, flatt))
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, IT c.Iterator[From]](elements IT, filter func(From) bool, flatt func(From) []To) c.Iterator[To] {
	return ptr.Of(it.FilterAndFlatt(elements, filter, flatt))
}

// Filter instantiates Iterator that checks elements by a filter and returns successful ones
func Filter[T any, IT c.Iterator[T]](elements IT, filter func(T) bool) c.Iterator[T] {
	return it.Filter(elements, filter)
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, IT c.Iterator[*T]](elements IT) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

// Reduce reduces elements to an one
func Reduce[T any, IT c.Iterator[T]](elements IT, by func(T, T) T) T {
	return it.Reduce(elements, by)
}

// ReduceKV reduces key/value elements to an one
func ReduceKV[K, V any, IT c.KVIterator[K, V]](elements IT, by c.Quaternary[K, V]) (K, V) {
	return it.ReduceKV(elements, by)
}

// ToSlice converts an Iterator to a slice
func ToSlice[T any](elements c.Iterator[T]) []T {
	return it.ToSlice[T](elements)
}

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable](elements c.Iterator[T], by func(T) K) c.MapPipe[K, T, map[K][]T] {
	return it.Group(elements, by)
}

// ForEach applies a walker to elements of an Iterator
func ForEach[T any, IT c.Iterator[T]](elements IT, walker func(T)) {
	it.ForEach(elements, walker)
}

// ForEachFiltered applies a walker to elements that satisfy a predicate condition
func ForEachFiltered[T any](elements c.Iterator[T], walker func(T), filter func(T) bool) {
	it.ForEachFiltered(elements, walker, filter)
}

// Sum returns the sum of all elements
func Sum[T c.Summable, IT c.Iterator[T]](elements IT) T {
	return it.Reduce(elements, op.Sum[T])
}

// First returns the first element that satisfies requirements of the predicate 'filter'
func First[T any, IT c.Iterator[T]](elements IT, filter func(T) bool) (T, bool) {
	return it.First(elements, filter)
}

// ToPairs converts a c.Iterator to a c.KVIterator using key and value extractors
func ToPairs[T, K, V any](elements c.Iterator[T], keyExtractor func(T) K, valExtractor func(T) V) c.KVIterator[K, V] {
	return ptr.Of(it.NewKeyValuer(elements, keyExtractor, valExtractor))
}
