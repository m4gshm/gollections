package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	m "github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

//ConvertKVsToMap converts a slice of key/value pairs to the Map.
func ConvertKVsToMap[K comparable, V any](elements []*c.KV[K, V]) *Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for _, kv := range elements {
		uniques[kv.Key()] = kv.Value()
	}
	return WrapMap(uniques)
}

//NewMap creates the Map with values copied from an map.
func NewMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

//WrapMap creates the Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	return &Map[K, V]{elements: elements}
}

//Map is the Collection implementation that provides element access by an unique key.
type Map[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ c.Map[int, any] = (*Map[int, any])(nil)
	_ fmt.Stringer    = (*Map[int, any])(nil)
)

//Begin creates a key/value iterator interface.
func (s *Map[K, V]) Begin() c.KVIterator[K, V] {
	iter := s.Head()
	return &iter
}

//Head creates a key/value iterator implementation started from the head.
func (s *Map[K, V]) Head() it.KV[K, V] {
	return it.NewKV(s.elements)
}

//Collect exports the content as a map.
func (s *Map[K, V]) Collect() map[K]V {
	out := make(map[K]V, len(s.elements))
	for key, val := range s.elements {
		out[key] = val
	}
	return out
}

//Sort transforms to the ordered Map contains sorted elements.
func (s *Map[K, V]) Sort(less func(k1, k2 K) bool) *ordered.Map[K, V] {
	return ordered.WrapMap(slice.Sort(map_.Keys(s.elements), less), s.elements)
}

//String is part of the Stringer interface for printing the string representation of this Map.
func (s *Map[K, V]) String() string {
	return m.ToString(s.elements)
}

//Track apply a tracker to touch key, value from the inside. To stop traking just return the m.Break.
func (s *Map[K, V]) Track(tracker func(K, V) error) error {
	return m.Track(s.elements, tracker)
}

//Track apply a tracker to touch each key, value from the inside.
func (s *Map[K, V]) TrackEach(tracker func(K, V)) {
	m.TrackEach(s.elements, tracker)
}

func (s *Map[K, V]) For(walker func(*c.KV[K, V]) error) error {
	return m.For(s.elements, walker)
}

func (s *Map[K, V]) ForEach(walker func(*c.KV[K, V])) {
	m.ForEach(s.elements, walker)
}

func (s *Map[K, V]) Len() int {
	return len(s.elements)
}

func (s *Map[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Map[K, V]) Contains(key K) bool {
	_, ok := s.elements[key]
	return ok
}

func (s *Map[K, V]) Get(key K) (V, bool) {
	val, ok := s.elements[key]
	return val, ok
}

func (s *Map[K, V]) Keys() c.Collection[K, []K, c.Iterator[K]] {
	return s.K()
}

func (s *Map[K, V]) K() *MapKeys[K, V] {
	return WrapKeys(s.elements)
}

func (s *Map[K, V]) Values() c.Collection[V, []V, c.Iterator[V]] {
	return s.V()
}

func (s *Map[K, V]) V() *MapValues[K, V] {
	return WrapVal(s.elements)
}

func (s *Map[K, V]) FilterKey(fit c.Predicate[K]) c.MapPipe[K, V, map[K]V] {
	var kvFit c.BiPredicate[K, V] = func(key K, val V) bool { return fit(key) }
	iter := s.Head()
	return it.NewKVPipe(it.FilterKV(&iter, kvFit), collect.Map[K, V])
}

func (s *Map[K, V]) MapKey(by c.Converter[K, K]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.MapKV(&iter, func(key K, val V) (K, V) { return by(key), val }), collect.Map[K, V])
}

func (s *Map[K, V]) FilterValue(fit c.Predicate[V]) c.MapPipe[K, V, map[K]V] {
	var kvFit c.BiPredicate[K, V] = func(key K, val V) bool { return fit(val) }
	iter := s.Head()
	return it.NewKVPipe(it.FilterKV(&iter, kvFit), collect.Map[K, V])
}

func (s *Map[K, V]) MapValue(by c.Converter[V, V]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.MapKV(&iter, func(key K, val V) (K, V) { return key, by(val) }), collect.Map[K, V])
}

func (s *Map[K, V]) Filter(filter c.BiPredicate[K, V]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.FilterKV(&iter, filter), collect.Map[K, V])
}

func (s *Map[K, V]) Map(by c.BiConverter[K, V, K, V]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.MapKV(&iter, by), collect.Map[K, V])
}

func (s *Map[K, V]) Reduce(by op.Quaternary[K, V]) (K, V) {
	iter := s.Head()
	return it.ReduceKV(&iter, by)
}
