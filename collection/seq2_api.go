package collection

import (
	convertkv "github.com/m4gshm/gollections/kv/convert"
	filterkv "github.com/m4gshm/gollections/kv/predicate"
	filter "github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/seq2"
)

// func (s Seq2[K, V]) Slice() []T {
// 	return seq.Slice(s)
// }

// func (s Seq2[K, V]) Reduce(merge func(a T, b T) T) T {
// 	return seq.Reduce(s, merge)
// }

func (s Seq2[K, V]) TrackEach(consumer func(K, V)) {
	seq2.TrackEach(s, consumer)
}

func (s Seq2[K, V]) Filter(filter func(K, V) bool) Seq2[K, V] {
	return seq2.Filter(s, filter)
}

// func (s Seq2[K, V]) Filt(predicate func(K, V) (bool, error)) Seq2[c.KV[K, V], error] {
// 	return seq2.Filt(s, predicate)
// }

// func (s Seq2[K, V]) HasAny(predicate func(T) bool) bool {
// 	return seq.HasAny(s, predicate)
// }

func (s Seq2[K, V]) Convert(converter func(K, V) (K, V)) Seq2[K, V] {
	return seq2.Convert(s, converter)
}

// func (s Seq2[K, V]) Append(out []T) []T {
// 	return seq.Append(s, out)
// }

// func (s Seq2[K, V]) ForEach(f func(T)) {
// 	seq.ForEach(s, f)
// }

func (s Seq2[K, V]) FilterKey(predicate filter.Predicate[K]) Seq2[K, V] {
	return s.Filter(filterkv.Key[V](predicate))
}

func (s Seq2[K, V]) FilterValue(predicate filter.Predicate[V]) Seq2[K, V] {
	return s.Filter(filterkv.Value[K](predicate))
}

func (s Seq2[K, V]) ConvertKey(converter func(K) K) Seq2[K, V] {
	return seq2.Convert(s, convertkv.Key[V](converter))
}

func (s Seq2[K, V]) ConvertValue(converter func(V) V) Seq2[K, V] {
	return seq2.Convert(s, convertkv.Value[K](converter))
}
