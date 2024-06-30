// Package seq provides helpers for “range-over-func” feature introduced in go 1.22.
package seq

import (
	"iter"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/op"
	"google.golang.org/genproto/googleapis/cloud/functions/v1"
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
	if seq == nil {
		return nil
	}
	return func(consumer func(T) bool) {
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
	if seq == nil {
		return func(yield func(To) bool) {}
	}
	return func(consumer func(To) bool) {
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
	if seq == nil {
		return func(yield func(K, V) bool) {}
	}
	return func(yield func(K, V) bool) {
		seq(func(v T) bool {
			return yield(converter(v))
		})
	}
}

// KeyValue converts an iterator of single elements to a key/value pairs iterator by applying key, value extractors to each iterable element.
func KeyValue[T, K, V any](seq iter.Seq[T], keyExtractor func(T) K, valExtractor func(T) V) iter.Seq2[K, V] {
	return ToSeq2(seq, func(t T) (K, V) { return keyExtractor(t), valExtractor(t) })
}

// Slice collects the elements of the 'seq' sequence into a new slice
func Slice[T any](seq iter.Seq[T]) []T {
	return SliceCap(seq, 0)
}

// SliceCap collects the elements of the 'seq' sequence into a new slice with predefined capacity
func SliceCap[T any](seq iter.Seq[T], cap int) (out []T) {
	if seq == nil {
		return nil
	}
	if cap > 0 {
		out = make([]T, 0, cap)
	}
	return Append(seq, out)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice
func Append[T any, TS ~[]T](seq iter.Seq[T], out TS) TS {
	if seq == nil {
		return nil
	}
	seq(func(v T) bool {
		out = append(out, v)
		return true
	})
	return out
}


func Reduce[T any](seq iter.Seq[T], merge func(T, T) T) T {
	result, _ := ReduceOK(seq, merge)
	return result
}

// ReduceOK reduces the elements retrieved by the 'seq' iterator into an one using the 'merge' function.
func ReduceOK[T any](seq iter.Seq[T], merge func(T, T) T) (result T, ok bool) {
	if seq == nil {
		return result, false
	}
	first := true
	seq(func(v T) bool {
		if first {
			result = v
		} else {
			result = merge(result, v)
		}
		first = false
		return true
	})
	return result, true
}

func Accum[T any](first T, seq iter.Seq[T], merge func(T, T) T) T {
	accumulator := first
	if seq == nil {
		return accumulator
	}

	seq(func (v T) bool {
		accumulator = merge(accumulator, v)
		return true
	})
	return accumulator
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any](seq iter.Seq[T], predicate func(T) bool) (v T, ok bool) {
	if seq == nil {
		return
	}
	seq(func(one T) bool {
		if predicate(one) {
			v = one
			ok = true
			return false
		}
		return true
	})
	return
}

// Sum returns the sum of all elements
func Sum[T c.Summable](seq iter.Seq[T]) (T, bool) {
	return ReduceOK(seq, op.Sum[T])
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func HasAny[T any](seq iter.Seq[T], predicate func(T) bool) bool {
	_, ok := First(seq, predicate)
	return ok
}

// Contains finds the first element that equal to the example and returns true
func Contains[T comparable](seq iter.Seq[T], example T) bool {
	if seq == nil {
		return false
	}
	contains := false
	seq(func(v T) bool {
		contains = v == example
		return !contains
	})
	return contains
}
