package it

import (
	"reflect"
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

func NewReflectKV[k comparable, v any](elements map[k]v) *ReflectKV[k, v] {
	refVal := reflect.ValueOf(elements)
	return &ReflectKV[k, v]{elements: elements, refVal: refVal, iter: refVal.MapRange()}
}

func NewOrderedKV[k comparable, v any](order []k, uniques map[k]v) *OrderedKV[k, v] {
	return &OrderedKV[k, v]{elements: New(order), uniques: uniques}
}

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

func (iter *KV[k, v]) HasNext() bool {
	if !iter.hiter.initialized() {
		mapiterinit(iter.maptype, iter.hmap, &iter.hiter)
	} else {
		if mapiterkey(&iter.hiter) == nil {
			return false
		}
		mapiternext(&iter.hiter)
	}
	return mapiterkey(&iter.hiter) != nil
}

func (iter *KV[k, v]) Next() (k, v) {
	iterkey := mapiterkey(&iter.hiter)
	if iterkey == nil {
		var key k
		var value v
		return key, value
	}
	iterelem := mapiterelem(&iter.hiter)
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

type ReflectKV[k comparable, v any] struct {
	elements map[k]v
	iter     *reflect.MapIter
	refVal   reflect.Value
}

var _ c.KVIterator[int, any] = (*ReflectKV[int, any])(nil)

func (iter *ReflectKV[k, v]) HasNext() bool {
	return iter.iter.Next()
}

func (iter *ReflectKV[k, v]) Next() (k, v) {
	key := iter.iter.Key().Interface().(k)
	value := iter.iter.Value().Interface().(v)
	return key, value
}

func (s *ReflectKV[k, v]) Reset() {
	s.iter.Reset(s.refVal)
}

type OrderedKV[k comparable, v any] struct {
	elements *Iter[k]
	uniques  map[k]v
}

var _ c.KVIterator[string, any] = (*OrderedKV[string, any])(nil)

func (s *OrderedKV[k, v]) HasNext() bool {
	return s.elements.HasNext()
}

func (s *OrderedKV[k, v]) Next() (k, v) {
	key := s.elements.Next()
	return key, s.uniques[key]
}

func NewKey[k comparable, v any](uniques map[k]v) *Key[k, v] {
	return &Key[k, v]{KV: NewKV(uniques)}
}

type Key[k comparable, v any] struct {
	*KV[k, v]
}

var _ c.Iterator[string] = (*Key[string, any])(nil)

func (iter *Key[k, v]) Next() k {
	key, _ := iter.KV.Next()
	return key
}

func NewVal[k comparable, v any](uniques map[k]v) *Val[k, v] {
	return &Val[k, v]{KV: NewKV(uniques)}
}

type Val[k comparable, v any] struct {
	*KV[k, v]
}

var _ c.Iterator[any] = (*Val[int, any])(nil)

func (iter *Val[k, v]) Next() v {
	_, val := iter.KV.Next()
	return val
}
