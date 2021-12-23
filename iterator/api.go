package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/slice"
)

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element
	HasNext() bool
	//retrieves next element
	Get() T
}

//Of wararg Iterator constructor based on slice
func Of[T any](elements ...T) Iterator[T] { return NewSlice(elements) }

//ToSlice converts Iterator to slice
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

// WrapMap Key, value Iterator constructor.
func WrapMap[K comparable, V any](values map[K]V) Iterator[*KV[K, V]] {
	return NewMap(values)
}

//Map creates a lazy Iterator that converts elements with a converter and returns them
// additionally filters 'From' elements by filters
func Map[From, To any](elements Iterator[From], by conv.Converter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	if len(filters) > 0 {
		return &ConvertFilterIter[From, To]{iter: elements, by: by, filters: filters}
	}
	return &ConvertIter[From, To]{iter: elements, by: by}
}

//Reduce reduces elements to an one
func Reduce[T any](elements Iterator[T], by conv.BinaryOp[T]) T {
	first := true
	var result T
	for elements.HasNext() {
		if first {
			result = elements.Get()
			first = false
		} else {
			result = by(result, elements.Get())
		}
	}
	return result
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
// additionally filters 'From' elements by filters
func Flatt[From, To any](elements Iterator[From], by slice.Flatter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	if len(filters) > 0 {
		return &FlattFilterIter[From, To]{iter: elements, by: by, filters: filters}
	}
	return &FlattIter[From, To]{iter: elements, by: by}
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any](elements Iterator[T], filters ...check.Predicate[T]) Iterator[T] {
	return &FilterIter[T]{iter: elements, filters: filters}
}

//NotNil reates a lazy Iterator that filters nullable elements
func NotNil[T any](elements Iterator[T]) Iterator[T] {
	return Filter(elements, check.NotNil[T])
}

//ForEach applies func on elements
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
