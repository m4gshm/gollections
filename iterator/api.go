package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/slice"
)

type Iterator[T any] interface {
	HasNext() bool
	Get() T
}

//Of wararg Iterator constructor
func Of[T any](elements ...T) Iterator[T] { return NewSlice(elements) }

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

func WrapMap[K comparable, V any](values map[K]V) Iterator[*KV[K, V]] {
	return NewMap(values)
}

//Map applies Converter to items and accumulate to result slice
func Map[From, To any](items Iterator[From], by conv.Converter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	return &ConvertIter[From, To]{iter: items, by: by, filters: filters}

}

//Flatt extracts embedded slices of items by Flatter and accumulate to result slice
func Flatt[From, To any](items Iterator[From], by slice.Flatter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	return &FlattIter[From, To]{iter: items, by: by, filters: filters}
}

//Filter tests items and adds that fit to result
func Filter[T any](items Iterator[T], filters ...check.Predicate[T]) Iterator[T] {
	return &FilterIter[T]{iter: items, filters: filters}
}

//NotNil filter nullable values
func NotNil[T any](items Iterator[T]) Iterator[T] {
	return Filter(items, check.NotNil[T])
}

func ForEach[T any, It Iterator[T]](elements It, apply func(T), filters ...check.Predicate[T]) {
	if len(filters) == 0 {
		for elements.HasNext() {
			apply(elements.Get())
		}
		return
	}
	for elements.HasNext() {
		v:=elements.Get()
		if check.IsFit(v, filters...) {
		apply(v)
		}
	}
}

