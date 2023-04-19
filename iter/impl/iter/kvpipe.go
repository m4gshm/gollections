package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
)

// NewKVPipe instantiates Iterator wrapper that converts the elements into key/value pairs and iterates over them.
func NewKVPipe[K comparable, V any, M map[K]V | map[K][]V, IT c.KVIterator[K, V]](it IT, collector MapCollector[K, V, M]) *KVIterPipe[K, V, M] {
	return &KVIterPipe[K, V, M]{KVIterator: it, collector: collector}
}

// KVIterPipe is the key/value Iterator based pipe implementation.
type KVIterPipe[K comparable, V any, M map[K]V | map[K][]V] struct {
	c.KVIterator[K, V]
	collector MapCollector[K, V, M]
}

var (
	_ c.MapPipe[string, any, map[string]any]   = (*KVIterPipe[string, any, map[string]any])(nil)
	_ c.MapPipe[string, any, map[string][]any] = (*KVIterPipe[string, any, map[string][]any])(nil)
)

func (s *KVIterPipe[K, V, M]) FilterKey(predicate func(K) bool) c.MapPipe[K, V, M] {
	return NewKVPipe(FilterKV(s, filter.Key[V](predicate)), s.collector)
}

func (s *KVIterPipe[K, V, M]) ConvertKey(by func(K) K) c.MapPipe[K, V, M] {
	return NewKVPipe(ConvertKV(s, convert.Key[V](by)), s.collector)
}

func (s *KVIterPipe[K, V, M]) FilterValue(predicate func(V) bool) c.MapPipe[K, V, M] {
	return NewKVPipe(FilterKV(s, filter.Value[K](predicate)), s.collector)
}

func (s *KVIterPipe[K, V, M]) ConvertValue(by func(V) V) c.MapPipe[K, V, M] {
	return NewKVPipe(ConvertKV(s, convert.Value[K](by)), s.collector)
}

func (s *KVIterPipe[K, V, M]) Filter(filter func(K, V) bool) c.MapPipe[K, V, M] {
	return NewKVPipe(FilterKV(s, filter), s.collector)
}

func (s *KVIterPipe[K, V, M]) Convert(by func(K, V) (K, V)) c.MapPipe[K, V, M] {
	return NewKVPipe(ConvertKV(s, by), s.collector)
}

func (s *KVIterPipe[K, V, M]) Track(tracker func(K, V) error) error {
	return loop.Track(s.Next, tracker)
}

func (s *KVIterPipe[K, V, M]) TrackEach(tracker func(K, V)) {
	loop.TrackEach(s.Next, tracker)
}

func (s *KVIterPipe[K, V, M]) Reduce(by c.Quaternary[K, V]) (K, V) {
	return loop.ReduceKV(s.Next, by)
}

func (s *KVIterPipe[K, V, M]) Begin() c.KVIterator[K, V] {
	return s
}

func (s *KVIterPipe[K, V, M]) Map() M {
	return s.collector(s)
}

// MapCollector is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result
type MapCollector[K comparable, V any, M map[K]V | map[K][]V] func(c.KVIterator[K, V]) M
