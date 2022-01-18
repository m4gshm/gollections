package it

import (
	"reflect"
	"unsafe"

	"github.com/m4gshm/gollections/typ"
)

func NewReflectKV[k comparable, v any](elements map[k]v) *ReflectKV[k, v] {
	refVal := reflect.ValueOf(elements)
	return &ReflectKV[k, v]{elements: elements, refVal: refVal, iter: refVal.MapRange()}
}

func NewOrderedKV[k comparable, v any](order []*k, uniques map[k]v) *OrderedKV[k, v] {
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
	err     error
}

var _ typ.KVIterator[any, any] = (*KV[any, any])(nil)

func (iter *KV[k, v]) HasNext() bool {
	if !iter.hiter.initialized() {
		mapiterinit(iter.maptype, iter.hmap, &iter.hiter)
	} else {
		if mapiterkey(&iter.hiter) == nil {
			iter.err = Exhausted
			return false
		}
		mapiternext(&iter.hiter)
	}
	return mapiterkey(&iter.hiter) != nil
}

func (iter *KV[k, v]) Get() (k, v, error) {
	if err := iter.err; err != nil {
		var key k
		var value v
		return key, value, err
	}
	iterkey := mapiterkey(&iter.hiter)
	if iterkey == nil {
		err := Exhausted
		iter.err = err
		var key k
		var value v
		return key, value, err
	}
	iterelem := mapiterelem(&iter.hiter)
	var key *k = (*k)(iterkey)
	var value *v = (*v)(iterelem)
	return *key, *value, nil
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
	err      error
}

var _ typ.KVIterator[any, any] = (*ReflectKV[any, any])(nil)

func (iter *ReflectKV[k, v]) HasNext() bool {
	next := iter.iter.Next()
	if !next {
		iter.err = Exhausted
	}
	return next
}

func (iter *ReflectKV[k, v]) Get() (k, v, error) {
	if err := iter.err; err != nil {
		var key k
		var value v
		return key, value, err
	}
	key := iter.iter.Key().Interface().(k)
	value := iter.iter.Value().Interface().(v)
	return key, value, nil
}

func (s *ReflectKV[k, v]) Reset() {
	s.iter.Reset(s.refVal)
	s.err = nil
}

type OrderedKV[k comparable, v any] struct {
	elements *Iter[*k]
	uniques  map[k]v
}

var _ typ.KVIterator[any, any] = (*OrderedKV[any, any])(nil)

func (s *OrderedKV[k, v]) HasNext() bool {
	return s.elements.HasNext()
}

func (s *OrderedKV[k, v]) Get() (k, v, error) {
	ref, err := s.elements.Get()
	if err != nil {
		var key k
		var value v
		return key, value, err
	}
	key := *ref
	return key, s.uniques[key], nil
}
