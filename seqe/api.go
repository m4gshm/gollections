// Package seqe provides convering, filtering, and reducing operations for the [seq.SeqE] interface.
package seqe

import (
	"github.com/m4gshm/gollections/internal/seq"
	"github.com/m4gshm/gollections/internal/seqe"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/always"
)

// Seq is an alias of an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[T any] = seq.Seq[T]

// SeqE is a specific iterator form that allows to retrieve a value with an error as second parameter of the iterator.
// It is used as a result of applying functions like seq.Conv, which may throw an error during iteration.
type SeqE[T any] = seq.SeqE[T]

// Union combines several sequences into one.
func Union[S ~SeqE[T], T any](seq ...S) seq.SeqE[T] {
	return seqe.Union(seq...)
}

// OfNextGet builds an iterator by iterating elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfNextGet[T any](hasNext func() bool, getNext func() (T, error)) seq.SeqE[T] {
	return seqe.OfNextGet(hasNext, getNext)
}

// OfNext builds an iterator by iterating elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The pushNext copy the element to the next pointer.
func OfNext[T any](hasNext func() bool, pushNext func(*T) error) seq.SeqE[T] {
	return OfNextGet(hasNext, func() (o T, err error) { return o, pushNext(&o) })
}

// OfSourceNextGet builds an iterator by iterating elements of the source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfSourceNextGet[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) seq.SeqE[T] {
	return OfNextGet(func() bool { return hasNext(source) }, func() (T, error) { return getNext(source) })
}

// OfSourceNext builds an iterator by iterating elements of the source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The pushNext copy the element to the next pointer.
func OfSourceNext[S, T any](source S, hasNext func(S) bool, pushNext func(S, *T) error) seq.SeqE[T] {
	return OfNext(func() bool { return hasNext(source) }, func(next *T) error { return pushNext(source, next) })
}

// OfIndexed builds a SeqE iterator by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](amount int, getAt func(int) (T, error)) SeqE[T] {
	return seqe.OfIndexed(amount, getAt)
}

// Top returns a sequence of top n elements.
func Top[S ~SeqE[T], T any](n int, seq S) seq.SeqE[T] {
	return seqe.Top(n, seq)
}

// Skip returns a sequence without first n elements.
func Skip[S ~SeqE[T], T any](n int, seq S) seq.SeqE[T] {
	return seqe.Skip(n, seq)
}

// Head returns the first element.
func Head[S ~SeqE[T], T any](seq S) (v T, ok bool, err error) {
	return First(seq, always.True)
}

// First returns the first element that satisfies the condition of the 'predicate' function.
func First[S ~SeqE[T], T any](seq S, predicate func(T) bool) (v T, ok bool, err error) {
	return seqe.First(seq, predicate)
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function.
func Firstt[S ~SeqE[T], T any](seq S, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	return seqe.Firstt(seq, predicate)
}

// Slice collects the elements of the 'seq' sequence into a new slice.
func Slice[S ~SeqE[T], T any](seq S) ([]T, error) {
	return SliceCap(seq, 0)
}

// SliceCap collects the elements of the 'seq' sequence into a new slice with predefined capacity.
func SliceCap[S ~SeqE[T], T any](seq S, capacity int) (out []T, e error) {
	return seqe.SliceCap(seq, capacity)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func Append[S ~SeqE[T], T any, TS ~[]T](seq S, out TS) (TS, error) {
	return seqe.Append(seq, out)
}

// Reduce reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reduce[S ~SeqE[T], T any](seq S, merge func(T, T) T) (T, error) {
	result, _, err := ReduceOK(seq, merge)
	return result, err
}

// ReduceOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceOK[S ~SeqE[T], T any](seq S, merge func(T, T) T) (result T, ok bool, err error) {
	return seqe.ReduceOK(seq, merge)
}

// Reducee reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reducee[S ~SeqE[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

// ReduceeOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceeOK[S ~SeqE[T], T any](seq S, merge func(T, T) (T, error)) (result T, ok bool, err error) {
	return seqe.ReduceeOK(seq, merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accum[T any, S ~SeqE[T]](first T, seq S, merge func(T, T) T) (accumulator T, err error) {
	return seqe.Accum(first, seq, merge)
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accumm[T any, S ~SeqE[T]](first T, seq S, merge func(T, T) (T, error)) (accumulator T, err error) {
	return seqe.Accumm(first, seq, merge)
}

// Sum returns the sum of all elements.
func Sum[S ~SeqE[T], T op.Summable](seq S) (out T, err error) {
	return Accum(out, seq, op.Sum[T])
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful.
func HasAny[S ~SeqE[T], T any](seq S, predicate func(T) bool) (bool, error) {
	_, ok, err := First(seq, predicate)
	return ok, err
}

// Contains finds the first element that equal to the example and returns true.
func Contains[S ~SeqE[T], T comparable](seq S, example T) (contains bool, err error) {
	return seqe.Contains(seq, example)
}

// Conv creates an iterator that applies the 'converter' function to each iterable element and returns value-error pairs.
// The error should be checked at every iteration step, like:
//
//	var integers iter.Seq2[int, error]
//	...
//	for s, err := range seqe.Conv(integers,  strconv.Itoa) {
//	    if err != nil {
//	        break
//	    }
//	    ...
//	}
func Conv[S ~SeqE[From], From, To any](seq S, converter func(From) (To, error)) seq.SeqE[To] {
	return seqe.Conv(seq, converter)
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func Convert[S ~SeqE[From], From, To any](seq S, converter func(From) To) seq.SeqE[To] {
	return seqe.Convert(seq, converter)
}

// ConvertOK creates an iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from the loop.
func ConvertOK[S ~SeqE[From], From, To any](seq S, converter func(from From) (To, bool)) seq.SeqE[To] {
	return seqe.ConvertOK(seq, converter)
}

// ConvOK creates a iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from iteration.
// It may also return an error to abort the iteration.
func ConvOK[S ~SeqE[From], From, To any](seq S, converter func(from From) (To, bool, error)) seq.SeqE[To] {
	return seqe.ConvOK(seq, converter)
}

// Flat is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var arrays seq.SeqE[[]int]
//	...
//	for e, err := range seqe.Flat(arrays, as.Is) {
//		if err != nil {
//			panic(err)
//		}
//	}
func Flat[S ~SeqE[From], STo ~[]To, From any, To any](seq S, flattener func(From) STo) seq.SeqE[To] {
	return seqe.Flat(seq, flattener)
}

// FlatSeq is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var arrays seq.SeqE[[]int]
//	...
//	for e, err := range seqe.FlatSeq(arrays, slices.Values) {
//		if err != nil {
//			panic(err)
//		}
//	}
func FlatSeq[S ~SeqE[From], STo ~Seq[To], From any, To any](seq S, flattener func(From) STo) seq.SeqE[To] {
	return seqe.FlatSeq(seq, flattener)
}

// Flatt is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var (
//		input     iter.Seq[[]string]
//		flattener func([]string) ([]int, error)
//		out       seq.SeqE[int]
//
//	)
//
//	flattener = convertEveryBy(strconv.Atoi)
//	out = seq.Flatt(input, flattener)
//	for i, err := range out {
//		if err != nil {
//			panic(err)
//		}
//		...
//	}
func Flatt[S ~SeqE[From], STo ~[]To, From any, To any](seq S, flattener func(From) (STo, error)) seq.SeqE[To] {
	return seqe.Flatt(seq, flattener)
}

// FlattSeq is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var (
//		input     iter.Seq[[]string]
//		flattener func([]string) seq.SeqE[int]
//		out       seq.SeqE[int]
//
//	)
//
//	flattener = convertEveryBy(strconv.Atoi)
//	out = seq.Flatt(input, flattener)
//	for i, err := range out {
//		if err != nil {
//			panic(err)
//		}
//		...
//	}
func FlattSeq[S ~SeqE[From], STo ~SeqE[To], From any, To any](seq S, flattener func(From) STo) seq.SeqE[To] {
	return seqe.FlattSeq(seq, flattener)
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~SeqE[T], T any](seq S, filter func(T) bool) seq.SeqE[T] {
	return seqe.Filter(seq, filter)
}

// Filt creates an erroreable iterator that iterates only those elements for which the 'filter' function returns true.
func Filt[S ~SeqE[T], T any](seq S, filter func(T) (bool, error)) seq.SeqE[T] {
	return seqe.Filt(seq, filter)
}

// Group collects the seq elements into a new map.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to a value.
func Group[S ~SeqE[T], T any, K comparable, V any](seq S, keyExtractor func(T) K, valExtractor func(T) V) (map[K][]V, error) {
	return seqe.Group(seq, keyExtractor, valExtractor)
}

// ForEach applies the 'consumer' function to the seq elements
func ForEach[T any](seq SeqE[T], consumer func(T)) error {
	return seqe.ForEach(seq, consumer)
}
