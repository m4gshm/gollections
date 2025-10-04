// Package seq extends [iter.Seq] API with convering, filtering, and reducing functionality.
package seq

import (
	"github.com/m4gshm/gollections/convert"
	s2 "github.com/m4gshm/gollections/internal/seq2"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/predicate/always"
	"golang.org/x/exp/constraints"
)

// Seq is an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[T any] seq[T]

type seq[T any] = func(func(T) bool)

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
type SeqE[T any] seqE[T]
type seqE[T any] = seq2[T, error]

// Seq2 is an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] seq2[K, V]

type seq2[K, V any] = func(func(K, V) bool)

// Of creates an iterator over the elements.
func Of[T any](elements ...T) Seq[T] {
	v := func(yield func(T) bool) {
		for _, v := range elements {
			if !yield(v) {
				break
			}
		}
	}
	return Seq[T](v)
}

// Of2 creates an index/value pairs iterator over the elements.
func Of2[T any](elements ...T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range elements {
			if !yield(i, v) {
				break
			}
		}
	}
}

// Union combines several sequences into one.
func Union[S ~seq[T], T any](seq ...S) Seq[T] {
	return func(yield func(T) bool) {
		for _, s := range seq {
			if s != nil {
				for v := range s {
					if !yield(v) {
						return
					}
				}
			}
		}
	}
}

// OfNextGet builds an iterator by iterating elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfNextGet[T any](hasNext func() bool, getNext func() T) Seq[T] {
	return func(yield func(T) bool) {
		for hasNext() {
			if o := getNext(); !yield(o) {
				return
			}
		}
	}
}

// OfNext builds an iterator by iterating elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The pushNext copy the element to the next pointer.
func OfNext[T any](hasNext func() bool, pushNext func(*T)) Seq[T] {
	return OfNextGet(hasNext, func() (o T) { pushNext(&o); return o })
}

// OfSourceNextGet builds an iterator by iterating elements of the source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfSourceNextGet[S, T any](source S, hasNext func(S) bool, getNext func(S) T) Seq[T] {
	return OfNextGet(func() bool { return hasNext(source) }, func() T { return getNext(source) })
}

// OfSourceNext builds an iterator by iterating elements of the source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The pushNext copy the element to the next pointer.
func OfSourceNext[S, T any](source S, hasNext func(S) bool, pushNext func(S, *T)) Seq[T] {
	return OfNext(func() bool { return hasNext(source) }, func(next *T) { pushNext(source, next) })
}

// OfIndexed builds a Seq iterator by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](amount int, getAt func(int) T) Seq[T] {
	return func(yield func(T) bool) {
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

// Series makes a sequence by applying the 'next' function to the previous step generated value.
func Series[T any](first T, next func(T) (T, bool)) Seq[T] {
	return func(yield func(T) bool) {
		if next == nil {
			return
		}
		current := first
		if !yield(current) {
			return
		}
		for {
			next, ok := next(current)
			if !ok {
				break
			}
			if !yield(next) {
				break
			}
			current = next
		}
	}
}

// RangeClosed creates a sequence that generates integers in the range defined by from and to inclusive
func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) Seq[T] {
	amount := toInclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	amount++
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

// Range creates a sequence that generates integers in the range defined by from and to exclusive
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
func ToSeq2[S ~seq[T], T, K, V any](seq S, converter func(T) (K, V)) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(v T) bool {
			return yield(converter(v))
		})
	}
}

// Top returns a sequence of top n elements.
func Top[S ~seq[T], T any](n int, seq S) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		m := n
		seq(func(t T) bool {
			if m == 0 {
				return false
			}
			m--
			return yield(t)
		})
	}
}

// Skip returns the seq without first n elements.
func Skip[S ~seq[T], T any](n int, seq S) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		m := n
		seq(func(t T) bool {
			if m == 0 {
				return yield(t)
			}
			m--
			return true
		})
	}
}

// While cuts tail elements of the seq that don't match the filter.
func While[S ~seq[T], T any](seq S, filter func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		seq(func(t T) bool {
			if !filter(t) {
				return false
			}
			return yield(t)
		})
	}
}

// SkipWhile returns a sequence without first elements of the seq that dont'math the filter.
func SkipWhile[S ~seq[T], T any](seq S, filter func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		started := false
		seq(func(t T) bool {
			if !started && filter(t) {
				return true
			}
			started = true
			return yield(t)
		})
	}
}

// Head returns the first element.
func Head[S ~seq[T], T any](seq S) (v T, ok bool) {
	return First(seq, always.True)
}

// First returns the first element that satisfies the condition.
func First[S ~seq[T], T any](seq S, condition func(T) bool) (v T, ok bool) {
	if seq == nil || condition == nil {
		return
	}
	seq(func(one T) bool {
		if condition(one) {
			v = one
			ok = true
			return false
		}
		return true
	})
	return
}

// Firstt returns the first element that satisfies the condition.
func Firstt[S ~seq[T], T any](seq S, condition func(T) (bool, error)) (v T, ok bool, err error) {
	if seq == nil || condition == nil {
		return v, false, nil
	}
	seq(func(one T) bool {
		ok, err = condition(one)
		if ok {
			v = one
			return false
		} else if err != nil {
			return false
		}
		return true
	})
	return v, ok, err
}

// Slice collects the elements of the 'seq' sequence into a new slice.
func Slice[S ~seq[T], T any](seq S) []T {
	return SliceCap(seq, 0)
}

// SliceCap collects the elements of the 'seq' sequence into a new slice with predefined capacity.
func SliceCap[S ~seq[T], T any](seq S, capacity int) (out []T) {
	if capacity > 0 {
		out = make([]T, 0, capacity)
	}
	return Append(seq, out)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func Append[S ~seq[T], TS ~[]T, T any](seq S, out TS) TS {
	if seq == nil {
		return out
	}
	seq(func(v T) bool {
		out = append(out, v)
		return true
	})
	return out
}

// Reduce reduces the elements of the seq into one using the 'merge' function.
func Reduce[S ~seq[T], T any](seq S, merge func(T, T) T) T {
	result, _ := ReduceOK(seq, merge)
	return result
}

// ReduceOK reduces the elements of the seq into one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceOK[S ~seq[T], T any](seq S, merge func(T, T) T) (result T, ok bool) {
	if seq == nil || merge == nil {
		return result, false
	}
	started := false
	seq(func(v T) bool {
		if !started {
			result = v
			started = true
		} else {
			result = merge(result, v)
		}
		return true
	})
	return result, started
}

// Reducee reduces the elements of the seq into one using the 'merge' function.
func Reducee[S ~seq[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

// ReduceeOK reduces the elements of the seq into one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceeOK[S ~seq[T], T any](seq S, merge func(T, T) (T, error)) (result T, ok bool, err error) {
	if seq == nil || merge == nil {
		return result, false, nil
	}
	started := false
	seq(func(v T) bool {
		if !started {
			result = v
			started = true
			return true
		}
		result, err = merge(result, v)
		return err == nil
	})
	return result, started, err
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accum[T any, S ~seq[T]](first T, seq S, merge func(T, T) T) T {
	accumulator := first
	if seq == nil || merge == nil {
		return accumulator
	}

	seq(func(v T) bool {
		accumulator = merge(accumulator, v)
		return true
	})
	return accumulator
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accumm[T any, S ~seq[T]](first T, seq S, merge func(T, T) (T, error)) (accumulator T, err error) {
	accumulator = first
	if seq == nil || merge == nil {
		return accumulator, nil
	}
	seq(func(v T) bool {
		accumulator, err = merge(accumulator, v)
		return err == nil
	})
	return accumulator, err

}

// Sum returns the sum of all elements.
func Sum[S ~seq[T], T op.Summable](seq S) (out T) {
	return Accum(out, seq, op.Sum[T])
}

// HasAny checks whether the seq contains an element that satisfies the condition.
func HasAny[S ~seq[T], T any](seq S, filter func(T) bool) bool {
	_, ok := First(seq, filter)
	return ok
}

// Contains finds the first element that equal to the example and returns true.
func Contains[S ~seq[T], T comparable](seq S, example T) bool {
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

// Conv creates an errorable seq that applies the 'converter' function to the iterable elements.
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
	return func(yield func(To) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(from From) bool {
			return yield(converter(from))
		})
	}
}

// ConvertNilSafe creates a seq that filters not nil elements, converts that ones, filters not nils after converting and returns them.
func ConvertNilSafe[S ~seq[*From], From, To any](seq S, converter func(*From) *To) Seq[*To] {
	return ConvertOK(seq, convert.NilSafe(converter))
}

// ConvertOK creates an iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from the sequence.
func ConvertOK[S ~seq[From], From, To any](seq S, converter func(from From) (To, bool)) Seq[To] {
	return func(yield func(To) bool) {
		if seq == nil || converter == nil {
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
func ConvOK[S ~seq[From], From, To any](seq S, converter func(from From) (To, bool, error)) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || converter == nil {
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
func Flat[S ~seq[From], STo ~[]To, From any, To any](seq S, flattener func(From) STo) Seq[To] {
	return func(yield func(To) bool) {
		if seq == nil || flattener == nil {
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
//	for e := range s.FlatSeq(arrays, slices.Values) {
//	    ...
//	}
func FlatSeq[S ~seq[From], STo ~seq[To], From any, To any](seq S, flattener func(From) STo) Seq[To] {
	return func(yield func(To) bool) {
		if seq == nil || flattener == nil {
			return
		}
		seq(func(v From) bool {
			if elementsTo := flattener(v); elementsTo != nil {
				for e := range elementsTo {
					if !yield(e) {
						return false
					}
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
	return func(yield func(To, error) bool) {
		if seq == nil || flattener == nil {
			return
		}
		seq(func(v From) bool {
			elementsTo, err := flattener(v)
			if err != nil && len(elementsTo) == 0 {
				var t To
				return yield(t, err)
			}
			for _, e := range elementsTo {
				if !yield(e, err) {
					return false
				}
			}
			return true
		})
	}
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
	return func(yield func(To, error) bool) {
		if seq == nil || flattener == nil {
			return
		}
		seq(func(v From) bool {
			if elementsTo := flattener(v); elementsTo != nil {
				for e, err := range elementsTo {
					if !yield(e, err) {
						return false
					}
				}
			}
			return true
		})
	}
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~seq[T], T any](seq S, filter func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil || filter == nil {
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
func Filt[S ~seq[T], T any](seq S, filter func(T) (bool, error)) SeqE[T] {
	return func(yield func(T, error) bool) {
		if seq == nil || filter == nil {
			return
		}
		seq(func(t T) bool {
			if ok, err := filter(t); ok || err != nil {
				return yield(t, err)
			}
			return true
		})
	}
}

// ToKV converts the seq iterator to a key/value pairs iterator by applying the key, value extractors to each iterable element.
func ToKV[S ~seq[T], T, K, V any](seq S, keyExtractor func(T) K, valExtractor func(T) V) Seq2[K, V] {
	return ToSeq2(seq, func(t T) (K, V) { return keyExtractor(t), valExtractor(t) })
}

// KeyValues converts the seq iterator to a key/value pairs iterator by applying the key, values extractors to each iterable element.
func KeyValues[S ~seq[T], T, K, V any](seq S, keyExtractor func(T) K, valsExtractor func(T) []V) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil || keyExtractor == nil || valsExtractor == nil {
			return
		}
		for t := range seq {
			k := keyExtractor(t)
			values := valsExtractor(t)
			for _, v := range values {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Group collects the seq elements into a new map.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to a value.
func Group[S ~seq[T], T any, K comparable, V any](seq S, keyExtractor func(T) K, valExtractor func(T) V) map[K][]V {
	return s2.Group(ToKV(seq, keyExtractor, valExtractor))
}

// NotNil returns teh seq without nil elements.
func NotNil[T any](seq Seq[*T]) Seq[*T] {
	return Filter(seq, not.Nil[T])
}

// ForEach applies the 'consumer' function to the seq elements.
func ForEach[T any](seq Seq[T], consumer func(T)) {
	if seq == nil {
		return
	}
	for v := range seq {
		consumer(v)
	}
}
