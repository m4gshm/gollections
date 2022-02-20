package slice

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
)

//Map creates the Iterator that converts elements with a converter and returns them.
func Map[From, To any, FS ~[]From](elements FS, by c.Converter[From, To]) *Convert[From, To] {
	return &Convert[From, To]{elements: elements, by: by}
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], by c.Converter[From, To]) *ConvertFit[From, To] {
	return &ConvertFit[From, To]{elements: elements, by: by, Fit: fit}
}

//Flatt creates the Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, FS ~[]From](elements FS, by c.Flatter[From, To]) *Flatten[From, To] {
	return &Flatten[From, To]{elements: elements, flatt: by}
}

//FlattFit additionally filters –'From' elements.
func FlattFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], flatt c.Flatter[From, To]) *FlattenFit[From, To] {
	return &FlattenFit[From, To]{elements: elements, flatt: flatt, fit: fit}
}

//Filter creates the Iterator that checks elements by filters and returns successful ones.
func Filter[T any, TS ~[]T](elements TS, filter c.Predicate[T]) *Fit[T] {
	return &Fit[T]{elements: elements, by: filter}
}

//NotNil creates the Iterator that filters nullable elements.
func NotNil[T any, TRS ~[]*T](elements TRS) *Fit[*T] {
	return Filter(elements, check.NotNil[T])
}
