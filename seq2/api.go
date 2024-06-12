// Package seq2 provides helpers for  “range-over-func” feature introduced in go 1.22.
package seq2

import "iter"

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
	return func(consumer func(K, V) bool) {
		if seq == nil {
			return
		}
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
	return func(consumer func(Kto, Vto) bool) {
		if seq == nil {
			return
		}
		seq(func(k Kfrom, v Vfrom) bool {
			return consumer(converter(k, v))
		})
	}
}

// Values converts a key/value pairs iterator to an iterator of just values.
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
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
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		if seq == nil {
			return
		}
		seq(func(k K, _ V) bool {
			return yield(k)
		})
	}
}
