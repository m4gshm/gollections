// Package iter provides map based iterator implementations
package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
)

// New returns the KVIterator based on map elements
func New[K comparable, V any](elements map[K]V) MapIter[K, V] {
	hmap := *(*unsafe.Pointer)(unsafe.Pointer(&elements))
	i := any(elements)
	maptype := *(*unsafe.Pointer)(unsafe.Pointer(&i))
	var iterator *hiter
	if hmap != nil {
		iterator = new(hiter)
	}
	return MapIter[K, V]{maptype: maptype, hmap: hmap, size: len(elements), iterator: iterator}
}

// MapIter is the embedded map based Iterator implementation
type MapIter[K comparable, V any] struct {
	iterator *hiter
	maptype  unsafe.Pointer
	hmap     unsafe.Pointer
	size     int
}

var _ c.KVIterator[int, any] = (*MapIter[int, any])(nil)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *MapIter[K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(i.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (i *MapIter[K, V]) TrackEach(traker func(key K, value V)) {
	loop.TrackEach(i.Next, traker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *MapIter[K, V]) Next() (key K, value V, ok bool) {
	if i == nil {
		return key, value, false
	}
	iterator := i.iterator
	if iterator == nil {
		return key, value, false
	}
	if !iterator.initialized() {
		mapiterinit(i.maptype, i.hmap, iterator)
	} else {
		mapiternext(iterator)
	}
	iterkey := mapiterkey(iterator)
	if iterkey == nil {
		return key, value, false
	}
	iterelem := mapiterelem(iterator)
	key = *(*K)(iterkey)
	value = *(*V)(iterelem)
	return key, value, true
}

// Cap returns the size of the map
func (i *MapIter[K, V]) Cap() int {
	if i == nil {
		return 0
	}
	return i.size
}

//go:linkname mapiterinit reflect.mapiterinit
func mapiterinit(maptype, hmap unsafe.Pointer, it *hiter)

func mapiterkey(it *hiter) unsafe.Pointer {
	return it.key
}

func mapiterelem(it *hiter) unsafe.Pointer {
	return it.elem
}

//go:linkname mapiternext reflect.mapiternext
func mapiternext(it *hiter)

// hiter's structure matches runtime.hiter's structure
type hiter struct {
	key         unsafe.Pointer
	elem        unsafe.Pointer
	t           unsafe.Pointer
	h           unsafe.Pointer
	buckets     unsafe.Pointer
	bptr        unsafe.Pointer
	overflow    *[]unsafe.Pointer
	oldoverflow *[]unsafe.Pointer
	startBucket uintptr
	offset      uint8
	wrapped     bool
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

func (h *hiter) initialized() bool {
	return h.t != nil
}

// NewKey it the Key constructor.
func NewKey[K comparable, V any](uniques map[K]V) KeyIter[K, V] {
	return KeyIter[K, V]{MapIter: New(uniques)}
}

// KeyIter is the Iterator implementation that provides iterating over keys of a key/value pairs iterator
type KeyIter[K comparable, V any] struct {
	MapIter[K, V]
}

var (
	_ c.Iterator[string] = (*KeyIter[string, any])(nil)
	_ c.Iterator[string] = KeyIter[string, any]{}
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i KeyIter[K, V]) For(walker func(element K) error) error {
	return loop.For(i.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i KeyIter[K, V]) ForEach(walker func(element K)) {
	loop.ForEach(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i KeyIter[K, V]) Next() (K, bool) {
	key, _, ok := i.MapIter.Next()
	return key, ok
}

// Cap returns the iterator capacity
func (i KeyIter[K, V]) Cap() int {
	return i.MapIter.Cap()
}

// NewVal is the Val constructor
func NewVal[K comparable, V any](uniques map[K]V) ValIter[K, V] {
	return ValIter[K, V]{MapIter: New(op.IfElse(uniques != nil, uniques, map[K]V{}))}
}

// ValIter is the Iterator implementation that provides iterating over values of a key/value pairs iterator
type ValIter[K comparable, V any] struct {
	MapIter[K, V]
}

var (
	_ c.Iterator[any] = (*ValIter[int, any])(nil)
	_ c.Iterator[any] = ValIter[int, any]{}
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f ValIter[K, V]) For(walker func(element V) error) error {
	return loop.For(f.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (f ValIter[K, V]) ForEach(walker func(element V)) {
	loop.ForEach(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i ValIter[K, V]) Next() (V, bool) {
	_, val, ok := i.MapIter.Next()
	return val, ok
}

// Cap returns the size of the map
func (i ValIter[K, V]) Cap() int {
	return i.MapIter.Cap()
}
