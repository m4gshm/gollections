package slice

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

//Map creates a lazy Iterator that converts elements with a converter and returns them
func Map[From, To any](elements []From, by conv.Converter[From, To]) *Convert[From, To] {
	return &Convert[From, To]{Elements: elements, By: by}
}

// additionally filters 'From' elements by filters
func MapFit[From, To any](elements []From, fit check.Predicate[From], by conv.Converter[From, To]) *ConvertFit[From, To] {
	return &ConvertFit[From, To]{Elements: elements, By: by, Fit: fit}
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, FS ~[]From](elements FS, by conv.Flatter[From, To]) *Flatten[From, To] {
	return &Flatten[From, To]{Elements: elements, Flatt: by}
}

// additionally checks 'From' elements by fit
func FlattFit[From, To any](elements []From, fit check.Predicate[From], flatt conv.Flatter[From, To]) *FlattenFit[From, To] {
	return &FlattenFit[From, To]{Elements: elements, Flatt: flatt, Fit: fit}
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any, TS []T](elements TS, filter check.Predicate[T]) *Fit[T] {
	return &Fit[T]{Elements: elements, By: filter}
}

//NotNil creates a lazy Iterator that filters nullable elements
func NotNil[T any](elements []*T) *Fit[*T] {
	return Filter(elements, check.NotNil[T])
}
