// Package seq2 extends [iter.Seq2] API with convering, filtering, and reducing functionality.
package seq2

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/map_/resolv"
)

// Seq is an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[T any] = func(func(T) bool)

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
type SeqE[T any] = Seq2[T, error]

// Seq2 is an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] = func(func(K, V) bool)

// Union combines several sequences into one.
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

// Head returns the first key\value pair.
func Head[S ~Seq2[K, V], K, V any](seq S) (k K, v V, ok bool) {
	return First(seq, func(K, V) bool { return true })
}

// HasAny checks whether the seq contains a key\value pair that satisfies the condition.
func HasAny[S ~Seq2[K, V], K, V any](seq S, condition func(K, V) bool) bool {
	_, _, ok := First(seq, condition)
	return ok
}

// First returns the first key\value pair that satisfies the condition.
func First[S ~Seq2[K, V], K, V any](seq S, condition func(K, V) bool) (k K, v V, ok bool) {
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
func Firstt[S ~Seq2[K, V], K, V any](seq S, condition func(K, V) (bool, error)) (k K, v V, ok bool, err error) {
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

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
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

// Filt creates an erroreable iterator that iterates only those key\value pairs for which the 'filter' function returns true.
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

// Convert creates an iterator that applies the 'converter' function to each iterable key\value pair.
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

// Conv creates an errorable seq that applies the 'converter' function to the iterable key\value pairs.
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

// Values converts a key/value pairs iterator to an iterator of just values.
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

// Keys converts a key/value pairs iterator to an iterator of just keys.
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

// Group collects the elements of the 'seq' sequence into a new map.
func Group[S ~Seq2[K, V], K comparable, V any](seq S) map[K][]V {
	return MapResolv(seq, resolv.Slice[K, V])
}

// MapResolv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
func MapResolv[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) map[K]VR {
	return AppendMapResolv(seq, resolver, nil)
}

// MapResolvOrder collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
// Returns a slice with the keys ordered by the time they were added and the resolved key\value map.
func MapResolvOrder[S ~Seq2[K, V], K comparable, V, VR any](seq S, resolver func(exists bool, key K, valResolv VR, val V) VR) ([]K, map[K]VR) {
	return AppendMapResolvOrder(seq, resolver, nil, nil)
}

// AppendMapResolv collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values.
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

// AppendMapResolvOrder collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values
// Additionaly populates the 'order' slice by the keys ordered by the time they were added and the resolved key\value map.
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

// TrackEach applies the 'consumer' function to the seq key\value pairs.
func TrackEach[S ~Seq2[K, V], K, V any](seq S, consumer func(K, V)) {
	if seq == nil {
		return
	}
	for k, v := range seq {
		consumer(k, v)
	}
}
