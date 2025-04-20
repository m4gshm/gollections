// Package seq provides helpers for “range-over-func” feature introduced in go 1.22.
package seq

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/op"
	"golang.org/x/exp/constraints"
)

// Seq is an alias of an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[V any] = func(yield func(V) bool)

// SeqE is a specific iterator form tha allows to retrieve a value with an error as second parameter of the iterator.
// It is used as a result of applying functions like seq.Conv, which may throw an error during iteration.
type SeqE[T any] = Seq2[T, error]

// Seq2 is an alias of an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] = func(yield func(K, V) bool)

// Of creates an iterator over the elements.
func Of[T any](elements ...T) Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range elements {
			if ok := yield(v); !ok {
				break
			}
		}
	}
}

// OfIndexed builds a Seq iterator by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](max int, getAt func(int) T) Seq[T] {
	return func(yield func(T) bool) {
		if getAt == nil {
			return
		}
		for i := range max {
			if ok := yield(getAt(i)); !ok {
				break
			}
		}
	}
}

// RangeClosed creates a loop that generates integers in the range defined by from and to inclusive
func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) Seq[T] {
	amount := toInclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	amount++
	nextElement := from
	i := T(0)
	return func(yield func(T) bool) {
		if ok := i < amount; ok {
			out := nextElement
			i++
			nextElement = nextElement + delta
			if !yield(out) {
				return
			}
		}
	}
}

// Range creates a loop that generates integers in the range defined by from and to exclusive
func Range[T constraints.Integer | rune](from T, toExclusive T) Seq[T] {
	amount := toExclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	return func(yield func(T) bool) {
		e := from
		for i := 0; i < int(amount); i++ {
			if !yield(e) {
				return
			}
			e = e + delta
		}

	}
}

// ToSeq2 converts an iterator of single elements to an iterator of key/value pairs by applying the 'converter' function to each iterable element.
func ToSeq2[S ~Seq[T], T, K, V any](seq S, converter func(T) (K, V)) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		seq(func(v T) bool {
			return yield(converter(v))
		})
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function.
func First[S ~Seq[T], T any](seq S, predicate func(T) bool) (v T, ok bool) {
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

// Firstt returns the first element that satisfies the condition of the 'predicate' function.
func Firstt[S ~Seq[T], T any](seq S, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	if seq == nil {
		return v, false, nil
	}
	seq(func(one T) bool {
		p, e := predicate(one)
		if e != nil {
			err = e
			return false
		} else if p {
			v = one
			ok = true
			return false
		}
		return true
	})
	return v, ok, err
}

// Slice collects the elements of the 'seq' sequence into a new slice.
func Slice[S ~Seq[T], T any](seq S) []T {
	return SliceCap(seq, 0)
}

// SliceCap collects the elements of the 'seq' sequence into a new slice with predefined capacity.
func SliceCap[S ~Seq[T], T any](seq S, cap int) (out []T) {
	if seq == nil {
		return nil
	}
	if cap > 0 {
		out = make([]T, 0, cap)
	}
	return Append(seq, out)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func Append[T any, TS ~[]T, S ~Seq[T]](seq S, out TS) TS {
	if seq == nil {
		return nil
	}
	seq(func(v T) bool {
		out = append(out, v)
		return true
	})
	return out
}

// Reduce reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reduce[S ~Seq[T], T any](seq S, merge func(T, T) T) T {
	result, _ := ReduceOK(seq, merge)
	return result
}

// ReduceOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceOK[S ~Seq[T], T any](seq S, merge func(T, T) T) (result T, ok bool) {
	if seq == nil {
		return result, false
	}
	started := false
	seq(func(v T) bool {
		if !started {
			result = v
		} else {
			result = merge(result, v)
		}
		started = true
		return true
	})
	return result, started
}

// Reducee reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reducee[S ~Seq[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

// ReduceeOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceeOK[S ~Seq[T], T any](seq S, merge func(T, T) (T, error)) (result T, ok bool, err error) {
	if seq == nil {
		return result, false, nil
	}
	started := false
	seq(func(v T) bool {
		if !started {
			result = v
		} else {
			result, err = merge(result, v)
			if err != nil {
				return false
			}
		}
		started = true
		return true
	})
	return result, started, err
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accum[T any, S ~Seq[T]](first T, seq S, merge func(T, T) T) T {
	accumulator := first
	if seq == nil {
		return accumulator
	}

	seq(func(v T) bool {
		accumulator = merge(accumulator, v)
		return true
	})
	return accumulator
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accumm[T any, S ~Seq[T]](first T, seq S, merge func(T, T) (T, error)) (accumulator T, err error) {
	accumulator = first
	if seq == nil {
		return accumulator, nil
	}
	seq(func(v T) bool {
		accumulator, err = merge(accumulator, v)
		return err == nil
	})
	return accumulator, err

}

// Sum returns the sum of all elements.
func Sum[S ~Seq[T], T c.Summable](seq S) (out T) {
	return Accum(out, seq, op.Sum[T])
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful.
func HasAny[S ~Seq[T], T any](seq S, predicate func(T) bool) bool {
	_, ok := First(seq, predicate)
	return ok
}

// Contains finds the first element that equal to the example and returns true.
func Contains[S ~Seq[T], T comparable](seq S, example T) bool {
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
func Conv[S ~Seq[From], From, To any](seq S, converter func(From) (To, error)) SeqE[To] {
	return SeqE[To](ToSeq2(seq, converter))
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func Convert[S ~Seq[From], From, To any](seq S, converter func(From) To) Seq[To] {
	return func(yield func(To) bool) {
		if seq == nil {
			return
		}
		seq(func(from From) bool {
			return yield(converter(from))
		})
	}
}

// ConvertOK creates an iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from the loop.
func ConvertOK[S ~Seq[From], From, To any](seq S, converter func(from From) (To, bool)) Seq[To] {
	return func(yield func(To) bool) {
		if seq == nil {
			return
		}
		seq(func(from From) bool {
			if to, ok := converter(from); ok {
				return yield(to)
			}
			return true
		})
	}
}

// ConvOK creates a iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from iteration.
// It may also return an error to abort the iteration.
func ConvOK[S ~Seq[From], From, To any](seq S, converter func(from From) (To, bool, error)) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil {
			return
		}
		seq(func(from From) bool {
			if to, ok, err := converter(from); ok || err != nil {
				return yield(to, err)
			}
			return true
		})
	}
}

// Flat is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var arrays iter.Seq[[]int]
//	...
//	for e := range seq.Flat(arrays, as.Is) {
//	    ...
//	}
func Flat[S ~Seq[From], STo ~[]To, From any, To any](seq S, flattener func(From) STo) Seq[To] {
	return func(yield func(To) bool) {
		if seq == nil {
			return
		}
		seq(func(v From) bool {
			elementsTo := flattener(v)
			for _, e := range elementsTo {
				if !yield(e) {
					return false
				}
			}
			return true
		})
	}
}

// FlatSeq is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var arrays iter.Seq[[]int]
//	...
//	for e := range seq.FlatSeq(arrays, slices.Values) {
//	    ...
//	}
func FlatSeq[S ~Seq[From], STo ~Seq[To], From any, To any](seq S, flattener func(From) STo) Seq[To] {
	return func(yield func(To) bool) {
		if seq == nil {
			return
		}
		seq(func(v From) bool {
			elementsTo := flattener(v)
			for e := range elementsTo {
				if !yield(e) {
					return false
				}
			}
			return true
		})
	}
}

// Flatt is used to iterate over a two-dimensional sequence in single dimension form, like:
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
func Flatt[S ~Seq[From], STo ~SeqE[To], From any, To any](seq S, flattener func(From) STo) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil {
			return
		}
		seq(func(v From) bool {
			elementsTo := flattener(v)
			for e, err := range elementsTo {
				if !yield(e, err) {
					return false
				}
			}
			return true
		})
	}
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~Seq[T], T any](seq S, filter func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		seq(func(e T) bool {
			if filter(e) {
				return yield(e)
			}
			return true
		})
	}
}

// Filt creates an erroreable iterator that iterates only those elements for which the 'filter' function returns true.
func Filt[S ~Seq[T], T any](seq S, filter func(T) (bool, error)) SeqE[T] {
	//delayed on next iteration step error
	var err error
	return func(yield func(T, error) bool) {
		if seq == nil {
			return
		}
		seq(func(e T) bool {
			if err != nil {
				return yield(e, err)
			}
			ok := false
			ok, err = filter(e)
			if ok {
				return yield(e, nil)
			}
			return true
		})
	}
}

// KeyValue converts the seq iterator to a key/value pairs iterator by applying the key, value extractors to each iterable element.
func KeyValue[S ~Seq[T], T, K, V any](seq S, keyExtractor func(T) K, valExtractor func(T) V) Seq2[K, V] {
	return ToSeq2(seq, func(t T) (K, V) { return keyExtractor(t), valExtractor(t) })
}

// KeyValues converts the seq iterator to a key/value pairs iterator by applying the key, values extractors to each iterable element.
func KeyValues[S ~Seq[T], T, K, V any](seq S, keyExtractor func(T) K, valsExtractor func(T) []V) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		for t := range seq {
			k := keyExtractor(t)
			values := valsExtractor(t)
			for _, v := range values {
				if ok := yield(k, v); !ok {
					return
				}
			}
		}
	}
}
