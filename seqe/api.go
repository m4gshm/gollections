package seqe

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/seq"
)

type Seq[T any] = seq.Seq[T]
type SeqE[T any] = seq.SeqE[T]
type Seq2[K, V any] = seq.Seq2[K, V]

func OfIndexed[T any](max int, getAt func(int) (T, error)) Seq2[T, error] {
	if getAt == nil {
		return empty2
	}
	return func(yield func(T, error) bool) {
		for i := range max {
			v, err := getAt(i)
			if ok := yield(v, err); !ok {
				break
			}
		}
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function.
func First[S ~SeqE[T], T any](seq S, predicate func(T) bool) (v T, ok bool, err error) {
	if seq == nil {
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
	if seq == nil {
		return v, false, nil
	}
	seq(func(one T, e error) bool {
		if e != nil {
			err = e
			return false
		} else if p, e := predicate(one); e != nil {
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
func Slice[S ~SeqE[T], T any](seq S) ([]T, error) {
	return SliceCap(seq, 0)
}

// SliceCap collects the elements of the 'seq' sequence into a new slice with predefined capacity.
func SliceCap[S ~SeqE[T], T any](seq S, cap int) (out []T, e error) {
	if seq == nil {
		return nil, nil
	}
	if cap > 0 {
		out = make([]T, 0, cap)
	}
	return Append(seq, out)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func Append[S ~SeqE[T], T any, TS ~[]T](seq S, out TS) (TS, error) {
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

// Reduce reduces the elements of the 'seq' sequence an one using the 'merge' function.
func Reduce[S ~SeqE[T], T any](seq S, merge func(T, T) T) (T, error) {
	result, _, err := ReduceOK(seq, merge)
	return result, err
}

// ReduceOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func ReduceOK[S ~SeqE[T], T any](seq S, merge func(T, T) T) (result T, ok bool, err error) {
	if seq == nil {
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
	if seq == nil {
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
	if seq == nil {
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
	if seq == nil {
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
	if seq == nil {
		return empty2
	}
	return func(yield func(To, error) bool) {
		seq(func(from From, err error) bool {
			var to To
			if err == nil {
				to, err = converter(from)
			}
			return yield(to, err)
		})
	}
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func Convert[S ~SeqE[From], From, To any](seq S, converter func(From) To) SeqE[To] {
	if seq == nil {
		return empty2
	}
	return func(consumer func(To, error) bool) {
		seq(func(from From, err error) bool {
			var to To
			if err == nil {
				to = converter(from)
			}
			return consumer(to, err)
		})
	}
}

func ConvertOK[S ~SeqE[From], From, To any](seq S, converter func(from From) (To, bool)) SeqE[To] {
	if seq == nil {
		return empty2
	}
	return func(consumer func(To, error) bool) {
		seq(func(from From, e error) bool {
			if e != nil {
				var to To
				return consumer(to, e)
			} else if to, ok := converter(from); ok {
				return consumer(to, nil)
			}
			return true
		})
	}
}

func ConvOK[S ~SeqE[From], From, To any](seq S, converter func(from From) (To, bool, error)) SeqE[To] {
	if seq == nil {
		return empty2
	}
	return func(yield func(To, error) bool) {
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
//	for e, err := range seqe.Flat(arrays, slices.Values) {
//	    ...
//	}
func Flat[S ~SeqE[From], STo ~Seq[To], From any, To any](seq S, flattener func(From) STo) SeqE[To] {
	if seq == nil {
		return empty2
	}
	return func(yield func(To, error) bool) {
		seq(func(v From, err error) bool {
			if err != nil {
				var t To
				return yield(t, err)
			}
			elementsTo := flattener(v)
			for e := range elementsTo {
				if !yield(e, err) {
					return false
				}
			}
			return true
		})
	}
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[S ~SeqE[T], T any](seq S, filter func(T) bool) SeqE[T] {
	if seq == nil {
		return empty2
	}
	return func(consumer func(T, error) bool) {
		seq(func(t T, err error) bool {
			if err != nil {
				return consumer(t, err)
			} else if filter(t) {
				return consumer(t, err)
			}
			return true
		})
	}
}

// Filt creates an erroreable iterator that iterates only those elements for which the 'filter' function returns true.
func Filt[S ~SeqE[T], T any](seq S, filter func(T) (bool, error)) SeqE[T] {
	if seq == nil {
		return empty2
	}
	return func(yield func(T, error) bool) {
		seq(func(t T, err error) bool {
			if err != nil {
				return yield(t, err)
			} else if ok, err := filter(t); ok || err != nil {
				return yield(t, err)
			}
			return true
		})
	}
}
