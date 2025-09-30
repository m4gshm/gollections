package seq2

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/map_/resolv"
)

type seq[T any] = func(func(T) bool)
type seqE[T any] = seq2[T, error]
type seq2[K, V any] = func(func(K, V) bool)

// Union combines several sequences into one.
func Union[S ~seq2[K, V], K, V any](seq ...S) seq2[K, V] {
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

// Head returns the first key\value pair.
func Head[S ~seq2[K, V], K, V any](seq S) (k K, v V, ok bool) {
	return First(seq, func(K, V) bool { return true })
}

// HasAny checks whether the seq contains an element that satisfies the condition.
func HasAny[S ~seq2[K, V], K, V any](seq S, condition func(K, V) bool) bool {
	_, _, ok := First(seq, condition)
	return ok
}

// First returns the first key\value pair that satisfies the condition.
func First[S ~seq2[K, V], K, V any](seq S, condition func(K, V) bool) (k K, v V, ok bool) {
	if seq == nil || condition == nil {
		return
	}
	seq(func(oneK K, oneV V) bool {
		if condition(oneK, oneV) {
			k = oneK
			v = oneV
			ok = true
			return false
		}
		return true
	})
	return
}

// Firstt returns the first key\value pair that satisfies the condition.
func Firstt[S ~seq2[K, V], K, V any](seq S, condition func(K, V) (bool, error)) (k K, v V, ok bool, err error) {
	if seq == nil || condition == nil {
		return
	}
	seq(func(oneK K, oneV V) bool {
		ok, err = condition(oneK, oneV)
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

func Filter[S ~seq2[K, V], K, V any](seq S, filter func(K, V) bool) seq2[K, V] {
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

// Filt creates an erroreable iterator that iterates only those key\value pairs for which the 'filter' function returns true.
func Filt[S ~seq2[K, V], K, V any](seq S, filter func(K, V) (bool, error)) seq2[c.KV[K, V], error] {
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

// Convert creates an iterator that applies the 'converter' function to each iterable key\value pair.
func Convert[S ~seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto)) seq2[Kto, Vto] {
	return func(consumer func(Kto, Vto) bool) {
		if seq == nil || converter == nil {
			return
		}
		seq(func(k Kfrom, v Vfrom) bool {
			return consumer(converter(k, v))
		})
	}
}

func Conv[S ~seq2[Kfrom, Vfrom], Kfrom, Vfrom, Kto, Vto any](seq S, converter func(Kfrom, Vfrom) (Kto, Vto, error)) seqE[c.KV[Kto, Vto]] {
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

// Values converts a key/value pairs iterator to an iterator of just values.
func Values[S ~seq2[K, V], K, V any](seq S) seq[V] {
	return func(yield func(V) bool) {
		if seq == nil {
			return
		}
		seq(func(_ K, v V) bool {
			return yield(v)
		})
	}
}

// Keys converts a key/value pairs iterator to an iterator of just keys.
func Keys[S ~seq2[K, V], K, V any](seq S) seq[K] {
	return func(yield func(K) bool) {
		if seq == nil {
			return
		}
		seq(func(k K, _ V) bool {
			return yield(k)
		})
	}
}

// Group collects the elements of the 'seq' sequence into a new map.
func Group[S ~seq2[K, V], K comparable, V any](seq S) map[K][]V {
	return MapResolv(seq, resolv.Slice[K, V])
}

// MapResolv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
func MapResolv[S ~seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) map[K]VR {
	return AppendMapResolv(seq, resolver, nil)
}

// MapResolvOrder collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
// Returns a slice with the keys ordered by the time they were added and the resolved key\value map.
func MapResolvOrder[S ~seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) ([]K, map[K]VR) {
	return AppendMapResolvOrder(seq, resolver, nil, nil)
}

// AppendMapResolv collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values.
func AppendMapResolv[S ~seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR, dest map[K]VR) map[K]VR {
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

// AppendMapResolvOrder collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values
// Additionaly populates the 'order' slice by the keys ordered by the time they were added and the resolved key\value map.
func AppendMapResolvOrder[S ~seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR, order []K, dest map[K]VR) ([]K, map[K]VR) {
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

// TrackEach applies the 'consumer' function to the seq elements.
func TrackEach[S ~seq2[K, V], K, V any](seq S, consumer func(K, V)) {
	if seq == nil {
		return
	}
	for k, v := range seq {
		consumer(k, v)
	}
}
