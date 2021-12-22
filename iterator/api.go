package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

type Iterator[T any] interface {
	HasNext() bool
	Get() T
	Next() (T, bool)
}

//Of wararg Iterator constructor
func Of[T any](elements ...T) Iterator[T] { return Wrap(elements) }

func SliceOf[T any](elements Iterator[T]) []T {
	s := make([]T, 0)

	for v, ok := elements.Next(); ok; v, ok = elements.Next() {
		s = append(s, v)
	}

	return s
}

//Wrap Iterator constructor
func Wrap[T any](elements []T) Iterator[T] {
	return It(elements)
}

func Iter[T any](elements []T) Iterator[T] {
	s:= &Slice[T]{elements: elements}
	s.Iterator = s
	return s
}

func WrapMap[K comparable, V any](values map[K]V) Iterator[*KV[K, V]] {
	return NewMap(values)
}

//Map applies Converter to items and accumulate to result slice
func Map[From, To any](items Iterator[From], by conv.Converter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	return &ConvertIter[From, To]{iter: items, by: by, filters: filters}

}

//Flatt extracts embedded slices of items by Flatter and accumulate to result slice
func Flatt[From, To any](items Iterator[From], by Flatter[From, To, Iterator[To]], filters ...check.Predicate[From]) Iterator[To] {
	return &FlattIter[From, To]{iter: items, by: by, filters: filters}
}

//Filter tests items and adds that fit to result
func Filter[T any](items Iterator[T], filters ...check.Predicate[T]) Iterator[T] {
	return &FilterIter[T]{iter: items, filters: filters}
}

//NotNil filter nullable values
func NotNil[T any](items Iterator[T]) Iterator[T] {
	var (
		result = make([]T, 0)
		add    = func(v T) { result = append(result, v) }
	)
	for v, ok := items.Next(); ok; v, ok = items.Next() {
		if check.NotNil(v) {
			add(v)
		}
	}
	return Wrap(result)
}

func ForEach[T any, It Iterator[T]](elements It, apply func(T), filters ...check.Predicate[T]) {
	if len(filters) == 0 {
		for v, ok := elements.Next(); ok; v, ok = elements.Next() {
			apply(v)
		}
		return
	}
	for v, ok := elements.Next(); ok; v, ok = elements.Next() {
		if check.IsFit(v, filters...) {
			apply(v)
		}
	}
}

//Flatter extracts Iterator of To
type Flatter[From, To any, ToIt Iterator[To]] conv.Converter[From, ToIt]
