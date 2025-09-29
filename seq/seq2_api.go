package seq

import (
	"github.com/m4gshm/gollections/c"
	s2 "github.com/m4gshm/gollections/internal/seq2"
	"github.com/m4gshm/gollections/kv/convert"
	"github.com/m4gshm/gollections/kv/predicate"
)

func (s Seq2[K, V]) Head() (K, V, bool) {
	return s2.Head(s)
}

func (s Seq2[K, V]) First(filter func(K, V) bool) (K, V, bool) {
	return s2.First(s, filter)
}

func (s Seq2[K, V]) Firstt(filter func(K, V) (bool, error)) (K, V, bool, error) {
	return s2.Firstt(s, filter)
}

func (s Seq2[K, V]) HasAny(filter func(K, V) bool) bool {
	return s2.HasAny(s, filter)
}

func (s Seq2[K, V]) Union(seqences ...seq2[K, V]) Seq2[K, V] {
	return s2.Union(append(append(make([]seq2[K, V], len(seqences)+1), s), seqences...)...)
}

func (s Seq2[K, V]) Filter(filter func(K, V) bool) Seq2[K, V] {
	return s2.Filter(s, filter)
}

func (s Seq2[K, V]) Filt(predicate func(K, V) (bool, error)) SeqE[c.KV[K, V]] {
	return s2.Filt(s, predicate)
}

func (s Seq2[K, V]) Convert(converter func(K, V) (K, V)) Seq2[K, V] {
	return s2.Convert(s, converter)
}

func (s Seq2[K, V]) Conv(converter func(K, V) (K, V, error)) SeqE[c.KV[K, V]] {
	return s2.Conv(s, converter)
}

func (s Seq2[K, V]) Keys() Seq[K] {
	return s2.Keys(s)
}

func (s Seq2[K, V]) Values() Seq[V] {
	return s2.Values(s)
}

func (s Seq2[K, V]) FilterKey(filter func(K) bool) Seq2[K, V] {
	return s.Filter(predicate.Key[V](filter))
}

func (s Seq2[K, V]) FilterValue(filter func(V) bool) Seq2[K, V] {
	return s.Filter(predicate.Value[K](filter))
}

func (s Seq2[K, V]) ConvertKey(converter func(K) K) Seq2[K, V] {
	return s2.Convert(s, convert.Key[V](converter))
}

func (s Seq2[K, V]) ConvertValue(converter func(V) V) Seq2[K, V] {
	return s2.Convert(s, convert.Value[K](converter))
}

func (s Seq2[K, V]) TrackEach(consumer func(K, V)) {
	s2.TrackEach(s, consumer)
}
