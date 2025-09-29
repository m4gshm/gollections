package seq

import (
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/always"
	"golang.org/x/exp/constraints"
)

type Seq[T any] = func(yield func(T) bool)

type SeqE[T any] = Seq2[T, error]

type Seq2[K, V any] = func(yield func(K, V) bool)

func Of[T any](elements ...T) Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range elements {
			if !yield(v) {
				break
			}
		}
	}
}

func Union[S ~Seq[T], T any](seq ...S) Seq[T] {
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

func OfNextGet[T any](hasNext func() bool, getNext func() T) Seq[T] {
	return func(yield func(T) bool) {
		for hasNext() {
			if o := getNext(); !yield(o) {
				return
			}
		}
	}
}

func OfNext[T any](hasNext func() bool, pushNext func(*T)) Seq[T] {
	return OfNextGet(hasNext, func() (o T) { pushNext(&o); return o })
}

func OfSourceNextGet[S, T any](source S, hasNext func(S) bool, getNext func(S) T) Seq[T] {
	return OfNextGet(func() bool { return hasNext(source) }, func() T { return getNext(source) })
}

func OfSourceNext[S, T any](source S, hasNext func(S) bool, pushNext func(S, *T)) Seq[T] {
	return OfNext(func() bool { return hasNext(source) }, func(next *T) { pushNext(source, next) })
}

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

func ToSeq2[S ~Seq[T], T, K, V any](seq S, converter func(T) (K, V)) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(v T) bool {
			return yield(converter(v))
		})
	}
}

func Top[S ~Seq[T], T any](n int, seq S) Seq[T] {
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

func Skip[S ~Seq[T], T any](n int, seq S) Seq[T] {
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

func While[S ~Seq[T], T any](seq S, predicate func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		seq(func(t T) bool {
			if !predicate(t) {
				return false
			}
			return yield(t)
		})
	}
}

func SkipWhile[S ~Seq[T], T any](seq S, predicate func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil {
			return
		}
		started := false
		seq(func(t T) bool {
			if !started && predicate(t) {
				return true
			}
			started = true
			return yield(t)
		})
	}
}

func Head[S ~Seq[T], T any](seq S) (v T, ok bool) {
	return First(seq, always.True)
}

func First[S ~Seq[T], T any](seq S, predicate func(T) bool) (v T, ok bool) {
	if seq == nil || predicate == nil {
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

func Firstt[S ~Seq[T], T any](seq S, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	if seq == nil || predicate == nil {
		return v, false, nil
	}
	seq(func(one T) bool {
		ok, err = predicate(one)
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

func Slice[S ~Seq[T], T any](seq S) []T {
	return SliceCap(seq, 0)
}

func SliceCap[S ~Seq[T], T any](seq S, capacity int) (out []T) {
	if seq == nil {
		return nil
	}
	if capacity > 0 {
		out = make([]T, 0, capacity)
	}
	return Append(seq, out)
}

func Append[S ~Seq[T], TS ~[]T, T any](seq S, out TS) TS {
	if seq == nil {
		return out
	}
	seq(func(v T) bool {
		out = append(out, v)
		return true
	})
	return out
}

func Reduce[S ~Seq[T], T any](seq S, merge func(T, T) T) T {
	result, _ := ReduceOK(seq, merge)
	return result
}

func ReduceOK[S ~Seq[T], T any](seq S, merge func(T, T) T) (result T, ok bool) {
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

func Reducee[S ~Seq[T], T any](seq S, merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

func ReduceeOK[S ~Seq[T], T any](seq S, merge func(T, T) (T, error)) (result T, ok bool, err error) {
	if seq == nil || merge == nil {
		return result, false, nil
	}
	started := false
	seq(func(v T) bool {
		if !started {
			result = v
			started = true
			return true
		} else {
			result, err = merge(result, v)
			return err == nil
		}

	})
	return result, started, err
}

func Accum[T any, S ~Seq[T]](first T, seq S, merge func(T, T) T) T {
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

func Accumm[T any, S ~Seq[T]](first T, seq S, merge func(T, T) (T, error)) (accumulator T, err error) {
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

func Sum[S ~Seq[T], T op.Summable](seq S) (out T) {
	return Accum(out, seq, op.Sum[T])
}

func HasAny[S ~Seq[T], T any](seq S, predicate func(T) bool) bool {
	_, ok := First(seq, predicate)
	return ok
}

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

func Conv[S ~Seq[From], From, To any](seq S, converter func(From) (To, error)) SeqE[To] {
	return SeqE[To](ToSeq2(seq, converter))
}

func Convert[S ~Seq[From], From, To any](seq S, converter func(From) To) Seq[To] {
	return func(yield func(To) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(from From) bool {
			return yield(converter(from))
		})
	}
}

func ConvertOK[S ~Seq[From], From, To any](seq S, converter func(from From) (To, bool)) Seq[To] {
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

func ConvOK[S ~Seq[From], From, To any](seq S, converter func(from From) (To, bool, error)) SeqE[To] {
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

func Flat[S ~Seq[From], STo ~[]To, From any, To any](seq S, flattener func(From) STo) Seq[To] {
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

func FlatSeq[S ~Seq[From], STo ~Seq[To], From any, To any](seq S, flattener func(From) STo) Seq[To] {
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

func Flatt[S ~Seq[From], STo ~[]To, From any, To any](seq S, flattener func(From) (STo, error)) SeqE[To] {
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

func FlattSeq[S ~Seq[From], STo ~SeqE[To], From any, To any](seq S, flattener func(From) STo) SeqE[To] {
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

func Filter[S ~Seq[T], T any](seq S, filter func(T) bool) Seq[T] {
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

func Filt[S ~Seq[T], T any](seq S, filter func(T) (bool, error)) SeqE[T] {
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

func ToKV[S ~Seq[T], T, K, V any](seq S, keyExtractor func(T) K, valExtractor func(T) V) Seq2[K, V] {
	return ToSeq2(seq, func(t T) (K, V) { return keyExtractor(t), valExtractor(t) })
}

func KeyValues[S ~Seq[T], T, K, V any](seq S, keyExtractor func(T) K, valsExtractor func(T) []V) Seq2[K, V] {
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

func ForEach[T any](seq Seq[T], consumer func(T)) {
	if seq == nil {
		return
	}
	for v := range seq {
		consumer(v)
	}
}
