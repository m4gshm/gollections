package seq

import (
	"github.com/m4gshm/gollections/c"
	s2 "github.com/m4gshm/gollections/internal/seq2"
)

// Head returns the first key\value pair.
func (s Seq2[K, V]) Head() (K, V, bool) {
	return s2.Head(s)
}

// First returns the first key\value pair that satisfies the condition.
func (s Seq2[K, V]) First(condition func(K, V) bool) (K, V, bool) {
	return s2.First(s, condition)
}

// Firstt returns the first key\value pair that satisfies the condition.
func (s Seq2[K, V]) Firstt(condition func(K, V) (bool, error)) (K, V, bool, error) {
	return s2.Firstt(s, condition)
}

// HasAny checks whether the seq contains a key\value pair that satisfies the condition.
func (s Seq2[K, V]) HasAny(condition func(K, V) bool) bool {
	return s2.HasAny(s, condition)
}

// Union combines several sequences into one.
func (s Seq2[K, V]) Union(seqences ...seq2[K, V]) Seq2[K, V] {
	return s2.Union(append(append(make([]seq2[K, V], len(seqences)+1), s), seqences...)...)
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func (s Seq2[K, V]) Filter(filter func(K, V) bool) Seq2[K, V] {
	return s2.Filter(s, filter)
}

// Filt creates an erroreable iterator that iterates only those key\value pairs for which the 'filter' function returns true.
func (s Seq2[K, V]) Filt(filter func(K, V) (bool, error)) SeqE[c.KV[K, V]] {
	return s2.Filt(s, filter)
}

// Convert creates an iterator that applies the 'converter' function to each iterable key\value pair.
func (s Seq2[K, V]) Convert(converter func(K, V) (K, V)) Seq2[K, V] {
	return s2.Convert(s, converter)
}

// Conv creates an errorable seq that applies the 'converter' function to the iterable key\value pairs.
func (s Seq2[K, V]) Conv(converter func(K, V) (K, V, error)) SeqE[c.KV[K, V]] {
	return s2.Conv(s, converter)
}

// Keys converts a key/value pairs iterator to an iterator of just keys.
func (s Seq2[K, V]) Keys() Seq[K] {
	return s2.Keys(s)
}

// Values converts a key/value pairs iterator to an iterator of just values.
func (s Seq2[K, V]) Values() Seq[V] {
	return s2.Values(s)
}

// FilterKey returns a seq consisting of key/value pairs where the key satisfies the condition of the 'filter' function.
func (s Seq2[K, V]) FilterKey(filter func(K) bool) Seq2[K, V] {
	return s2.FilterKey(s, filter)
}

// FilterValue returns a seq consisting of key/value pairs where the value satisfies the condition of the 'filter' function.
func (s Seq2[K, V]) FilterValue(filter func(V) bool) Seq2[K, V] {
	return s2.FilterValue(s, filter)
}

// ConvertKey returns a seq that applies the 'converter' function to keys.
func (s Seq2[K, V]) ConvertKey(converter func(K) K) Seq2[K, V] {
	return s2.ConvertKey(s, converter)
}

// ConvKey returns a seq that applies the 'converter' function to keys.
func (s Seq2[K, V]) ConvKey(converter func(K) (K, error)) SeqE[c.KV[K, V]] {
	return s2.ConvKey(s, converter)
}

// ConvertValue returns a seq that applies the 'converter' function to values.
func (s Seq2[K, V]) ConvertValue(converter func(V) V) Seq2[K, V] {
	return s2.ConvertValue(s, converter)
}

// ConvValue returns a seq that applies the 'converter' function to values.
func (s Seq2[K, V]) ConvValue(converter func(V) (V, error)) SeqE[c.KV[K, V]] {
	return s2.ConvValue(s, converter)
}

// TrackEach applies the 'consumer' function to the seq key\value pairs.
func (s Seq2[K, V]) TrackEach(consumer func(K, V)) {
	s2.TrackEach(s, consumer)
}
