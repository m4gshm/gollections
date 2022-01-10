package it

import (
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

//Map creates a lazy Iterator that converts elements with a converter and returns them
func Map[From, To any, IT typ.Iterator[From]](elements IT, by typ.Converter[From, To]) *Convert[From, To] {
	return &Convert[From, To]{Iter: elements, By: by}
}

// additionally filters 'From' elements by filters
func MapFit[From, To any, IT typ.Iterator[From]](elements IT, fit typ.Predicate[From], by typ.Converter[From, To]) *ConvertFit[From, To] {
	return &ConvertFit[From, To]{Iter: elements, By: by, Fit: fit}
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, IT typ.Iterator[From]](elements IT, by typ.Flatter[From, To]) *Flatten[From, To] {
	return &Flatten[From, To]{Iter: elements, Flatt: by}
}

// additionally checks 'From' elements by fit
func FlattFit[From, To any, IT typ.Iterator[From]](elements IT, fit typ.Predicate[From], flatt typ.Flatter[From, To]) *FlattenFit[From, To] {
	return &FlattenFit[From, To]{Iter: elements, Flatt: flatt, Fit: fit}
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any, IT typ.Iterator[T]](elements IT, filter typ.Predicate[T]) *Fit[T] {
	return &Fit[T]{Iter: elements, By: filter}
}

//NotNil creates a lazy Iterator that filters nullable elements
func NotNil[T any, IT typ.Iterator[*T]](elements IT) *Fit[*T] {
	return Filter(elements, check.NotNil[T])
}

//ForEach applies func on elements
func ForEach[T any, IT typ.Iterator[T]](elements IT, apply func(T)) {
	for elements.HasNext() {
		apply(elements.Get())
	}
}

func ForEachFit[T any, IT typ.Iterator[T]](elements IT, apply func(T), fit typ.Predicate[T]) {
	for elements.HasNext() {
		if v := elements.Get(); fit(v) {
			apply(v)
		}
	}
}

//Reduce reduces elements to an one
func Reduce[T any, IT typ.Iterator[T]](elements IT, by op.Binary[T]) T {
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

//Slice converts Iterator to slice
func Slice[T any, IT typ.Iterator[T]](elements IT) []T {
	s := make([]T, 0)

	for elements.HasNext() {
		s = append(s, elements.Get())
	}

	return s
}
