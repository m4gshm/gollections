package seqe

import (
	"github.com/m4gshm/gollections/internal/seq"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/always"
)

type Seq[T any] = seq.Seq[T]

type SeqE[T any] = seq.SeqE[T]

type Seq2[K, V any] = seq.Seq2[K, V]

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

func OfNextGet[T any](hasNext func() bool, getNext func() (T, error)) SeqE[T] {
	return func(yield func(T, error) bool) {
		for hasNext() {
			if o, err := getNext(); !yield(o, err) {
				return
			}
		}
	}
}

func OfNext[T any](hasNext func() bool, pushNext func(*T) error) SeqE[T] {
	return OfNextGet(hasNext, func() (o T, err error) { return o, pushNext(&o) })
}

func OfSourceNextGet[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) SeqE[T] {
	return OfNextGet(func() bool { return hasNext(source) }, func() (T, error) { return getNext(source) })
}

func OfSourceNext[S, T any](source S, hasNext func(S) bool, pushNext func(S, *T) error) SeqE[T] {
	return OfNext(func() bool { return hasNext(source) }, func(next *T) error { return pushNext(source, next) })
}

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

func Head[S ~SeqE[T], T any](seq S) (v T, ok bool, err error) {
	return First(seq, always.True)
}

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

func Slice[S ~SeqE[T], T any](seq S) ([]T, error) {
	return SliceCap(seq, 0)
}

func SliceCap[S ~SeqE[T], T any](seq S, capacity int) (out []T, e error) {
	if seq == nil {
		return nil, nil
	}
	if capacity > 0 {
		out = make([]T, 0, capacity)
	}
	return Append(seq, out)
}

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

func Reduce[S ~SeqE[T], T any](seq S, merge func(T, T) T) (T, error) {
	result, _, err := ReduceOK(seq, merge)
	return result, err
}

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

func Reducee[S ~SeqE[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

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

func Sum[S ~SeqE[T], T op.Summable](seq S) (out T, err error) {
	return Accum(out, seq, op.Sum[T])
}

func HasAny[S ~SeqE[T], T any](seq S, predicate func(T) bool) (bool, error) {
	_, ok, err := First(seq, predicate)
	return ok, err
}

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

func ConvOK[S ~SeqE[From], From, To any](seq S, converter func(from From) (To, bool, error)) SeqE[To] {
	return func(yield func(To, error) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(from From, err error) bool {
			if err != nil {
				var to To
				return yield(to, err)
			} else if to, ok, err := converter(from); ok || err != nil {
				return yield(to, err)
			}
			return true
		})
	}
}

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

func Group[S ~SeqE[T], T any, K comparable, V any](seq S, keyExtractor func(T) K, valExtractor func(T) V) (map[K][]V, error) {
	groups := map[K][]V{}
	for e, err := range seq {
		if err != nil {
			return groups, err
		}
		key := keyExtractor(e)
		group := groups[key]
		if group == nil {
			group = make([]V, 0)
		}
		groups[key] = append(group, valExtractor((e)))
	}
	return groups, nil
}

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
