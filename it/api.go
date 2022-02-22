package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	impl "github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
)

//Of creates the Iterator of predefined elements.
func Of[T any](elements ...T) *impl.Iter[T] {
	iter := impl.NewHead(elements)
	return &iter
}

//Wrap creates the Iterator using sclie as the elements source.
func Wrap[T any, TS ~[]T](elements TS) *impl.Iter[T] {
	iter := impl.NewHead(elements)
	return &iter
}

//Map creates the Iterator that converts elements with a converter and returns them.
func Map[From, To any, IT c.Iterator[From]](elements IT, by c.Converter[From, To]) c.Iterator[To] {
	return impl.Map(elements, by)
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any, IT c.Iterator[From]](elements IT, fit c.Predicate[From], by c.Converter[From, To]) c.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates the Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, IT c.Iterator[From]](elements IT, by c.Flatter[From, To]) c.Iterator[To] {
	return impl.Flatt(elements, by)
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any, IT c.Iterator[From]](elements IT, fit c.Predicate[From], flatt c.Flatter[From, To]) c.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates the Iterator that checks elements by a filter and returns successful ones.
func Filter[T any, IT c.Iterator[T]](elements IT, filter c.Predicate[T]) impl.Fit[T, IT] {
	return impl.Filter(elements, filter)
}

//NotNil creates the Iterator that filters nullable elements.
func NotNil[T any, IT c.Iterator[*T]](elements IT) impl.Fit[*T, IT] {
	return Filter(elements, check.NotNil[T])
}

//Reduce reduces elements to an one.
func Reduce[T any, IT c.Iterator[T]](elements IT, by op.Binary[T]) T {
	return impl.Reduce(elements, by)
}

//ReduceKV reduces key/value elements to an one.
func ReduceKV[K, V any, IT c.KVIterator[K, V]](elements IT, by op.Quaternary[K, V]) (K, V) {
	return impl.ReduceKV(elements, by)
}

//Slice converts an Iterator to a slice.
func Slice[T any](elements c.Iterator[T]) []T {
	return impl.Slice[T](elements)
}

//Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable](elements c.Iterator[T], by c.Converter[T, K]) c.MapPipe[K, T, map[K][]T] {
	return impl.Group(elements, by)
}

//ForEach applies a walker to elements of an Iterator.
func ForEach[T any, IT c.Iterator[T]](elements IT, walker func(T)) {
	impl.ForEach(elements, walker)
}

//ForEachFit applies a walker to elements that satisfy a predicate condition.
func ForEachFit[T any](elements c.Iterator[T], walker func(T), fit c.Predicate[T]) {
	impl.ForEachFit(elements, walker, fit)
}
