// Package seq provides helpers for “range-over-func” feature introduced in go 1.22.
package seq

import (
	"iter"
)

// Of creates an iterator over the elements.
func Of[T any](elements ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range elements {
			if ok := yield(v); !ok {
				break
			}
		}
	}
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[T any](seq iter.Seq[T], filter func(T) bool) iter.Seq[T] {
	return func(consumer func(T) bool) {
		if seq == nil {
			return
		}
		seq(func(e T) bool {
			if filter(e) {
				return consumer(e)
			}
			return true
		})
	}
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func Convert[From, To any](seq iter.Seq[From], converter func(From) To) iter.Seq[To] {
	return func(consumer func(To) bool) {
		if seq == nil {
			return
		}
		seq(func(from From) bool {
			return consumer(converter(from))
		})
	}
}

// Conv creates an iterator that applies the 'converter' function to each iterable element and returns value-error pairs.
// The error should be checked at every iteration step, like:
//
//	var integers []int
//	...
//	for s, err := range seq.Conv(integers,  strconv.Itoa) {
//	    if err != nil {
//	        break
//	    }
//	    ...
//	}
func Conv[From, To any](seq iter.Seq[From], converter func(From) (To, error)) iter.Seq2[To, error] {
	return ToSeq2(seq, converter)
}

// ToSeq2 converts an iterator of single elements to an iterator of key/value pairs by applying the 'converter' function to each iterable element.
func ToSeq2[T, K, V any](seq iter.Seq[T], converter func(T) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		seq(func(v T) bool {
			return yield(converter(v))
		})
	}
}

// KeyValue converts an iterator of single elements to a key/value pairs iterator by applying key, value extractors to each iterable element.
func KeyValue[T, K, V any](seq iter.Seq[T], keyExtractor func(T) K, valExtractor func(T) V) iter.Seq2[K, V] {
	return ToSeq2(seq, func(t T) (K, V) { return keyExtractor(t), valExtractor(t) })
}
