// Package seq2 extends [iter.Seq2] API with convering, filtering, and reducing functionality.
package seq2

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	s2 "github.com/m4gshm/gollections/internal/seq2"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/seq"
)

// Seq2 is an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] = func(func(K, V) bool)

// Of creates an index/value pairs iterator over the elements.
func Of[T any](elements ...T) seq.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range elements {
			if !yield(i, v) {
				break
			}
		}
	}
}

// OfMap creates an key/value pairs iterator over the elements map.
func OfMap[K comparable, V any](elements map[K]V) seq.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range elements {
			if !yield(k, v) {
				break
			}
		}
	}
}

// Union combines several sequences into one.
func Union[S ~Seq2[K, V], K, V any](seq ...S) seq.Seq2[K, V] {
	return s2.Union(seq...)
}

// OfIndexed builds an indexed Seq2 iterator by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](amount int, getAt func(int) T) seq.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if getAt == nil {
			return
		}
		for i := range amount {
			if !yield(i, getAt(i)) {
				break
			}
		}
	}
}

func OfIndexedKV[K, V any](amount int, getAt func(int) (K, V)) seq.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if getAt == nil {
			return
		}
		for i := range amount {
			if !yield(getAt(i)) {
				break
			}
		}
	}
}

func OfIndexedPair[K, V any](amount int, getKey func(int) K, getValue func(int) V) seq.Seq2[K, V] {
	return OfIndexedKV(amount, func(i int) (K, V) { return getKey(i), getValue(i) })
}

// Series makes a sequence by applying the 'next' function to the previous step generated value.
func Series[T any](first T, next func(int, T) (T, bool)) seq.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if next == nil {
			return
		}
		i := 0
		current := first
		if !yield(i, current) {
			return
		}
		for {
			i++
			next, ok := next(i, current)
			if !ok {
				break
			}
			if !yield(i, next) {
				break
			}
			current = next
		}
	}
}

// RangeClosed creates a sequence that generates integers in the range defined by from and to inclusive
func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) seq.Seq2[int, T] {
	amount := toInclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	amount++
	return func(yield func(int, T) bool) {
		e := from
		for i := 0; i < int(amount); i++ {
			if !yield(i, e) {
				return
			}
			e = e + delta
		}
	}
}

// Range creates a sequence that generates integers in the range defined by from and to exclusive.
func Range[T constraints.Integer | rune](from T, toExclusive T) seq.Seq2[int, T] {
	amount := toExclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	return func(yield func(int, T) bool) {
		e := from
		for i := 0; i < int(amount); i++ {
			if !yield(i, e) {
				return
			}
			e = e + delta
		}
	}
}

// ToSeq converts an iterator of key/value pairs elements to an iterator of single elements by applying the 'converter' function to each iterable pair.
func ToSeq[S ~Seq2[K, V], T, K, V any](seq S, converter func(K, V) T) seq.Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(k K, v V) bool {
			return yield(converter(k, v))
		})
	}
}

// Top returns a sequence of top n key\value pairs.
func Top[S ~Seq2[K, V], K, V any](n int, seq S) S {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		m := n
		seq(func(k K, v V) bool {
			if m == 0 {
				return false
			}
			m--
			return yield(k, v)
		})
	}
}

// Skip returns a sequence without first n elements.
func Skip[S ~Seq2[K, V], K, V any](n int, seq S) seq.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		m := n
		seq(func(k K, v V) bool {
			if m == 0 {
				return yield(k, v)
			}
			m--
			return true
		})
	}
}

// While cuts tail elements of the seq that don't match the filter.
func While[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) seq.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		seq(func(k K, v V) bool {
			if !filter(k, v) {
				return false
			}
			return yield(k, v)
		})
	}
}

// SkipWhile returns a sequence without first elements of the seq that dont'math the filter.
func SkipWhile[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) seq.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		started := false
		seq(func(k K, v V) bool {
			if !started && filter(k, v) {
				return true
			}
			started = true
			return yield(k, v)
		})
	}
}

// Head returns the first key\value pair.
func Head[S ~Seq2[K, V], K, V any](seq S) (k K, v V, ok bool) {
	return First(seq, func(K, V) bool { return true })
}

// First returns the first key\value pair that satisfies the condition.
func First[S ~Seq2[K, V], K, V any](seq S, condition func(K, V) bool) (k K, v V, ok bool) {
	return s2.First(seq, condition)
}

// Firstt returns the first key\value pair that satisfies the condition.
func Firstt[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) (bool, error)) (k K, v V, ok bool, err error) {
	return s2.Firstt(seq, filter)
}

// Reduce reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reduce[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) T) T {
	result, _ := ReduceOK(seq, merge)
	return result
}

// ReduceOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceOK[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) T) (result T, ok bool) {
	if seq == nil || merge == nil {
		return result, false
	}
	started := false
	seq(func(k K, v V) bool {
		result = merge(op.IfElse(!started, nil, &result), k, v)
		started = true
		return true
	})
	return result, started
}

// Reducee reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reducee[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

// ReduceeOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceeOK[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) (T, error)) (result T, ok bool, err error) {
	if seq == nil || merge == nil {
		return result, false, nil
	}
	started := false
	seq(func(k K, v V) bool {
		result, err = merge(op.IfElse(!started, nil, &result), k, v)
		started = true
		return err == nil
	})
	return result, started, err
}

func HasAny[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) bool {
	return s2.HasAny(seq, filter)
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) seq.Seq2[K, V] {
	return s2.Filter(seq, filter)
}

// Filt creates an erroreable iterator that iterates only those key\value pairs for which the 'filter' function returns true.
func Filt[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) (bool, error)) seq.SeqE[c.KV[K, V]] {
	return s2.Filt(seq, filter)
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func Convert[S ~Seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto)) seq.Seq2[Kto, Vto] {
	return s2.Convert(seq, converter)
}

func Conv[S ~Seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto, error)) seq.SeqE[c.KV[Kto, Vto]] {
	return s2.Conv(seq, converter)
}

// Values converts a key/value pairs iterator to an iterator of just values.
func Values[S ~Seq2[K, V], K, V any](seq S) seq.Seq[V] {
	return s2.Values(seq)
}

// Keys converts a key/value pairs iterator to an iterator of just keys.
func Keys[S ~Seq2[K, V], K, V any](seq S) seq.Seq[K] {
	return s2.Keys(seq)
}

// Group collects the elements of the 'seq' sequence into a new map.
func Group[S ~Seq2[K, V], K comparable, V any](seq S) map[K][]V {
	return s2.Group(seq)
}

// Map collects key\value elements into a new map by iterating over the elements.
func Map[S ~Seq2[K, V], K comparable, V any](seq S) map[K]V {
	return s2.MapResolv(seq, resolv.First[K, V])
}

// MapResolv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
func MapResolv[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) map[K]VR {
	return s2.MapResolv(seq, resolver)
}

// MapResolvOrder collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
// Returns a slice with the keys ordered by the time they were added and the resolved key\value map.
func MapResolvOrder[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) ([]K, map[K]VR) {
	return s2.MapResolvOrder(seq, resolver)
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
