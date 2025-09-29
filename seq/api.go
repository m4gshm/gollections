// Package seq extends [iter.Seq] API with convering, filtering, and reducing functionality.
package seq

import (
	s "github.com/m4gshm/gollections/internal/seq"
	s2 "github.com/m4gshm/gollections/internal/seq2"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/always"
	"golang.org/x/exp/constraints"
)

// Seq is an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[T any] s.Seq[T]
type seq[T any] = s.Seq[T]

// SeqE is a specific iterator form that allows to retrieve a value with an error as second parameter of the iterator.
// It is used as a result of applying functions like seq.Conv, which may throw an error during iteration.
// At each iteration step, it is necessary to check for the occurrence of an error.
//
//	for e, err := range seqence {
//	    if err != nil {
//	        break
//	    }
//	    ...
//	}
type SeqE[T any] s.SeqE[T]
type seqE[T any] = s.SeqE[T]

// Seq2 is an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] s.Seq2[K, V]
type seq2[K, V any] = s.Seq2[K, V]

// Of creates an iterator over the elements.
func Of[T any](elements ...T) Seq[T] {
	return s.Of(elements...)
}

// Union combines several sequences into one.
func Union[S ~seq[T], T any](seq ...S) Seq[T] {
	return s.Union(seq...)
}

// OfNextGet builds an iterator by iterating elements of a source.
// The hasNext specifies a filter that tests existing of a next element in the source.
// The getNext extracts the element.
func OfNextGet[T any](hasNext func() bool, getNext func() T) Seq[T] {
	return s.OfNextGet(hasNext, getNext)
}

// OfNext builds an iterator by iterating elements of a source.
// The hasNext specifies a filter that tests existing of a next element in the source.
// The pushNext copy the element to the next pointer.
func OfNext[T any](hasNext func() bool, pushNext func(*T)) Seq[T] {
	return OfNextGet(hasNext, func() (o T) { pushNext(&o); return o })
}

// OfSourceNextGet builds an iterator by iterating elements of the source.
// The hasNext specifies a filter that tests existing of a next element in the source.
// The getNext extracts the element.
func OfSourceNextGet[S, T any](source S, hasNext func(S) bool, getNext func(S) T) Seq[T] {
	return OfNextGet(func() bool { return hasNext(source) }, func() T { return getNext(source) })
}

// OfSourceNext builds an iterator by iterating elements of the source.
// The hasNext specifies a filter that tests existing of a next element in the source.
// The pushNext copy the element to the next pointer.
func OfSourceNext[S, T any](source S, hasNext func(S) bool, pushNext func(S, *T)) Seq[T] {
	return OfNext(func() bool { return hasNext(source) }, func(next *T) { pushNext(source, next) })
}

// OfIndexed builds a Seq iterator by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](amount int, getAt func(int) T) Seq[T] {
	return s.OfIndexed(amount, getAt)
}

// Series makes a sequence by applying the 'next' function to the previous step generated value.
func Series[T any](first T, next func(T) (T, bool)) Seq[T] {
	return s.Series(first, next)
}

// RangeClosed creates a sequence that generates integers in the range defined by from and to inclusive
func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) Seq[T] {
	return s.RangeClosed(from, toInclusive)
}

// Range creates a sequence that generates integers in the range defined by from and to exclusive
func Range[T constraints.Integer | rune](from T, toExclusive T) Seq[T] {
	return s.Range(from, toExclusive)
}

// ToSeq2 converts an iterator of single elements to an iterator of key/value pairs by applying the 'converter' function to each iterable element.
func ToSeq2[S ~seq[T], T, K, V any](seq S, converter func(T) (K, V)) Seq2[K, V] {
	return s.ToSeq2(seq, converter)
}

// Top returns a sequence of top n elements.
func Top[S ~seq[T], T any](n int, seq S) Seq[T] {
	return s.Top(n, seq)
}

// Skip returns a sequence without first n elements.
func Skip[S ~seq[T], T any](n int, seq S) Seq[T] {
	return s.Skip(n, seq)
}

// While cuts tail elements of the seq that don't match the filter.
func While[S ~seq[T], T any](seq S, filter func(T) bool) Seq[T] {
	return s.While(seq, filter)
}

// SkipWhile returns a sequence without first elements of the seq that dont'math the filter.
func SkipWhile[S ~seq[T], T any](seq S, filter func(T) bool) Seq[T] {
	return s.SkipWhile(seq, filter)
}

// Head returns the first element.
func Head[S ~seq[T], T any](seq S) (v T, ok bool) {
	return First(seq, always.True)
}

// First returns the first element that satisfies the condition of the 'filter' function.
func First[S ~seq[T], T any](seq S, filter func(T) bool) (v T, ok bool) {
	return s.First(seq, filter)
}

// Firstt returns the first element that satisfies the condition of the 'filter' function.
func Firstt[S ~seq[T], T any](seq S, filter func(T) (bool, error)) (v T, ok bool, err error) {
	return s.Firstt(seq, filter)
}

// Slice collects the elements of the 'seq' sequence into a new slice.
func Slice[S ~seq[T], T any](seq S) []T {
	return SliceCap(seq, 0)
}

// SliceCap collects the elements of the 'seq' sequence into a new slice with predefined capacity.
func SliceCap[S ~seq[T], T any](seq S, capacity int) (out []T) {
	return s.SliceCap(seq, capacity)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func Append[S ~seq[T], TS ~[]T, T any](seq S, out TS) TS {
	return s.Append(seq, out)
}

// Reduce reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reduce[S ~seq[T], T any](seq S, merge func(T, T) T) T {
	result, _ := ReduceOK(seq, merge)
	return result
}

// ReduceOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceOK[S ~seq[T], T any](seq S, merge func(T, T) T) (result T, ok bool) {
	return s.ReduceOK(seq, merge)
}

// Reducee reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reducee[S ~seq[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

// ReduceeOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceeOK[S ~seq[T], T any](seq S, merge func(T, T) (T, error)) (result T, ok bool, err error) {
	return s.ReduceeOK(seq, merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accum[T any, S ~seq[T]](first T, seq S, merge func(T, T) T) T {
	return s.Accum(first, seq, merge)
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accumm[T any, S ~seq[T]](first T, seq S, merge func(T, T) (T, error)) (accumulator T, err error) {
	return s.Accumm(first, seq, merge)
}

// Sum returns the sum of all elements.
func Sum[S ~seq[T], T op.Summable](seq S) (out T) {
	return s.Sum(seq)
}

// HasAny finds the first element that satisfies the 'filter' function condition and returns true if successful.
func HasAny[S ~seq[T], T any](seq S, filter func(T) bool) bool {
	return s.HasAny(seq, filter)
}

// Contains finds the first element that equal to the example and returns true.
func Contains[S ~seq[T], T comparable](seq S, example T) bool {
	return s.Contains(seq, example)
}

// Conv creates an iterator that applies the 'converter' function to each iterable element and returns value-error pairs.
// The error should be checked at every iteration step, like:
//
//	var integers iter.Seq[int]
//	...
//	for s, err := range seq.Conv(integers,  strconv.Itoa) {
//	    if err != nil {
//	        break
//	    }
//	    ...
//	}
func Conv[S ~seq[From], From, To any](seq S, converter func(From) (To, error)) SeqE[To] {
	return SeqE[To](ToSeq2(seq, converter))
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func Convert[S ~seq[From], From, To any](seq S, converter func(From) To) Seq[To] {
	return s.Convert(seq, converter)
}

// ConvertOK creates an iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from the loop.
func ConvertOK[S ~seq[From], From, To any](seq S, converter func(from From) (To, bool)) Seq[To] {
	return s.ConvertOK(seq, converter)
}

// ConvOK creates a iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from iteration.
// It may also return an error to abort the iteration.
func ConvOK[S ~seq[From], From, To any](seq S, converter func(from From) (To, bool, error)) SeqE[To] {
	return s.ConvOK(seq, converter)
}

// Flat is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var arrays iter.Seq[[]int]
//	...
//	for e := range seq.Flat(arrays, as.Is) {
//	    ...
//	}
func Flat[S ~seq[From], STo ~[]To, From any, To any](seq S, flattener func(From) STo) Seq[To] {
	return s.Flat(seq, flattener)
}

// FlatSeq is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var arrays iter.Seq[[]int]
//	...
//	for e := range s.FlatSeq(arrays, slices.Values) {
//	    ...
//	}
func FlatSeq[S ~seq[From], STo ~seq[To], From any, To any](seq S, flattener func(From) STo) Seq[To] {
	return s.FlatSeq(seq, flattener)
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
func Flatt[S ~seq[From], STo ~[]To, From any, To any](seq S, flattener func(From) (STo, error)) SeqE[To] {
	return s.Flatt(seq, flattener)
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
func FlattSeq[S ~seq[From], STo ~seqE[To], From any, To any](seq S, flattener func(From) STo) SeqE[To] {
	return s.FlattSeq(seq, flattener)
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~seq[T], T any](seq S, filter func(T) bool) Seq[T] {
	return s.Filter(seq, filter)
}

// Filt creates an erroreable iterator that iterates only those elements for which the 'filter' function returns true.
func Filt[S ~seq[T], T any](seq S, filter func(T) (bool, error)) SeqE[T] {
	return s.Filt(seq, filter)
}

// ToKV converts the seq iterator to a key/value pairs iterator by applying the key, value extractors to each iterable element.
func ToKV[S ~seq[T], T, K, V any](seq S, keyExtractor func(T) K, valExtractor func(T) V) Seq2[K, V] {
	return ToSeq2(seq, func(t T) (K, V) { return keyExtractor(t), valExtractor(t) })
}

// KeyValues converts the seq iterator to a key/value pairs iterator by applying the key, values extractors to each iterable element.
func KeyValues[S ~seq[T], T, K, V any](seq S, keyExtractor func(T) K, valsExtractor func(T) []V) Seq2[K, V] {
	return s.KeyValues(seq, keyExtractor, valsExtractor)
}

// Group collects the seq elements into a new map.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to a value.
func Group[S ~seq[T], T any, K comparable, V any](seq S, keyExtractor func(T) K, valExtractor func(T) V) map[K][]V {
	return s2.Group(ToKV(seq, keyExtractor, valExtractor))
}

// ForEach applies the 'consumer' function to the seq elements
func ForEach[T any](seq Seq[T], consumer func(T)) {
	s.ForEach(seq, consumer)
}
