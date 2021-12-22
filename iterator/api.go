package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/slice"
)

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	HasNext() bool
	Get() T
}

//Of wararg Iterator constructor based on slice
func Of[T any](elements ...T) Iterator[T] { return NewSlice(elements) }

//ToSlice Iterator to slice converter
func ToSlice[T any](elements Iterator[T]) []T {
	s := make([]T, 0)

	for elements.HasNext() {
		s = append(s, elements.Get())
	}

	return s
}

//Wrap Iterator constructor
func Wrap[T any](elements []T) Iterator[T] {
	return NewSlice(elements)
}

//Iter Iterator constructor
func Iter[T any](elements []T) Iterator[T] {
	return NewSlice(elements)
}

//NewSlice default constructor
func NewSlice[T any](elements []T) *Slice[T] {
	return &Slice[T]{elements: elements}
}

//WrapMap Key, value Iterator constructor.
func WrapMap[K comparable, V any](values map[K]V) Iterator[*KV[K, V]] {
	return NewMap(values)
}

//Map applies 'by' Converter of 'elements' Iterator and accumulate to result Iterator
func Map[From, To any](elements Iterator[From], by conv.Converter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	return &ConvertIter[From, To]{iter: elements, by: by, filters: filters}
}

//Flatt extracts embedded slices of elements by Flatter and accumulate to result Iterator
func Flatt[From, To any](elements Iterator[From], by slice.Flatter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	return &FlattIter[From, To]{iter: elements, by: by, filters: filters}
}

//Filter tests elements and adds that fit to result
func Filter[T any](elements Iterator[T], filters ...check.Predicate[T]) Iterator[T] {
	return &FilterIter[T]{iter: elements, filters: filters}
}

//NotNil filter nullable elements
func NotNil[T any](elements Iterator[T]) Iterator[T] {
	return Filter(elements, check.NotNil[T])
}

func ForEach[T any, It Iterator[T]](elements It, apply func(T), filters ...check.Predicate[T]) {
	if len(filters) == 0 {
		for elements.HasNext() {
			apply(elements.Get())
		}
		return
	}
	for elements.HasNext() {
		v := elements.Get()
		if check.IsFit(v, filters...) {
			apply(v)
		}
	}
}
