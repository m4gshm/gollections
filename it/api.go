package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
)

//Of instantiates Iterator of predefined elements.
func Of[T any](elements ...T) c.Iterator[T] {
	iter := it.NewHead(elements)
	return &iter
}

//Wrap instantiates Iterator using sclie as the elements source.
func Wrap[T any, TS ~[]T](elements TS) c.Iterator[T] {
	iter := it.NewHead(elements)
	return &iter
}

//Map instantiates Iterator that converts elements with a converter and returns them.
func Map[From, To any, IT c.Iterator[From]](elements IT, by c.Converter[From, To]) c.Iterator[To] {
	return it.Map(elements, by)
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any, IT c.Iterator[From]](elements IT, fit c.Predicate[From], by c.Converter[From, To]) c.Iterator[To] {
	return it.MapFit(elements, fit, by)
}

//Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, IT c.Iterator[From]](elements IT, by c.Flatter[From, To]) c.Iterator[To] {
	iter := it.Flatt(elements, by)
	return &iter
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any, IT c.Iterator[From]](elements IT, fit c.Predicate[From], flatt c.Flatter[From, To]) c.Iterator[To] {
	iter := it.FlattFit(elements, fit, flatt)
	return &iter
}

//Filter instantiates Iterator that checks elements by a filter and returns successful ones.
func Filter[T any, IT c.Iterator[T]](elements IT, filter c.Predicate[T]) c.Iterator[T] {
	return it.Filter(elements, filter)
}

//NotNil instantiates Iterator that filters nullable elements.
func NotNil[T any, IT c.Iterator[*T]](elements IT) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

//Reduce reduces elements to an one.
func Reduce[T any, IT c.Iterator[T]](elements IT, by op.Binary[T]) T {
	return it.Reduce(elements, by)
}

//ReduceKV reduces key/value elements to an one.
func ReduceKV[K, V any, IT c.KVIterator[K, V]](elements IT, by op.Quaternary[K, V]) (K, V) {
	return it.ReduceKV(elements, by)
}

//Slice converts an Iterator to a slice.
func Slice[T any](elements c.Iterator[T]) []T {
	return it.Slice[T](elements)
}

//Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable](elements c.Iterator[T], by c.Converter[T, K]) c.MapPipe[K, T, map[K][]T] {
	return it.Group(elements, by)
}

//ForEach applies a walker to elements of an Iterator.
func ForEach[T any, IT c.Iterator[T]](elements IT, walker func(T)) {
	it.ForEach(elements, walker)
}

//ForEachFit applies a walker to elements that satisfy a predicate condition.
func ForEachFit[T any](elements c.Iterator[T], walker func(T), fit c.Predicate[T]) {
	it.ForEachFit(elements, walker, fit)
}
