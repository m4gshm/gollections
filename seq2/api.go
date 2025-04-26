// Package seq2 extends [iter.Seq2] API with convering, filtering, and reducing functionality.
package seq2

import (
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/seq"
	"golang.org/x/exp/constraints"
)

// Seq is an alias of an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[T any] = seq.Seq[T]

// Seq2 is an alias of an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] = seq.Seq2[K, V]

// Of creates an index/value pairs iterator over the elements.
func Of[T any](elements ...T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range elements {
			if !yield(i, v) {
				break
			}
		}
	}
}

// OfMap creates an key/value pairs iterator over the elements map.
func OfMap[K comparable, V any](elements map[K]V) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range elements {
			if !yield(k, v) {
				break
			}
		}
	}
}

// OfIndexed builds an indexed Seq2 iterator by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](amount int, getAt func(int) T) Seq2[int, T] {
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

// Series makes a sequence by applying the 'next' function to the previous step generated value.
func Series[T any](first T, next func(T) (T, bool)) Seq2[int, T] {
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
			next, ok := next(current)
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
func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) Seq2[int, T] {
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

// Range creates a sequence that generates integers in the range defined by from and to exclusive
func Range[T constraints.Integer | rune](from T, toExclusive T) Seq2[int, T] {
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


// Filter creates a rangefunc that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil || filter == nil {
			return
		}
		seq(func(k K, v V) bool {
			if filter(k, v) {
				return yield(k, v)
			}
			return true
		})
	}
}

// Convert creates a rangefunc that applies the 'converter' function to each iterable element.
func Convert[S ~Seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto)) Seq2[Kto, Vto] {
	return func(consumer func(Kto, Vto) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(k Kfrom, v Vfrom) bool {
			return consumer(converter(k, v))
		})
	}
}

// Values converts a key/value pairs iterator to an iterator of just values.
func Values[S ~Seq2[K, V], K, V any](seq S) Seq[V] {
	return func(yield func(V) bool) {
		if seq == nil {
			return
		}
		seq(func(_ K, v V) bool {
			return yield(v)
		})
	}
}

// Keys converts a key/value pairs iterator to an iterator of just keys.
func Keys[S ~Seq2[K, V], K, V any](seq S) Seq[K] {
	return func(yield func(K) bool) {
		if seq == nil {
			return
		}
		seq(func(k K, _ V) bool {
			return yield(k)
		})
	}
}

// Group converts the elements of the 'seq' sequence into a new map, extracting a key for each element applying the converter 'keyExtractor'.
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
	if seq == nil || resolver == nil {
		return nil
	}
	if dest == nil {
		dest = map[K]VR{}
	}
	seq(func(k K, v V) bool {
		exists, ok := dest[k]
		dest[k] = resolver(ok, k, exists, v)
		return true
	})
	return dest
}

// AppendMapResolvOrder collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values
// Additionaly populates the 'order' slice by the keys ordered by the time they were added and the resolved key\value map.
func AppendMapResolvOrder[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR, order []K, dest map[K]VR) ([]K, map[K]VR) {
	if seq == nil || resolver == nil {
		return nil, nil
	}
	if dest == nil {
		dest = map[K]VR{}
	}
	seq(func(k K, v V) bool {
		exists, ok := dest[k]
		dest[k] = resolver(ok, k, exists, v)
		if !ok {
			order = append(order, k)
		}
		return true
	})
	return order, dest
}
