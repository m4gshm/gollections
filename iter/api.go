package iter

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
	sliceIter "github.com/m4gshm/gollections/slice/iter"
)

// Of instantiates Iterator of predefined elements
func Of[T any](elements ...T) c.Iterator[T] {
	return New(elements)
}

// New instantiates Iterator using a slice as the elements source
func New[T any](elements []T) c.Iterator[T] {
	return Wrap(elements)
}

// OfLoop creates an IteratorBreakable instance that loops over elements of a source
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfLoop[S, T any](source S, hasNext func(S) bool, next func(S) (T, error)) c.IteratorBreakable[T] {
	l := loop.NewIter(source, hasNext, next)
	return &l
}

// Wrap instantiates Iterator using a slice as the elements source
func Wrap[TS ~[]T, T any](elements TS) *sliceIter.SliceIter[T] {
	h := sliceIter.NewHead(elements)
	return &h
}

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[From, To any, IT c.Iterator[From]](elements IT, converter func(From) To) c.Iterator[To] {
	return loop.Convert(elements.Next, converter)
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[From, To any, IT c.Iterator[From]](elements IT, filter func(From) bool, converter func(From) To) c.Iterator[To] {
	return loop.FilterAndConvert(elements.Next, filter, converter)
}

// Flatt instantiates Iterator that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To any, IT c.Iterator[From]](elements IT, flatt func(From) []To) c.Iterator[To] {
	f := loop.Flatt(elements.Next, flatt)
	return &f
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, IT c.Iterator[From]](elements IT, filter func(From) bool, flatt func(From) []To) c.Iterator[To] {
	f := loop.FilterAndFlatt(elements.Next, filter, flatt)
	return &f
}

// Filter instantiates Iterator that checks elements by a filter and returns successful ones
func Filter[T any, IT c.Iterator[T]](elements IT, filter func(T) bool) c.Iterator[T] {
	f := loop.Filter(elements.Next, filter)
	return f
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, IT c.Iterator[*T]](elements IT) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

// Reduce reduces elements to an one
func Reduce[T any](elements c.Iterator[T], by func(T, T) T) T {
	return loop.Reduce(elements.Next, by)
}

// ReduceKV reduces key/value elements to an one
func ReduceKV[K, V any](elements c.KVIterator[K, V], by func(k1 K, v1 V, k2 K, v2 V) (K, V)) (K, V) {
	return kvloop.Reduce(elements.Next, by)
}

// ToSlice converts an Iterator to a slice
func ToSlice[T any](elements c.Iterator[T]) []T {
	return loop.ToSlice(elements.Next)
}

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable](elements c.Iterator[T], by func(T) K) c.MapStream[K, T, map[K][]T] {
	return kvloop.Stream(loop.NewKeyValuer(elements.Next, by, as.Is[T]).Next, kvloop.Group[K, T])
}

// ForEach applies the 'walker' function to elements of an Iterator
func ForEach[T any](elements c.Iterator[T], walker func(T)) {
	loop.ForEach(elements.Next, walker)
}

// ForEachFiltered applies the 'walker' function to elements that satisfy a predicate condition
func ForEachFiltered[T any](elements c.Iterator[T], walker func(T), filter func(T) bool) {
	loop.ForEachFiltered(elements.Next, walker, filter)
}

// Sum returns the sum of all elements
func Sum[T c.Summable](elements c.Iterator[T]) T {
	return loop.Reduce(elements.Next, op.Sum[T])
}

// First returns the first element that satisfies requirements of the predicate 'filter'
func First[T any, IT c.Iterator[T]](elements IT, filter func(T) bool) (T, bool) {
	return loop.First(elements.Next, filter)
}

// ToPairs converts a c.Iterator to a c.KVIterator using key and value extractors
func ToPairs[T, K, V any](elements c.Iterator[T], keyExtractor func(T) K, valExtractor func(T) V) c.KVIterator[K, V] {
	kv := loop.NewKeyValuer(elements.Next, keyExtractor, valExtractor)
	return &kv
}
