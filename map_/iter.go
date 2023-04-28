// Package iter provides map based iterator implementations
package map_

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
)

// NewIter returns the MapIter based on map elements
func NewIter[K comparable, V any](elements map[K]V) Iter[K, V] {
	hmap := *(*unsafe.Pointer)(unsafe.Pointer(&elements))
	i := any(elements)
	maptype := *(*unsafe.Pointer)(unsafe.Pointer(&i))
	var iterator *hiter
	if hmap != nil {
		iterator = new(hiter)
	}
	return Iter[K, V]{maptype: maptype, hmap: hmap, size: len(elements), iterator: iterator}
}

// Iter is the embedded map based Iterator implementation
type Iter[K comparable, V any] struct {
	iterator *hiter
	maptype  unsafe.Pointer
	hmap     unsafe.Pointer
	size     int
}

var _ c.KVIterator[int, any] = (*Iter[int, any])(nil)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *Iter[K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(i.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (i *Iter[K, V]) TrackEach(traker func(key K, value V)) {
	loop.TrackEach(i.Next, traker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *Iter[K, V]) Next() (key K, value V, ok bool) {
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
func (i *Iter[K, V]) Cap() int {
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

// NewKeyIter instantiates a map keys iterator
func NewKeyIter[K comparable, V any](uniques map[K]V) KeyIter[K, V] {
	return KeyIter[K, V]{Iter: NewIter(uniques)}
}

// KeyIter is the Iterator implementation that provides iterating over keys of a key/value pairs iterator
type KeyIter[K comparable, V any] struct {
	Iter[K, V]
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
	key, _, ok := i.Iter.Next()
	return key, ok
}

// Cap returns the iterator capacity
func (i KeyIter[K, V]) Cap() int {
	return i.Iter.Cap()
}

// NewVal is the Val constructor
func NewVal[K comparable, V any](uniques map[K]V) NewValIter[K, V] {
	return NewValIter[K, V]{Iter: NewIter(op.IfElse(uniques != nil, uniques, map[K]V{}))}
}

// NewKeyIter instantiates a map values iterator
type NewValIter[K comparable, V any] struct {
	Iter[K, V]
}

var (
	_ c.Iterator[any] = (*NewValIter[int, any])(nil)
	_ c.Iterator[any] = NewValIter[int, any]{}
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f NewValIter[K, V]) For(walker func(element V) error) error {
	return loop.For(f.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (f NewValIter[K, V]) ForEach(walker func(element V)) {
	loop.ForEach(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i NewValIter[K, V]) Next() (V, bool) {
	_, val, ok := i.Iter.Next()
	return val, ok
}

// Cap returns the size of the map
func (i NewValIter[K, V]) Cap() int {
	return i.Iter.Cap()
}
