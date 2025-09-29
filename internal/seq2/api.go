package seq2

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/internal/seq"
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/op"
	"golang.org/x/exp/constraints"
)

type Seq[T any] = seq.Seq[T]

type SeqE[T any] = Seq2[T, error]

type Seq2[K, V any] = seq.Seq2[K, V]

func Of[T any](elements ...T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range elements {
			if !yield(i, v) {
				break
			}
		}
	}
}

func OfMap[K comparable, V any](elements map[K]V) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range elements {
			if !yield(k, v) {
				break
			}
		}
	}
}

func Union[S ~Seq2[K, V], K, V any](seq ...S) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, s := range seq {
			if s != nil {
				for k, v := range s {
					if !yield(k, v) {
						return
					}
				}
			}
		}
	}
}

func OfIndexed[T any](amount int, getAt func(int) T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if getAt == nil {
			return
		}
		for i := range amount {
			if !yield(i, getAt(i)) {
				break
			}
		}
	}
}

func OfIndexedKV[K, V any](amount int, getAt func(int) (K, V)) Seq2[K, V] {
	return func(yield func(K, V) bool) {
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

func OfIndexedPair[K, V any](amount int, getKey func(int) K, getValue func(int) V) Seq2[K, V] {
	return OfIndexedKV(amount, func(i int) (K, V) { return getKey(i), getValue(i) })
}

func Series[T any](first T, next func(int, T) (T, bool)) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if next == nil {
			return
		}
		i := 0
		current := first
		if !yield(i, current) {
			return
		}
		for {
			i++
			next, ok := next(i, current)
			if !ok {
				break
			}
			if !yield(i, next) {
				break
			}
			current = next
		}
	}
}

func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) Seq2[int, T] {
	amount := toInclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	amount++
	return func(yield func(int, T) bool) {
		e := from
		for i := 0; i < int(amount); i++ {
			if !yield(i, e) {
				return
			}
			e = e + delta
		}
	}
}

func Range[T constraints.Integer | rune](from T, toExclusive T) Seq2[int, T] {
	amount := toExclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	return func(yield func(int, T) bool) {
		e := from
		for i := 0; i < int(amount); i++ {
			if !yield(i, e) {
				return
			}
			e = e + delta
		}
	}
}

func ToSeq[S ~Seq2[K, V], T, K, V any](seq S, converter func(K, V) T) Seq[T] {
	return func(yield func(T) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(k K, v V) bool {
			return yield(converter(k, v))
		})
	}
}

func Top[S ~Seq2[K, V], K, V any](n int, seq S) S {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		m := n
		seq(func(k K, v V) bool {
			if m == 0 {
				return false
			}
			m--
			return yield(k, v)
		})
	}
}

func Skip[S ~Seq2[K, V], K, V any](n int, seq S) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		m := n
		seq(func(k K, v V) bool {
			if m == 0 {
				return yield(k, v)
			}
			m--
			return true
		})
	}
}

func While[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		seq(func(k K, v V) bool {
			if !filter(k, v) {
				return false
			}
			return yield(k, v)
		})
	}
}

func SkipWhile[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil {
			return
		}
		started := false
		seq(func(k K, v V) bool {
			if !started && filter(k, v) {
				return true
			}
			started = true
			return yield(k, v)
		})
	}
}

func Head[S ~Seq2[K, V], K, V any](seq S) (k K, v V, ok bool) {
	return First(seq, func(K, V) bool { return true })
}

func HasAny[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) bool {
	_, _, ok := First(seq, filter)
	return ok
}

func First[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) (k K, v V, ok bool) {
	if seq == nil || filter == nil {
		return
	}
	seq(func(oneK K, oneV V) bool {
		if filter(oneK, oneV) {
			k = oneK
			v = oneV
			ok = true
			return false
		}
		return true
	})
	return
}

func Firstt[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) (bool, error)) (k K, v V, ok bool, err error) {
	if seq == nil || filter == nil {
		return
	}
	seq(func(oneK K, oneV V) bool {
		ok, err = filter(oneK, oneV)
		if ok {
			k = oneK
			v = oneV
			return false
		} else if err != nil {
			return false
		}
		return true
	})
	return k, v, ok, err
}

func Reduce[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) T) T {
	result, _ := ReduceOK(seq, merge)
	return result
}

func ReduceOK[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) T) (result T, ok bool) {
	if seq == nil || merge == nil {
		return result, false
	}
	started := false
	seq(func(k K, v V) bool {
		result = merge(op.IfElse(!started, nil, &result), k, v)
		started = true
		return true
	})
	return result, started
}

func Reducee[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) (T, error)) (T, error) {
	result, _, err := ReduceeOK(seq, merge)
	return result, err
}

func ReduceeOK[S ~Seq2[K, V], K, V, T any](seq S, merge func(prev *T, k K, v V) (T, error)) (result T, ok bool, err error) {
	if seq == nil || merge == nil {
		return result, false, nil
	}
	started := false
	seq(func(k K, v V) bool {
		result, err = merge(op.IfElse(!started, nil, &result), k, v)
		started = true
		return err == nil
	})
	return result, started, err
}

func Filter[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) bool) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if seq == nil || filter == nil {
			return
		}
		seq(func(k K, v V) bool {
			if filter(k, v) {
				return yield(k, v)
			}
			return true
		})
	}
}

func Filt[S ~Seq2[K, V], K, V any](seq S, filter func(K, V) (bool, error)) Seq2[c.KV[K, V], error] {
	return func(yield func(c.KV[K, V], error) bool) {
		if seq == nil || filter == nil {
			return
		}
		seq(func(k K, v V) bool {
			if ok, err := filter(k, v); ok || err != nil {
				return yield(kv.New(k, v), err)
			}
			return true
		})
	}
}

func Convert[S ~Seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto)) Seq2[Kto, Vto] {
	return func(consumer func(Kto, Vto) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(k Kfrom, v Vfrom) bool {
			return consumer(converter(k, v))
		})
	}
}

func Conv[S ~Seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto, error)) SeqE[c.KV[Kto, Vto]] {
	return func(consumer func(c.KV[Kto, Vto], error) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(k Kfrom, v Vfrom) bool {
			kto, vto, err := converter(k, v)
			return consumer(kv.New(kto, vto), err)
		})
	}
}

func Values[S ~Seq2[K, V], K, V any](seq S) Seq[V] {
	return func(yield func(V) bool) {
		if seq == nil {
			return
		}
		seq(func(_ K, v V) bool {
			return yield(v)
		})
	}
}

func Keys[S ~Seq2[K, V], K, V any](seq S) Seq[K] {
	return func(yield func(K) bool) {
		if seq == nil {
			return
		}
		seq(func(k K, _ V) bool {
			return yield(k)
		})
	}
}

func Group[S ~Seq2[K, V], K comparable, V any](seq S) map[K][]V {
	return MapResolv(seq, resolv.Slice[K, V])
}

func Map[S ~Seq2[K, V], K comparable, V any](seq S) map[K]V {
	return MapResolv(seq, resolv.First[K, V])
}

func MapResolv[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) map[K]VR {
	return AppendMapResolv(seq, resolver, nil)
}

func MapResolvOrder[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) ([]K, map[K]VR) {
	return AppendMapResolvOrder(seq, resolver, nil, nil)
}

func AppendMapResolv[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR, dest map[K]VR) map[K]VR {
	if seq == nil || resolver == nil {
		return nil
	}
	if dest == nil {
		dest = map[K]VR{}
	}
	seq(func(k K, v V) bool {
		exists, ok := dest[k]
		dest[k] = resolver(ok, k, exists, v)
		return true
	})
	return dest
}

func AppendMapResolvOrder[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR, order []K, dest map[K]VR) ([]K, map[K]VR) {
	if seq == nil || resolver == nil {
		return nil, nil
	}
	if dest == nil {
		dest = map[K]VR{}
	}
	seq(func(k K, v V) bool {
		exists, ok := dest[k]
		dest[k] = resolver(ok, k, exists, v)
		if !ok {
			order = append(order, k)
		}
		return true
	})
	return order, dest
}

func TrackEach[S ~Seq2[K, V], K, V any](seq S, consumer func(K, V)) {
	if seq == nil {
		return
	}
	for k, v := range seq {
		consumer(k, v)
	}
}
