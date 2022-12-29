package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/kvit"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
)

// ConvertKVsToMap converts a slice of key/value pairs to the Map.
func ConvertKVsToMap[K comparable, V any](elements []c.KV[K, V]) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for _, kv := range elements {
		uniques[kv.Key()] = kv.Value()
	}
	return WrapMap(uniques)
}

// NewMap instantiates Map with values copied from an map.
func NewMap[K comparable, V any](elements map[K]V) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// WrapMap instantiates Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) Map[K, V] {
	return Map[K, V]{elements: elements}
}

// Map is the Collection implementation that provides element access by an unique key.
type Map[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ c.Map[int, any] = (*Map[int, any])(nil)
	_ c.Map[int, any] = Map[int, any]{}
	_ fmt.Stringer    = (*Map[int, any])(nil)
	_ fmt.Stringer    = Map[int, any]{}
)

// Begin creates a key/value iterator interface.
func (m Map[K, V]) Begin() c.KVIterator[K, V] {
	return ptr.Of(m.Head())
}

// Head creates a key/value iterator implementation started from the head.
func (m Map[K, V]) Head() it.EmbedMapKVIter[K, V] {
	return it.NewEmbedMapKV(m.elements)
}

// Collect exports the content as a map.
func (m Map[K, V]) Collect() map[K]V {
	return map_.Clone(m.elements)
}

// Sort transforms to the ordered Map contains sorted elements.
func (m Map[K, V]) Sort(less slice.Less[K]) ordered.Map[K, V] {
	return m.sortBy(sort.Slice, less)
}

func (m Map[K, V]) StableSort(less slice.Less[K]) ordered.Map[K, V] {
	return m.sortBy(sort.SliceStable, less)
}

func (m Map[K, V]) sortBy(sorter slice.Sorter, less slice.Less[K]) ordered.Map[K, V] {
	c := map_.Keys(m.elements)
	slice.Sort(map_.Keys(m.elements), sorter, less)
	return ordered.WrapMap(c, m.elements)
}

// String is part of the Stringer interface for printing the string representation of this Map.
func (m Map[K, V]) String() string {
	return map_.ToString(m.elements)
}

// Track apply a tracker to touch key, value from the inside. To stop traking just return the map_.Break.
func (m Map[K, V]) Track(tracker func(K, V) error) error {
	return map_.Track(m.elements, tracker)
}

// TrackEach apply a tracker to touch each key, value from the inside.
func (m Map[K, V]) TrackEach(tracker func(K, V)) {
	map_.TrackEach(m.elements, tracker)
}

func (m Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	return map_.For(m.elements, walker)
}

func (m Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	map_.ForEach(m.elements, walker)
}

func (m Map[K, V]) Len() int {
	return len(m.elements)
}

func (m Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m Map[K, V]) Contains(key K) bool {
	_, ok := m.elements[key]
	return ok
}

func (m Map[K, V]) Get(key K) (V, bool) {
	val, ok := m.elements[key]
	return val, ok
}

func (m Map[K, V]) Keys() c.Collection[K, []K, c.Iterator[K]] {
	return m.K()
}

func (m Map[K, V]) K() MapKeys[K, V] {
	return WrapKeys(m.elements)
}

func (m Map[K, V]) Values() c.Collection[V, []V, c.Iterator[V]] {
	return m.V()
}

func (m Map[K, V]) V() MapValues[K, V] {
	return WrapVal(m.elements)
}

func (m Map[K, V]) FilterKey(fit c.Predicate[K]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.FilterKV(ptr.Of(m.Head()), c.FitKey[K, V](fit)), kvit.ToMap[K, V])
}

func (m Map[K, V]) MapKey(by c.Converter[K, K]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.MapKV(ptr.Of(m.Head()), func(key K, val V) (K, V) { return by(key), val }), kvit.ToMap[K, V])
}

func (m Map[K, V]) FilterValue(fit c.Predicate[V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.FilterKV(ptr.Of(m.Head()), c.FitValue[K](fit)), kvit.ToMap[K, V])
}

func (m Map[K, V]) MapValue(by c.Converter[V, V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.MapKV(ptr.Of(m.Head()), func(key K, val V) (K, V) { return key, by(val) }), kvit.ToMap[K, V])
}

func (m Map[K, V]) Filter(filter c.BiPredicate[K, V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.FilterKV(ptr.Of(m.Head()), filter), kvit.ToMap[K, V])
}

// Map creates an Iterator that applies a converter to iterable elements.
// The 'by' converter is applied only when the 'Next' method of returned Iterator is called and does not change the elements of the map.
func (m Map[K, V]) Map(by c.BiConverter[K, V, K, V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.MapKV(ptr.Of(m.Head()), by), kvit.ToMap[K, V])
}

// Reduce reduces key\value pairs to an one.
func (m Map[K, V]) Reduce(by c.Quaternary[K, V]) (K, V) {
	return it.ReduceKV(ptr.Of(m.Head()), by)
}
