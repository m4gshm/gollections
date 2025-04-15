// Package seq2 provides helpers for “range-over-func” feature introduced in go 1.22.
package seq2

import (
	"iter"

	"github.com/m4gshm/gollections/map_/resolv"
)

// Of creates an index/value pairs iterator over the elements.
func Of[T any](elements ...T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range elements {
			if ok := yield(i, v); !ok {
				break
			}
		}
	}
}

// OfMap creates an key/value pairs iterator over the elements map.
func OfMap[K comparable, V any](elements map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range elements {
			if ok := yield(k, v); !ok {
				break
			}
		}
	}
}

// Filter creates a rangefunc that iterates only those elements for which the 'filter' function returns true.
func Filter[K, V any](seq iter.Seq2[K, V], filter func(K, V) bool) iter.Seq2[K, V] {
	if seq == nil {
		return func(_ func(K, V) bool) {}
	}
	return func(consumer func(K, V) bool) {
		seq(func(k K, v V) bool {
			if filter(k, v) {
				return consumer(k, v)
			}
			return true
		})
	}
}

// Convert creates a rangefunc that applies the 'converter' function to each iterable element.
func Convert[Kfrom, Vfrom, Kto, Vto any](seq iter.Seq2[Kfrom, Vfrom], converter func(Kfrom, Vfrom) (Kto, Vto)) iter.Seq2[Kto, Vto] {
	if seq == nil {
		return func(_ func(Kto, Vto) bool) {}
	}
	return func(consumer func(Kto, Vto) bool) {
		seq(func(k Kfrom, v Vfrom) bool {
			return consumer(converter(k, v))
		})
	}
}

// Values converts a key/value pairs iterator to an iterator of just values.
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	if seq == nil {
		return func(_ func(V) bool) {}
	}
	return func(yield func(V) bool) {
		seq(func(_ K, v V) bool {
			return yield(v)
		})
	}
}

// Keys converts a key/value pairs iterator to an iterator of just keys.
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	if seq == nil {
		return func(_ func(K) bool) {}
	}
	return func(yield func(K) bool) {
		seq(func(k K, _ V) bool {
			return yield(k)
		})
	}
}

// Slice collects the elements of the 'seq' sequence into a new slice.
func Slice[T any](seq iter.Seq2[T, error]) ([]T, error) {
	return SliceCap(seq, 0)
}

// SliceCap collects the elements of the 'seq' sequence into a new slice with predefined capacity.
func SliceCap[T any](seq iter.Seq2[T, error], cap int) (out []T, e error) {
	if seq == nil {
		return nil, nil
	}
	if cap > 0 {
		out = make([]T, 0, cap)
	}
	return Append(seq, out)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func Append[T any, TS ~[]T](seq iter.Seq2[T, error], out TS) (TS, error) {
	if seq == nil {
		return nil, nil
	}
	var errOur error
	seq(func(v T, e error) bool {
		if e != nil {
			errOur = e
			return false
		}
		out = append(out, v)
		return true
	})
	return out, errOur
}

// Group converts the elements of the 'seq' sequence into a new map, extracting a key for each element applying the converter 'keyExtractor'.
func Group[K comparable, V any](seq iter.Seq2[K, V]) map[K][]V {
	return MapResolv(seq, resolv.Slice[K, V])
}

// Map collects key\value elements into a new map by iterating over the elements.
func Map[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	return MapResolv(seq, resolv.First[K, V])
}

// MapResolv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
func MapResolv[K comparable, V, VR any](seq iter.Seq2[K, V], resolver func(bool, K, VR, V) VR) map[K]VR {
	return AppendMapResolv(seq, resolver, nil)
}

// MapResolvOrder collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
// Returns a slice with the keys ordered by the time they were added and the resolved key\value map.
func MapResolvOrder[K comparable, V, VR any](seq iter.Seq2[K, V], resolver func(bool, K, VR, V) VR) ([]K, map[K]VR) {
	return AppendMapResolvOrder(seq, resolver, nil, nil)
}

// AppendMapResolv collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values.
func AppendMapResolv[K comparable, V, VR any](seq iter.Seq2[K, V], resolver func(bool, K, VR, V) VR, dest map[K]VR) map[K]VR {
	if seq == nil {
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
func AppendMapResolvOrder[K comparable, V, VR any](seq iter.Seq2[K, V], resolver func(bool, K, VR, V) VR, order []K, dest map[K]VR) ([]K, map[K]VR) {
	if seq == nil {
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
