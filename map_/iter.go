package map_

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv/collection"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
)

// NewIter returns the Iter based on map elements
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

var _ collection.Iterator[int, any] = (*Iter[int, any])(nil)

// All is used to iterate through the iterator using `for ... range`.
func (i *Iter[K, V]) All(consumer func(key K, value V) bool) {
	kvloop.All(i.Next, consumer)
}

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning Break
func (i *Iter[K, V]) Track(traker func(key K, value V) error) error {
	return kvloop.Track(i.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (i *Iter[K, V]) TrackEach(traker func(key K, value V)) {
	kvloop.TrackEach(i.Next, traker)
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

// Size returns the size of the map
func (i *Iter[K, V]) Size() int {
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

// All is used to iterate through the iterator using `for ... range`.
func (i KeyIter[K, V]) All(consumer func(element K) bool) {
	loop.All(i.Next, consumer)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning Break
func (i KeyIter[K, V]) For(consumer func(element K) error) error {
	return loop.For(i.Next, consumer)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i KeyIter[K, V]) ForEach(consumer func(element K)) {
	loop.ForEach(i.Next, consumer)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i KeyIter[K, V]) Next() (K, bool) {
	key, _, ok := i.Iter.Next()
	return key, ok
}

// Size returns the iterator capacity
func (i KeyIter[K, V]) Size() int {
	return i.Iter.Size()
}

// NewValIter is the main values iterator constructor
func NewValIter[K comparable, V any](uniques map[K]V) ValIter[K, V] {
	return ValIter[K, V]{Iter: NewIter(op.IfElse(uniques != nil, uniques, map[K]V{}))}
}

// ValIter is a map values iterator
type ValIter[K comparable, V any] struct {
	Iter[K, V]
}

var (
	_ c.Iterator[any] = (*ValIter[int, any])(nil)
	_ c.Iterator[any] = ValIter[int, any]{}
)

// All is used to iterate through the iterator using `for ... range`.
func (i ValIter[K, V]) All(consumer func(element V) bool) {
	loop.All(i.Next, consumer)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning Break
func (i ValIter[K, V]) For(consumer func(element V) error) error {
	return loop.For(i.Next, consumer)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i ValIter[K, V]) ForEach(consumer func(element V)) {
	loop.ForEach(i.Next, consumer)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i ValIter[K, V]) Next() (V, bool) {
	_, val, ok := i.Iter.Next()
	return val, ok
}

// Size returns the size of the map
func (i ValIter[K, V]) Size() int {
	return i.Iter.Size()
}
