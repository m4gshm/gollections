package seqe

import (
	"github.com/m4gshm/gollections/predicate/always"
)

type seq[T any] = func(func(T) bool)
type seqE[T any] = seq2[T, error]
type seq2[K, V any] = func(func(K, V) bool)

// Union combines several sequences into one.
func Union[S ~seqE[T], T any](seq ...S) seqE[T] {
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

func Top[S ~seqE[T], T any](n int, seq S) seqE[T] {
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

func Skip[S ~seqE[T], T any](n int, seq S) seqE[T] {
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

func Head[S ~seqE[T], T any](seq S) (v T, ok bool, err error) {
	return First(seq, always.True)
}

func First[S ~seqE[T], T any](seq S, predicate func(T) bool) (v T, ok bool, err error) {
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

func Firstt[S ~seqE[T], T any](seq S, predicate func(T) (bool, error)) (v T, ok bool, err error) {
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

func Slice[S ~seqE[T], T any](seq S) ([]T, error) {
	return SliceCap(seq, 0)
}

func SliceCap[S ~seqE[T], T any](seq S, capacity int) (out []T, e error) {
	if seq == nil {
		return nil, nil
	}
	if capacity > 0 {
		out = make([]T, 0, capacity)
	}
	return Append(seq, out)
}

func Append[S ~seqE[T], T any, TS ~[]T](seq S, out TS) (TS, error) {
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

func Reduce[S ~seqE[T], T any](seq S, merge func(T, T) T) (T, error) {
	result, _, err := ReduceOK(seq, merge)
	return result, err
}

func ReduceOK[S ~seqE[T], T any](seq S, merge func(T, T) T) (result T, ok bool, err error) {
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

func Reducee[S ~seqE[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

func ReduceeOK[S ~seqE[T], T any](seq S, merge func(T, T) (T, error)) (result T, ok bool, err error) {
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

func Accum[T any, S ~seqE[T]](first T, seq S, merge func(T, T) T) (accumulator T, err error) {
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

func Accumm[T any, S ~seqE[T]](first T, seq S, merge func(T, T) (T, error)) (accumulator T, err error) {
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

func HasAny[S ~seqE[T], T any](seq S, predicate func(T) bool) (bool, error) {
	_, ok, err := First(seq, predicate)
	return ok, err
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
func Conv[S ~seqE[From], From, To any](seq S, converter func(From) (To, error)) seqE[To] {
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
func Convert[S ~seqE[From], From, To any](seq S, converter func(From) To) seqE[To] {
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
func Filter[S ~seqE[T], T any](seq S, filter func(T) bool) seqE[T] {
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
func Filt[S ~seqE[T], T any](seq S, filter func(T) (bool, error)) seqE[T] {
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
func ForEach[T any](seq seqE[T], consumer func(T)) error {
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
