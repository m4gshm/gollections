// Package collection consists of common operations of Iterable based collections
package collection

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/comparer"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/seq"
)

// Head returns the first element.
func Head[IT c.Range[T], T any](collection IT) (T, bool) {
	return seq.Head(collection.All)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func Convert[IT c.Range[From], From, To any](collection IT, converter func(From) To) seq.Seq[To] {
	return seq.Convert(collection.All, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func Conv[IT c.Range[From], From, To any](collection IT, converter func(From) (To, error)) seq.SeqE[To] {
	return seq.Conv(collection.All, converter)
}

// FilterAndConvert returns a seq that filters source elements and converts them
func FilterAndConvert[IT c.Range[From], From, To any](collection IT, filter func(From) bool, converter func(From) To) seq.Seq[To] {
	return seq.Convert(seq.Filter(collection.All, filter), converter)
}

// Flat returns a seq that converts the collection elements into slices and then flattens them to one level
func Flat[IT c.Range[From], From, To any](collection IT, by func(From) []To) seq.Seq[To] {
	return seq.Flat(collection.All, by)
}

// Flatt returns a errorable seq that converts the collection elements into slices and then flattens them to one level
func Flatt[IT c.Range[From], From, To any](collection IT, flattener func(From) ([]To, error)) seq.SeqE[To] {
	return seq.Flatt(collection.All, flattener)
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[IT c.Range[From], From, To any](collection IT, filter func(From) bool, flattener func(From) []To) seq.Seq[To] {
	return seq.Flat(seq.Filter(collection.All, filter), flattener)
}

// Filter instantiates a seq that checks elements by the 'filter' function and returns successful ones.
func Filter[IT c.Range[T], T any](collection IT, filter func(T) bool) seq.Seq[T] {
	return seq.Filter(collection.All, filter)
}

// Filter instantiates a seq that checks elements by the 'filter' function and returns successful ones.
func Filt[IT c.Range[T], T any](collection IT, filter func(T) (bool, error)) seq.SeqE[T] {
	return seq.Filt(collection.All, filter)
}

// NotNil instantiates a seq that filters nullable elements
func NotNil[IT c.Range[*T], T any](collection IT) seq.Seq[*T] {
	return Filter(collection, not.Nil[T])
}

// // PtrVal creates a seq that transform pointers to the values referenced referenced by those pointers.
// // Nil pointers are transformet to zero values.
// func PtrVal[IT c.Range[*T], T any](collection IT) loop.Loop[T] {
// 	return loop.PtrVal(collection.All)
// }

// // NoNilPtrVal creates a seq that transform only not nil pointers to the values referenced referenced by those pointers.
// // Nil pointers are ignored.
// func NoNilPtrVal[T any, IT Iterable[*T]](collection IT) loop.Loop[T] {
// 	return loop.NoNilPtrVal(collection.All)
// }

// // NilSafe creates a seq that filters not nil elements, converts that ones, filters not nils after converting and returns them
// func NilSafe[From, To any, IT Iterable[*From]](collection IT, converter func(*From) *To) loop.Loop[*To] {
// 	h := collection.All
// 	return loopconvert.NilSafe(h, converter)
// }

// KeyValue transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValue[IT c.Range[T], T any, K comparable, V any](collection IT, keyExtractor func(T) K, valExtractor func(T) V) seq.Seq2[K, V] {
	return seq.ToKV(collection.All, keyExtractor, valExtractor)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[IT c.Range[T], T any](collection IT, predicate func(T) bool) (v T, ok bool) {
	return seq.First(collection.All, predicate)
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function
func Firstt[IT c.Range[T], T any](collection IT, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	return seq.Firstt(collection.All, predicate)
}

// Sort sorts the specified sortable collection that contains orderable elements
func Sort[SC any, Cmp ~func(T, T) int, C interface {
	Sort(Cmp) SC
}, T any, O constraints.Ordered](collection C, order func(T) O) SC {
	return collection.Sort(comparer.Of(order))
}

// Reduce reduces the 'collection' elements into an one using the 'merge' function.
// If the 'collection' is empty, the zero value of 'T' type is returned.
func Reduce[IT c.Range[T], T any](collection IT, merge func(T, T) T) T {
	return seq.Reduce(collection.All, merge)
}

// Reducee reduces the 'collection' elements into an one using the 'merge' function.
// Returns ok==false if the 'collection' is empty.
func Reducee[IT c.Range[T], T any](collection IT, merge func(T, T) (T, error)) (T, error) {
	return seq.Reducee(collection.All, merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'collection'.
func Accum[IT c.Range[T], T any](first T, collection IT, merge func(T, T) T) T {
	return seq.Accum(first, collection.All, merge)
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'collection'.
func Accumm[IT c.Range[T], T any](first T, collection IT, merge func(T, T) (T, error)) (T, error) {
	return seq.Accumm(first, collection.All, merge)
}

// IsEmpty returns true if the collection is empty
func IsEmpty[C interface{ Len() int }](collection C) bool {
	return collection.Len() == 0
}
