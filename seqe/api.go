// Package seqe provides convering, filtering, and reducing operations for the seq.SeqE interface.
package seqe

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/seq"
)

// Seq is an alias of an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[T any] = seq.Seq[T]

// SeqE is a specific iterator form that allows to retrieve a value with an error as second parameter of the iterator.
// It is used as a result of applying functions like seq.Conv, which may throw an error during iteration.
type SeqE[T any] = seq.SeqE[T]

// Seq2 is an alias of an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] = seq.Seq2[K, V]

// OfIndexed builds a SeqE iterator by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](amount int, getAt func(int) (T, error)) Seq2[T, error] {
	return func(yield func(T, error) bool) {
		if getAt == nil {
			return
		}
		for i := range amount {
			v, err := getAt(i)
			if !yield(v, err) {
				break
			}
		}
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function.
func First[S ~SeqE[T], T any](seq S, predicate func(T) bool) (v T, ok bool, err error) {
	if seq == nil || predicate == nil {
		return
	}
	seq(func(one T, e error) bool {
		if e != nil {
			err = e
			ok = false
			return false
		} else if predicate(one) {
			v = one
			ok = true
			return false
		}
		return true
	})
	return
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function.
func Firstt[S ~SeqE[T], T any](seq S, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	if seq == nil || predicate == nil {
		return v, false, nil
	}
	seq(func(one T, e error) bool {
		if e != nil {
			err = e
			return false
		} else if ok, err = predicate(one); ok {
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
func Slice[S ~SeqE[T], T any](seq S) ([]T, error) {
	return SliceCap(seq, 0)
}

// SliceCap collects the elements of the 'seq' sequence into a new slice with predefined capacity.
func SliceCap[S ~SeqE[T], T any](seq S, capacity int) (out []T, e error) {
	if seq == nil {
		return nil, nil
	}
	if capacity > 0 {
		out = make([]T, 0, capacity)
	}
	return Append(seq, out)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func Append[S ~SeqE[T], T any, TS ~[]T](seq S, out TS) (TS, error) {
	if seq == nil {
		return out, nil
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

// Reduce reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reduce[S ~SeqE[T], T any](seq S, merge func(T, T) T) (T, error) {
	result, _, err := ReduceOK(seq, merge)
	return result, err
}

// ReduceOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceOK[S ~SeqE[T], T any](seq S, merge func(T, T) T) (result T, ok bool, err error) {
	if seq == nil || merge == nil {
		return result, false, nil
	}
	started := false
	seq(func(v T, e error) bool {
		if e != nil {
			err = e
			return false
		} else if !started {
			result = v
		} else {
			result = merge(result, v)
		}
		started = true
		return true
	})
	return result, started, err
}

// Reducee reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reducee[S ~SeqE[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

// ReduceeOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceeOK[S ~SeqE[T], T any](seq S, merge func(T, T) (T, error)) (result T, ok bool, err error) {
	if seq == nil || merge == nil {
		return result, false, nil
	}
	started := false
	seq(func(v T, e error) bool {
		if e != nil {
			err = e
			return false
		} else if !started {
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
func Accum[T any, S ~SeqE[T]](first T, seq S, merge func(T, T) T) (accumulator T, err error) {
	accumulator = first
	if seq == nil || merge == nil {
		return
	}
	seq(func(v T, e error) bool {
		err = e
		if err != nil {
			return false
		}
		accumulator = merge(accumulator, v)
		return true
	})
	return
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func Accumm[T any, S ~SeqE[T]](first T, seq S, merge func(T, T) (T, error)) (accumulator T, err error) {
	accumulator = first
	if seq == nil || merge == nil {
		return accumulator, nil
	}
	seq(func(v T, e error) bool {
		err = e
		if err == nil {
			accumulator, err = merge(accumulator, v)
		}
		return err == nil
	})
	return accumulator, err
}

// Sum returns the sum of all elements.
func Sum[S ~SeqE[T], T c.Summable](seq S) (out T, err error) {
	return Accum(out, seq, op.Sum[T])
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful.
func HasAny[S ~SeqE[T], T any](seq S, predicate func(T) bool) (bool, error) {
	_, ok, err := First(seq, predicate)
	return ok, err
}

// Contains finds the first element that equal to the example and returns true.
func Contains[S ~SeqE[T], T comparable](seq S, example T) (contains bool, err error) {
	if seq == nil {
		return
	}
	seq(func(v T, e error) bool {
		err = e
		if err != nil {
			return false
		}
		contains = v == example
		return !contains
	})
	return
}

// Conv creates an iterator that applies the 'converter' function to each iterable element and returns value-error pairs.
// The error should be checked at every iteration step, like:
//
//	var integers iter.Seq2[int, error]
//	...
//	for s, err := range seqe.Conv(integers,  strconv.Itoa) {
//	    if err != nil {
//	        break
//	    }
//	    ...
//	}
func Conv[S ~SeqE[From], From, To any](seq S, converter func(From) (To, error)) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(from From, err error) bool {
			if err != nil {
				var to To
				return yield(to, err)
			}
			return yield(converter(from))
		})
	}
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func Convert[S ~SeqE[From], From, To any](seq S, converter func(From) To) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(from From, err error) bool {
			if err != nil {
				var to To
				return yield(to, err)
			}
			return yield(converter(from), err)
		})
	}
}

// ConvertOK creates an iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from the loop.
func ConvertOK[S ~SeqE[From], From, To any](seq S, converter func(from From) (To, bool)) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(from From, err error) bool {
			if err != nil {
				var to To
				return yield(to, err)
			} else if to, ok := converter(from); ok {
				return yield(to, err)
			}
			return true
		})
	}
}

// ConvOK creates a iterator that applies the 'converter' function to each iterable element.
// The converter may returns a value or ok=false to exclude the value from iteration.
// It may also return an error to abort the iteration.
func ConvOK[S ~SeqE[From], From, To any](seq S, converter func(from From) (To, bool, error)) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(from From, e error) bool {
			if to, ok, err := converter(from); ok || err != nil {
				return yield(to, err)
			}
			return true
		})
	}
}

// Flat is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var arrays seq.SeqE[[]int]
//	...
//	for e, err := range seqe.Flat(arrays, as.Is) {
//		if err != nil {
//			panic(err)
//		}
//	}
func Flat[S ~SeqE[From], STo ~[]To, From any, To any](seq S, flattener func(From) STo) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || flattener == nil {
			return
		}
		seq(func(v From, err error) bool {
			if err != nil {
				var t To
				if !yield(t, err) {
					return false
				}
			}
			elementsTo := flattener(v)
			for _, e := range elementsTo {
				if !yield(e, err) {
					return false
				}
			}
			return true
		})
	}
}

// FlatSeq is used to iterate over a two-dimensional sequence in single dimension form, like:
//
//	var arrays seq.SeqE[[]int]
//	...
//	for e, err := range seqe.FlatSeq(arrays, slices.Values) {
//		if err != nil {
//			panic(err)
//		}
//	}
func FlatSeq[S ~SeqE[From], STo ~Seq[To], From any, To any](seq S, flattener func(From) STo) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || flattener == nil {
			return
		}
		seq(func(v From, err error) bool {
			if err != nil {
				var t To
				return yield(t, err)
			}
			if elementsTo := flattener(v); elementsTo != nil {
				for e := range elementsTo {
					if !yield(e, err) {
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
func Flatt[S ~SeqE[From], STo ~[]To, From any, To any](seq S, flattener func(From) (STo, error)) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || flattener == nil {
			return
		}
		seq(func(v From, err error) bool {
			if err != nil {
				var to To
				return yield(to, err)
			}
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
func FlattSeq[S ~SeqE[From], STo ~SeqE[To], From any, To any](seq S, flattener func(From) STo) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || flattener == nil {
			return
		}
		seq(func(v From, err error) bool {
			if err != nil {
				var to To
				return yield(to, err)
			}
			if elementsTo := flattener(v); elementsTo != nil {
				for to, err := range elementsTo {
					if !yield(to, err) {
						return false
					}
				}
			}
			return true
		})
	}
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~SeqE[T], T any](seq S, filter func(T) bool) SeqE[T] {
	return func(yield func(T, error) bool) {
		if seq == nil || filter == nil {
			return
		}
		seq(func(t T, err error) bool {
			if err != nil || filter(t) {
				return yield(t, err)
			}
			return true
		})
	}
}

// Filt creates an erroreable iterator that iterates only those elements for which the 'filter' function returns true.
func Filt[S ~SeqE[T], T any](seq S, filter func(T) (bool, error)) SeqE[T] {
	return func(yield func(T, error) bool) {
		if seq == nil || filter == nil {
			return
		}
		seq(func(t T, err error) bool {
			if err != nil {
				return yield(t, err)
			}
			if ok, err := filter(t); ok || err != nil {
				return yield(t, err)
			}
			return true
		})
	}
}
