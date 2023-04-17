package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
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

func (s *KVIterPipe[K, V, M]) FilterKey(filter func(K) bool) c.MapPipe[K, V, M] {
	kvFit := func(key K, val V) bool { return filter(key) }
	return NewKVPipe(FilterKV(s, kvFit), s.collector)
}

func (s *KVIterPipe[K, V, M]) ConvertKey(by func(K) K) c.MapPipe[K, V, M] {
	return NewKVPipe(ConvertKV(s, func(key K, val V) (K, V) { return by(key), val }), s.collector)
}

func (s *KVIterPipe[K, V, M]) FilterValue(filter func(V) bool) c.MapPipe[K, V, M] {
	kvFit := func(key K, val V) bool { return filter(val) }
	return NewKVPipe(FilterKV(s, kvFit), s.collector)
}

func (s *KVIterPipe[K, V, M]) ConvertValue(by func(V) V) c.MapPipe[K, V, M] {
	return NewKVPipe(ConvertKV(s, func(key K, val V) (K, V) { return key, by(val) }), s.collector)
}

func (s *KVIterPipe[K, V, M]) Filter(filter func(K, V) bool) c.MapPipe[K, V, M] {
	return NewKVPipe(FilterKV(s, filter), s.collector)
}

func (s *KVIterPipe[K, V, M]) Convert(by func(K, V) (K, V)) c.MapPipe[K, V, M] {
	return NewKVPipe(ConvertKV(s, by), s.collector)
}

func (s *KVIterPipe[K, V, M]) Track(tracker func(K, V) error) error {
	for key, val, ok := s.Next(); ok; key, val, ok = s.Next() {
		if err := tracker(key, val); err != nil {
			return err
		}
	}
	return nil
}

func (s *KVIterPipe[K, V, M]) TrackEach(tracker func(K, V)) {

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
