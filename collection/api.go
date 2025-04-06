// Package collection consists of common operations of Iterable based collections
package collection

import (
	"golang.org/x/exp/constraints"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/comparer"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
	loopconvert "github.com/m4gshm/gollections/loop/convert"
	"github.com/m4gshm/gollections/op/check/not"
)

// Convert returns a loop that applies the 'converter' function to the collection elements
func Convert[From, To any, IT Iterable[From]](collection IT, converter func(From) To) loop.Loop[To] {
	b := collection.Loop()
	return loop.Convert(b, converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func Conv[From, To any, IT Iterable[From]](collection IT, converter func(From) (To, error)) breakLoop.Loop[To] {
	b := collection.Loop()
	return loop.Conv(b, converter)
}

// FilterAndConvert returns a loop that filters source elements and converts them
func FilterAndConvert[From, To any, IT Iterable[From]](collection IT, filter func(From) bool, converter func(From) To) loop.Loop[To] {
	b := collection.Loop()
	f := loop.FilterAndConvert(b, filter, converter)
	return f
}

// Flat returns a loop that converts the collection elements into slices and then flattens them to one level
func Flat[From, To any, IT Iterable[From]](collection IT, by func(From) []To) loop.Loop[To] {
	b := collection.Loop()
	f := loop.Flat(b, by)
	return f
}

// Flatt returns a breakable loop that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable, IT Iterable[From]](collection IT, flattener func(From) ([]To, error)) breakLoop.Loop[To] {
	return loop.Flatt(collection.Loop(), flattener)
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[From, To any, IT Iterable[From]](collection IT, filter func(From) bool, flattener func(From) []To) loop.Loop[To] {
	b := collection.Loop()
	f := loop.FilterAndFlat(b, filter, flattener)
	return f
}

// Filter instantiates a loop that checks elements by the 'filter' function and returns successful ones
func Filter[T any, IT Iterable[T]](collection IT, filter func(T) bool) loop.Loop[T] {
	b := collection.Loop()
	f := loop.Filter(b, filter)
	return f
}

// NotNil instantiates a loop that filters nullable elements
func NotNil[T any, IT Iterable[*T]](collection IT) loop.Loop[*T] {
	return Filter(collection, not.Nil[T])
}

// PtrVal creates a loop that transform pointers to the values referenced referenced by those pointers.
// Nil pointers are transformet to zero values.
func PtrVal[T any, IT Iterable[*T]](collection IT) loop.Loop[T] {
	return loop.PtrVal(collection.Loop())
}

// NoNilPtrVal creates a loop that transform only not nil pointers to the values referenced referenced by those pointers.
// Nil pointers are ignored.
func NoNilPtrVal[T any, IT Iterable[*T]](collection IT) loop.Loop[T] {
	return loop.NoNilPtrVal(collection.Loop())
}

// NilSafe creates a loop that filters not nil elements, converts that ones, filters not nils after converting and returns them
func NilSafe[From, To any, IT Iterable[*From]](collection IT, converter func(*From) *To) loop.Loop[*To] {
	h := collection.Loop()
	return loopconvert.NilSafe(h, converter)
}

// KeyValue transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValue[T any, K comparable, V any, IT Iterable[T]](collection IT, keyExtractor func(T) K, valExtractor func(T) V) kvloop.Loop[K, V] {
	h := collection.Loop()
	return loop.KeyValue(h, keyExtractor, valExtractor)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any, IT Iterable[T]](collection IT, predicate func(T) bool) (v T, ok bool) {
	i := collection.Loop()
	return loop.First(i, predicate)
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function
func Firstt[T any, IT Iterable[T]](collection IT, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	return loop.Firstt(collection.Loop(), predicate)
}

// Sort sorts the specified sortable collection that contains orderable elements
func Sort[SC any, Cmp ~func(T, T) int, C interface {
	Sort(Cmp) SC
}, T any, O constraints.Ordered](collection C, order func(T) O) SC {
	return collection.Sort(comparer.Of(order))
}

// Reduce reduces the 'collection' elements into an one using the 'merge' function.
// If the 'collection' is empty, the zero value of 'T' type is returned.
func Reduce[T any, IT Iterable[T]](collection IT, merge func(T, T) T) T {
	return loop.Reduce(collection.Loop(), merge)
}

// Reducee reduces the 'collection' elements into an one using the 'merge' function.
// Returns ok==false if the 'collection' is empty.
func Reducee[T any, IT Iterable[T]](collection IT, merge func(T, T) (T, error)) (T, error) {
	return loop.Reducee(collection.Loop(), merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'collection'.
func Accum[T any, IT Iterable[T]](first T, collection IT, merge func(T, T) T) T {
	return loop.Accum(first, collection.Loop(), merge)
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'collection'.
func Accumm[T any, IT Iterable[T]](first T, collection IT, merge func(T, T) (T, error)) (T, error) {
	return loop.Accumm(first, collection.Loop(), merge)
}

func IsEmpty[C interface{ Len() int }](collection C) bool {
	return collection.Len() == 0
}
