package seq

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/internal/seq2"
	convertkv "github.com/m4gshm/gollections/kv/convert"
	filterkv "github.com/m4gshm/gollections/kv/predicate"
)

func (s Seq2[K, V]) Filter(filter func(K, V) bool) Seq2[K, V] {
	return seq2.Filter(s, filter)
}

func (s Seq2[K, V]) Filt(predicate func(K, V) (bool, error)) SeqE[c.KV[K, V]] {
	return seq2.Filt(s, predicate)
}

func (s Seq2[K, V]) Convert(converter func(K, V) (K, V)) Seq2[K, V] {
	return seq2.Convert(s, converter)
}

func (s Seq2[K, V]) Conv(converter func(K, V) (K, V, error)) SeqE[c.KV[K, V]] {
	return seq2.Conv(s, converter)
}

func (s Seq2[K, V]) FilterKey(predicate func(K) bool) Seq2[K, V] {
	return s.Filter(filterkv.Key[V](predicate))
}

func (s Seq2[K, V]) FilterValue(predicate func(V) bool) Seq2[K, V] {
	return s.Filter(filterkv.Value[K](predicate))
}

func (s Seq2[K, V]) ConvertKey(converter func(K) K) Seq2[K, V] {
	return seq2.Convert(s, convertkv.Key[V](converter))
}

func (s Seq2[K, V]) ConvertValue(converter func(V) V) Seq2[K, V] {
	return seq2.Convert(s, convertkv.Value[K](converter))
}

func (s Seq2[K, V]) TrackEach(consumer func(K, V)) {
	seq2.TrackEach(s, consumer)
}
