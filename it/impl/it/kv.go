package it

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

func NewOrderedKV[k comparable, v any](order []k, uniques map[k]v) *OrderedKV[k, v] {
	return &OrderedKV[k, v]{elements: New(order), uniques: uniques}
}

//NewKV returns the KVIterator based on map elements.
func NewKV[k comparable, v any](elements map[k]v) *KV[k, v] {
	m := elements
	hmap := *(*unsafe.Pointer)(unsafe.Pointer(&m))
	i := any(m)
	maptype := *(*unsafe.Pointer)(unsafe.Pointer(&i))

	return &KV[k, v]{maptype: maptype, hmap: hmap}
}

type KV[k comparable, v any] struct {
	hiter
	maptype unsafe.Pointer
	hmap    unsafe.Pointer
}

var _ c.KVIterator[int, any] = (*KV[int, any])(nil)

func (i *KV[k, v]) HasNext() bool {
	if !i.hiter.initialized() {
		mapiterinit(i.maptype, i.hmap, &i.hiter)
	} else {
		if mapiterkey(&i.hiter) == nil {
			return false
		}
		mapiternext(&i.hiter)
	}
	return mapiterkey(&i.hiter) != nil
}

func (i *KV[k, v]) Next() (k, v) {
	iterkey := mapiterkey(&i.hiter)
	if iterkey == nil {
		var key k
		var value v
		return key, value
	}
	iterelem := mapiterelem(&i.hiter)
	var key *k = (*k)(iterkey)
	var value *v = (*v)(iterelem)
	return *key, *value
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

// hiter's structure matches runtime.hiter's structure.
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

type OrderedKV[k comparable, v any] struct {
	elements *Iter[k]
	uniques  map[k]v
}

var _ c.KVIterator[string, any] = (*OrderedKV[string, any])(nil)

func (i *OrderedKV[k, v]) HasNext() bool {
	return i.elements.HasNext()
}

func (i *OrderedKV[k, v]) Next() (k, v) {
	key := i.elements.Next()
	return key, i.uniques[key]
}

func NewKey[k comparable, v any](uniques map[k]v) *Key[k, v] {
	return &Key[k, v]{KV: NewKV(uniques)}
}

type Key[k comparable, v any] struct {
	*KV[k, v]
}

var _ c.Iterator[string] = (*Key[string, any])(nil)

func (i *Key[k, v]) Next() k {
	key, _ := i.KV.Next()
	return key
}

func NewVal[k comparable, v any](uniques map[k]v) *Val[k, v] {
	return &Val[k, v]{KV: NewKV(uniques)}
}

type Val[k comparable, v any] struct {
	*KV[k, v]
}

var _ c.Iterator[any] = (*Val[int, any])(nil)

func (i *Val[k, v]) Next() v {
	_, val := i.KV.Next()
	return val
}
