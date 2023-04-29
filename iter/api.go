// Package iter provides generic constructors and helpers for iterators
package iter

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// Of instantiates an iterator of predefined elements
func Of[T any](elements ...T) *slice.Iter[T] {
	return slice.NewIter(elements)
}

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[From, To any, I c.Iterator[From]](elements I, converter func(From) To) loop.ConvertIter[From, To] {
	return loop.Convert(elements.Next, converter)
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[From, To any, I c.Iterator[From]](elements I, filter func(From) bool, converter func(From) To) loop.ConvertFitIter[From, To] {
	return loop.FilterAndConvert(elements.Next, filter, converter)
}

// Flatt instantiates Iterator that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To any, I c.Iterator[From]](elements I, flatt func(From) []To) *loop.FlatIter[From, To] {
	f := loop.Flatt(elements.Next, flatt)
	return &f
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, I c.Iterator[From]](elements I, filter func(From) bool, flatt func(From) []To) *loop.FlattenFitIter[From, To] {
	f := loop.FilterAndFlatt(elements.Next, filter, flatt)
	return &f
}

// Filter instantiates Iterator that checks elements by a filter and returns successful ones
func Filter[T any, I c.Iterator[T]](elements I, filter func(T) bool) loop.FitIter[T] {
	f := loop.Filter(elements.Next, filter)
	return f
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, I c.Iterator[*T]](elements I) loop.FitIter[*T] {
	return Filter(elements, check.NotNil[T])
}

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable, I c.Iterator[T]](elements I, by func(T) K) stream.Iter[K, T, map[K][]T] {
	return stream.New(loop.NewKeyValuer(elements.Next, by, as.Is[T]).Next, kvloop.Group[K, T])
}

// First returns the first element that satisfies requirements of the predicate 'filter'
func First[T any, I c.Iterator[T]](elements I, filter func(T) bool) (T, bool) {
	return loop.First(elements.Next, filter)
}

// ToKV converts a c.Iterator to a c.KVIterator using key and value extractors
func ToKV[T, K, V any, IT c.Iterator[T]](elements IT, keyExtractor func(T) K, valExtractor func(T) V) loop.KeyValuer[T, K, V] {
	kv := loop.NewKeyValuer(elements.Next, keyExtractor, valExtractor)
	return kv
}
