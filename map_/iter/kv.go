package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice/iter"
)

// NewOrdered is the OrderedKV constructor
func NewOrdered[K comparable, V any](uniques map[K]V, elements iter.ArrayIter[K]) OrderedEmbedMapKVIter[K, V] {
	return OrderedEmbedMapKVIter[K, V]{elements: elements, uniques: uniques}
}

// New returns the KVIterator based on map elements
func New[K comparable, V any](elements map[K]V) EmbedMapKVIter[K, V] {
	hmap := *(*unsafe.Pointer)(unsafe.Pointer(&elements))
	i := any(elements)
	maptype := *(*unsafe.Pointer)(unsafe.Pointer(&i))
	var iterator *hiter
	if hmap != nil {
		iterator = new(hiter)
	}
	return EmbedMapKVIter[K, V]{maptype: maptype, hmap: hmap, size: len(elements), iterator: iterator}
}

// EmbedMapKVIter is the embedded map based Iterator implementation
type EmbedMapKVIter[K comparable, V any] struct {
	iterator *hiter
	maptype  unsafe.Pointer
	hmap     unsafe.Pointer
	size     int
}

var _ c.KVIterator[int, any] = (*EmbedMapKVIter[int, any])(nil)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *EmbedMapKVIter[K, V]) Next() (key K, value V, ok bool) {
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
func (i *EmbedMapKVIter[K, V]) Cap() int {
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

// OrderedEmbedMapKVIter is the ordered key/value pairs Iterator implementation
type OrderedEmbedMapKVIter[K comparable, V any] struct {
	elements iter.ArrayIter[K]
	uniques  map[K]V
}

var _ c.KVIterator[string, any] = (*OrderedEmbedMapKVIter[string, any])(nil)

// Next returns the next key/value pair.
// The ok result indicates whether the pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *OrderedEmbedMapKVIter[K, V]) Next() (key K, val V, ok bool) {
	if i != nil {
		if key, ok = i.elements.Next(); ok {
			val = i.uniques[key]
		}
	}
	return key, val, ok
}

// Cap returns the iterator capacity
func (i *OrderedEmbedMapKVIter[K, V]) Cap() int {
	return i.elements.Cap()
}

// NewKey it the Key constructor.
func NewKey[K comparable, V any](uniques map[K]V) Key[K, V] {
	return Key[K, V]{EmbedMapKVIter: New(uniques)}
}

// Key is the Iterator implementation that provides iterating over keys of a key/value pairs iterator
type Key[K comparable, V any] struct {
	EmbedMapKVIter[K, V]
}

var (
	_ c.Iterator[string] = (*Key[string, any])(nil)
	_ c.Iterator[string] = Key[string, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i Key[K, V]) Next() (K, bool) {
	key, _, ok := i.EmbedMapKVIter.Next()
	return key, ok
}

// Cap returns the iterator capacity
func (i Key[K, V]) Cap() int {
	return i.EmbedMapKVIter.Cap()
}

// NewVal is the Val constructor
func NewVal[K comparable, V any](uniques map[K]V) Val[K, V] {
	return Val[K, V]{EmbedMapKVIter: New(op.IfElse(uniques != nil, uniques, map[K]V{}))}
}

// Val is the Iterator implementation that provides iterating over values of a key/value pairs iterator
type Val[K comparable, V any] struct {
	EmbedMapKVIter[K, V]
}

var (
	_ c.Iterator[any] = (*Val[int, any])(nil)
	_ c.Iterator[any] = Val[int, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i Val[K, V]) Next() (V, bool) {
	_, val, ok := i.EmbedMapKVIter.Next()
	return val, ok
}

// Cap returns the size of the map
func (i Val[K, V]) Cap() int {
	return i.EmbedMapKVIter.Cap()
}
