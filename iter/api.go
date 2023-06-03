// Package iter provides generic constructors and helpers for iterators
package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert/as"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/slice"
)

// Of instantiates an iterator of predefined elements
func Of[T any](elements ...T) *slice.Iter[T] {
	return slice.NewIter(elements)
}

// Convert instantiates an iterator that converts elements with a converter and returns them
func Convert[From, To any, I c.Iterator[From]](elements I, converter func(From) To) loop.ConvertIter[From, To] {
	return loop.Convert(elements.Next, converter)
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[From, To any, I c.Iterator[From]](elements I, filter func(From) bool, converter func(From) To) loop.ConvertFiltIter[From, To] {
	return loop.FilterAndConvert(elements.Next, filter, converter)
}

// Flat instantiates an iterator that converts the collection elements into slices and then flattens them to one level
func Flat[From, To any, I c.Iterator[From]](elements I, flattener func(From) []To) *loop.FlatIter[From, To] {
	return loop.Flat(elements.Next, flattener)
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[From, To any, I c.Iterator[From]](elements I, filter func(From) bool, flattener func(From) []To) *loop.FlatFilterIter[From, To] {
	return loop.FilterAndFlat(elements.Next, filter, flattener)
}

// Filter instantiates an iterator that checks elements by a filter and returns successful ones
func Filter[T any, I c.Iterator[T]](elements I, filter func(T) bool) loop.FiltIter[T] {
	return loop.Filter(elements.Next, filter)
}

// NotNil instantiates an iterator that filters nullable elements
func NotNil[T any, I c.Iterator[*T]](elements I) loop.FiltIter[*T] {
	return Filter(elements, not.Nil[T])
}

// Reduce reduces elements to an one
func Reduce[T any, I c.Iterator[T]](elements I, by func(T, T) T) T {
	return loop.Reduce(elements.Next, by)
}

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable, I c.Iterator[T]](elements I, by func(T) K) stream.Iter[K, T, map[K][]T] {
	return stream.New(loop.KeyValue(elements.Next, by, as.Is[T]).Next, kvloop.Group[K, T])
}

// First returns the first element that satisfies requirements of the predicate 'filter'
func First[T any, I c.Iterator[T]](elements I, filter func(T) bool) (T, bool) {
	return loop.First(elements.Next, filter)
}

// KeyValue converts a c.Iterator to a kv.KVIterator using key and value extractors
func KeyValue[T, K, V any, IT c.Iterator[T]](elements IT, keyExtractor func(T) K, valExtractor func(T) V) loop.KeyValuer[T, K, V] {
	kv := loop.KeyValue(elements.Next, keyExtractor, valExtractor)
	return kv
}

func Start[T any, IT c.Iterator[T]](elements IT) (IT, T, bool) {
	element, ok := elements.Next()
	return elements, element, ok
}
