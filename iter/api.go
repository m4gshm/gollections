package iter

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/typ"
)

//Of Iterator constructor
func Of[T any](elements ...T) typ.Iterator[T] {
	return &Slice[T]{Elements: elements}
}

//Wrap slice to Iterator
func Wrap[T any](elements []T) typ.Iterator[T] {
	return &Slice[T]{Elements: elements}
}

//ToSlice converts Iterator to slice
func ToSlice[T any](elements typ.Iterator[T]) []T {
	s := make([]T, 0)

	for elements.HasNext() {
		s = append(s, elements.Get())
	}

	return s
}

// WrapMap Key, value Iterator constructor.
func WrapMap[K comparable, V any](values map[K]V) typ.Iterator[*KV[K, V]] {
	return NewMap(values)
}

//Map creates a lazy Iterator that converts elements with a converter and returns them
func Map[From, To any](elements typ.Iterator[From], by conv.Converter[From, To]) typ.Iterator[To] {
	return &Convert[From, To]{iter: elements, by: by}
}

// additionally filters 'From' elements by filters
func MapFit[From, To any](elements typ.Iterator[From], fit check.Predicate[From], by conv.Converter[From, To]) typ.Iterator[To] {
	return &ConvertFit[From, To]{iter: elements, by: by, fit: fit}
}

//Reduce reduces elements to an one
func Reduce[T any](elements typ.Iterator[T], by op.Binary[T]) T {
	if !elements.HasNext() {
		var nothing T
		return nothing
	}
	result := elements.Get()
	for elements.HasNext() {
		result = by(result, elements.Get())
	}
	return result
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any](elements typ.Iterator[From], by conv.Flatter[From, To]) typ.Iterator[To] {
	return &Flatten[From, To]{Iter: elements, Flatt: by}
}

// additionally checks 'From' elements by fit
func FlattFit[From, To any](elements typ.Iterator[From], fit check.Predicate[From], flatt conv.Flatter[From, To]) typ.Iterator[To] {
	return &FlattenFit[From, To]{Iter: elements, Flatt: flatt, Fit: fit}
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any](elements typ.Iterator[T], filter check.Predicate[T]) typ.Iterator[T] {
	return &Fit[T]{Iter: elements, By: filter}
}

//NotNil creates a lazy Iterator that filters nullable elements
func NotNil[T any](elements typ.Iterator[T]) typ.Iterator[T] {
	return Filter(elements, check.NotNil[T])
}

//ForEach applies func on elements
func ForEach[T any, It typ.Iterator[T]](elements It, apply func(T)) {
	for elements.HasNext() {
		apply(elements.Get())
	}
}

func ForEachFit[T any, It typ.Iterator[T]](elements It, apply func(T), fit check.Predicate[T]) {
	for elements.HasNext() {
		if v := elements.Get(); fit(v) {
			apply(v)
		}
	}
}
