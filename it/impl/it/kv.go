package it

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

// ToKVIter converts an iterator of c.KV elements
func ToKVIter[K, V any, IT c.Iterator[c.KV[K, V]]](elements IT) KVIter[K, V] {
	return KVIter[K, V]{elements: elements}
}

// NewOrderedKV is the OrderedKV constructor
func NewOrderedKV[K comparable, V any](uniques map[K]V, elements ArrayIter[K]) OrderedKV[K, V] {
	return OrderedKV[K, V]{elements: elements, uniques: uniques}
}

// NewKV returns the KVIterator based on map elements
func NewKV[K comparable, V any](elements map[K]V) EmbedMapKV[K, V] {
	m := elements
	hmap := *(*unsafe.Pointer)(unsafe.Pointer(&m))
	i := any(m)
	maptype := *(*unsafe.Pointer)(unsafe.Pointer(&i))

	return EmbedMapKV[K, V]{maptype: maptype, hmap: hmap, size: len(elements), iter: new(hiter)}
}

// EmbedMapKV is the embedded map based Iterator implementation
type EmbedMapKV[K comparable, V any] struct {
	iter    *hiter
	maptype unsafe.Pointer
	hmap    unsafe.Pointer
	size    int
}

var _ c.KVIterator[int, any] = (*EmbedMapKV[int, any])(nil)

func (i *EmbedMapKV[K, V]) Next() (K, V, bool) {
	if !i.iter.initialized() {
		mapiterinit(i.maptype, i.hmap, i.iter)
	} else {
		mapiternext(i.iter)
	}
	iterkey := mapiterkey(i.iter)
	if iterkey == nil {
		var key K
		var value V
		return key, value, false
	}
	iterelem := mapiterelem(i.iter)
	var key *K = (*K)(iterkey)
	var value *V = (*V)(iterelem)
	return *key, *value, true
}

func (i *EmbedMapKV[K, V]) Cap() int {
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

// OrderedKV is the ordered key/value pairs Iterator implementation
type OrderedKV[K comparable, V any] struct {
	elements ArrayIter[K]
	uniques  map[K]V
}

var _ c.KVIterator[string, any] = (*OrderedKV[string, any])(nil)

func (i *OrderedKV[K, V]) Next() (K, V, bool) {
	if key, ok := i.elements.Next(); ok {
		return key, i.uniques[key], true
	}
	var k K
	var v V
	return k, v, false
}

func (i *OrderedKV[K, V]) Cap() int {
	return i.elements.Cap()
}

// KVIter is the wrapper of an iterator of c.KV elements
type KVIter[K, V any] struct {
	elements c.Iterator[c.KV[K, V]]
}

var _ c.KVIterator[string, any] = (*KVIter[string, any])(nil)

func (i *KVIter[K, V]) Next() (K, V, bool) {
	if kv, ok := i.elements.Next(); ok {
		return kv.K, kv.V, true
	}
	var k K
	var v V
	return k, v, false
}

// NewKey it the Key constructor.
func NewKey[K comparable, V any](uniques map[K]V) Key[K, V] {
	return Key[K, V]{EmbedMapKV: NewKV(uniques)}
}

// Key is the Iterator implementation that provides iterating over keys of a key/value pairs iterator
type Key[K comparable, V any] struct {
	EmbedMapKV[K, V]
}

var (
	_ c.Iterator[string] = (*Key[string, any])(nil)
	_ c.Iterator[string] = Key[string, any]{}
)

func (i Key[K, V]) Next() (K, bool) {
	key, _, ok := i.EmbedMapKV.Next()
	return key, ok
}

func (i Key[K, V]) Cap() int {
	return i.EmbedMapKV.Cap()
}

// NewVal is the Val constructor
func NewVal[K comparable, V any](uniques map[K]V) Val[K, V] {
	return Val[K, V]{EmbedMapKV: NewKV(uniques)}
}

// Val is the Iterator implementation that provides iterating over values of a key/value pairs iterator
type Val[K comparable, V any] struct {
	EmbedMapKV[K, V]
}

var _ c.Iterator[any] = (*Val[int, any])(nil)

func (i Val[K, V]) Next() (V, bool) {
	_, val, ok := i.EmbedMapKV.Next()
	return val, ok
}

func (i *Val[K, V]) Cap() int {
	return i.EmbedMapKV.Cap()
}
