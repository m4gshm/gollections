// Package seqe provides convering, filtering, and reducing operations for the [seq.SeqE] interface.
package seqe

import (
	"github.com/m4gshm/gollections/predicate/always"
)

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
type SeqE[T any] = func(func(T, error) bool)

// Union combines several sequences into one.
func Union[S ~SeqE[T], T any](seq ...S) SeqE[T] {
	return func(yield func(T, error) bool) {
		for _, s := range seq {
			if s != nil {
				for v, err := range s {
					if !yield(v, err) {
						return
					}
				}
			}
		}
	}
}

// Top returns a sequence of top n elements.
func Top[S ~SeqE[T], T any](n int, seq S) SeqE[T] {
	return func(yield func(T, error) bool) {
		if seq == nil {
			return
		}
		m := n
		seq(func(t T, err error) bool {
			if m == 0 {
				return false
			}
			m--
			return yield(t, err)
		})
	}
}

// Skip returns the seq without first n elements.
func Skip[S ~SeqE[T], T any](n int, seq S) SeqE[T] {
	return func(yield func(T, error) bool) {
		if seq == nil {
			return
		}
		m := n
		seq(func(t T, err error) bool {
			if m == 0 {
				return yield(t, err)
			}
			m--
			return true
		})
	}
}

// While cuts tail elements of the seq that don't match the filter.
func While[S ~SeqE[T], T any](seq S, filter func(T) bool) SeqE[T] {
	return func(yield func(T, error) bool) {
		if seq == nil {
			return
		}
		seq(func(t T, err error) bool {
			if !filter(t) {
				return false
			}
			return yield(t, err)
		})
	}
}

// SkipWhile returns a sequence without first elements of the seq that dont'math the filter.
func SkipWhile[S ~SeqE[T], T any](seq S, filter func(T) bool) SeqE[T] {
	return func(yield func(T, error) bool) {
		if seq == nil {
			return
		}
		started := false
		seq(func(t T, err error) bool {
			if !started && filter(t) {
				return true
			}
			started = true
			return yield(t, err)
		})
	}
}

// Head returns the first element.
func Head[S ~SeqE[T], T any](seq S) (v T, ok bool, err error) {
	return First(seq, always.True)
}

// First returns the first element that satisfies the condition.
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

// Firstt returns the first element that satisfies the condition.
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
func SliceCap[S ~SeqE[T], T any](seq S, capacity int) (out []T, err error) {
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

// Reduce reduces the elements of the seq into one using the 'merge' function.
func Reduce[S ~SeqE[T], T any](seq S, merge func(T, T) T) (T, error) {
	result, _, err := ReduceOK(seq, merge)
	return result, err
}

// ReduceOK reduces the elements of the seq into one using the 'merge' function.
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

// Reducee reduces the elements of the seq into one using the 'merge' function.
func Reducee[S ~SeqE[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

// ReduceeOK reduces the elements of the seq into one using the 'merge' function.
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

// HasAny checks whether the seq contains an element that satisfies the condition.
func HasAny[S ~SeqE[T], T any](seq S, predicate func(T) bool) (bool, error) {
	_, ok, err := First(seq, predicate)
	return ok, err
}

// Conv creates an errorable seq that applies the 'converter' function to the collection elements.
// The error should be checked at every iteration step, like:
//
//	var integers iter.Seq2[int, error]
//	...
//	for s, err := range seqe.Conv(integers, strconv.Itoa) {
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

// ForEach applies the 'consumer' function to the seq elements
func ForEach[T any](seq SeqE[T], consumer func(T)) error {
	if seq == nil {
		return nil
	}
	for v, err := range seq {
		if err != nil {
			return err
		}
		consumer(v)
	}
	return nil
}
