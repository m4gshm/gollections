package slice

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	impl "github.com/m4gshm/gollections/it/impl/slice"
)

//Map creates the Iterator that converts elements with a converter and returns them.
func Map[From, To any, FS ~[]From](elements FS, by c.Converter[From, To]) c.Iterator[To] {
	return impl.Map(elements, by)
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], by c.Converter[From, To]) c.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates the Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, FS ~[]From](elements FS, by c.Flatter[From, To]) c.Iterator[To] {
	return impl.Flatt(elements, by)
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], flatt c.Flatter[From, To]) c.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates the Iterator that checks elements by filters and returns successful ones.
func Filter[T any, TS ~[]T](elements TS, filter c.Predicate[T]) c.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates the Iterator that filters nullable elements.
func NotNil[T any, TRS ~[]*T](elements TRS) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}
