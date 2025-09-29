// Package seq2 extends [iter.Seq2] API with convering, filtering, and reducing functionality.
package seq2

import (
	"github.com/m4gshm/gollections/c"
	s "github.com/m4gshm/gollections/internal/seq"
	s2 "github.com/m4gshm/gollections/internal/seq2"
	"github.com/m4gshm/gollections/map_/resolv"
	"golang.org/x/exp/constraints"
)

// Seq is an alias of an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[T any] = s.Seq[T]

// Seq2 is an alias of an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] = s.Seq2[K, V]

// Of creates an index/value pairs iterator over the elements.
func Of[T any](elements ...T) Seq2[int, T] {
	return s2.Of(elements...)
}

// OfMap creates an key/value pairs iterator over the elements map.
func OfMap[K comparable, V any](elements map[K]V) Seq2[K, V] {
	return s2.OfMap(elements)
}

// Union combines several sequences into one.
func Union[S ~Seq2[K, V], K, V any](seq ...S) Seq2[K, V] {
	return s2.Union(seq...)
}

// OfIndexed builds an indexed Seq2 iterator by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](amount int, getAt func(int) T) Seq2[int, T] {
	return s2.OfIndexed(amount, getAt)
}

func OfIndexedKV[K, V any](amount int, getAt func(int) (K, V)) Seq2[K, V] {
	return s2.OfIndexedKV(amount, getAt)
}

func OfIndexedPair[K, V any](amount int, getKey func(int) K, getValue func(int) V) Seq2[K, V] {
	return OfIndexedKV(amount, func(i int) (K, V) { return getKey(i), getValue(i) })
}

// Series makes a sequence by applying the 'next' function to the previous step generated value.
func Series[T any](first T, next func(int, T) (T, bool)) Seq2[int, T] {
	return s2.Series(first, next)
}

// RangeClosed creates a sequence that generates integers in the range defined by from and to inclusive
func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) Seq2[int, T] {
	return s2.RangeClosed(from, toInclusive)
}

// Range creates a sequence that generates integers in the range defined by from and to exclusive
func Range[T constraints.Integer | rune](from T, toExclusive T) Seq2[int, T] {
	return s2.Range(from, toExclusive)
}

// ToSeq converts an iterator of key/value pairs elements to an iterator of single elements by applying the 'converter' function to each iterable pair.
func ToSeq[S ~Seq2[K, V], T, K, V any](seq S, converter func(K, V) T) Seq[T] {
	return s2.ToSeq(seq, converter)
}

// Top returns a sequence of top n key\value pairs.
func Top[S ~Seq2[K, V], K, V any](n int, seq S) S {
	return s2.Top(n, seq)
}

// Skip returns a sequence without first n elements.
func Skip[S ~Seq2[K, V], K, V any](n int, seq S) Seq2[K, V] {
	return s2.Skip(n, seq)
}

// While cuts tail elements of the seq that don't match the filter.
func While[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) Seq2[K, V] {
	return s2.While(seq, filter)
}

// SkipWhile returns a sequence without first elements of the seq that dont'math the filter.
func SkipWhile[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) Seq2[K, V] {
	return s2.SkipWhile(seq, filter)
}

// Head returns the first key\value pair.
func Head[S ~Seq2[K, V], K, V any](seq S) (k K, v V, ok bool) {
	return First(seq, func(K, V) bool { return true })
}

// First returns the first key\value pair that satisfies the condition of the 'filter' function.
func First[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) (k K, v V, ok bool) {
	return s2.First(seq, filter)
}

// Reduce reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reduce[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) T) T {
	result, _ := ReduceOK(seq, merge)
	return result
}

// ReduceOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceOK[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) T) (result T, ok bool) {
	return s2.ReduceOK(seq, merge)
}

// Reducee reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reducee[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

// ReduceeOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceeOK[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) (T, error)) (result T, ok bool, err error) {
	return s2.ReduceeOK(seq, merge)
}

func HasAny[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) bool {
	return s2.HasAny(seq, filter)
}

// Filter creates a rangefunc that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) Seq2[K, V] {
	return s2.Filter(seq, filter)
}

func Filt[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) (bool, error)) Seq2[c.KV[K, V], error] {
	return s2.Filt(seq, filter)
}

// Convert creates a rangefunc that applies the 'converter' function to each iterable element.
func Convert[S ~Seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto)) Seq2[Kto, Vto] {
	return s2.Convert(seq, converter)
}

func Conv[S ~Seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto, error)) Seq2[c.KV[Kto, Vto], error] {
	return s2.Conv(seq, converter)
}

// Values converts a key/value pairs iterator to an iterator of just values.
func Values[S ~Seq2[K, V], K, V any](seq S) Seq[V] {
	return s2.Values(seq)
}

// Keys converts a key/value pairs iterator to an iterator of just keys.
func Keys[S ~Seq2[K, V], K, V any](seq S) Seq[K] {
	return s2.Keys(seq)
}

// Group collects the elements of the 'seq' sequence into a new map.
func Group[S ~Seq2[K, V], K comparable, V any](seq S) map[K][]V {
	return MapResolv(seq, resolv.Slice[K, V])
}

// Map collects key\value elements into a new map by iterating over the elements.
func Map[S ~Seq2[K, V], K comparable, V any](seq S) map[K]V {
	return MapResolv(seq, resolv.First[K, V])
}

// MapResolv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
func MapResolv[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) map[K]VR {
	return AppendMapResolv(seq, resolver, nil)
}

// MapResolvOrder collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
// Returns a slice with the keys ordered by the time they were added and the resolved key\value map.
func MapResolvOrder[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) ([]K, map[K]VR) {
	return AppendMapResolvOrder(seq, resolver, nil, nil)
}

// AppendMapResolv collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values.
func AppendMapResolv[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR, dest map[K]VR) map[K]VR {
	return s2.AppendMapResolv(seq, resolver, dest)
}

// AppendMapResolvOrder collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values
// Additionaly populates the 'order' slice by the keys ordered by the time they were added and the resolved key\value map.
func AppendMapResolvOrder[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR, order []K, dest map[K]VR) ([]K, map[K]VR) {
	return s2.AppendMapResolvOrder(seq, resolver, order, dest)
}

// TrackEach applies the 'consumer' function to the seq elements.
func TrackEach[S ~Seq2[K, V], K, V any](seq S, consumer func(K, V)) {
	s2.TrackEach(seq, consumer)
}
