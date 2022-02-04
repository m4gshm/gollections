package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	impl "github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
)

const NoStarted = impl.NoStarted

//Of - the Iterator constructor.
func Of[T any](elements ...T) c.Iterator[T] {
	return impl.New(elements)
}

//Wrap the same as 'Of' but for a slice.
func Wrap[T any](elements []T) c.Iterator[T] {
	return impl.New(elements)
}

//Pipe returns the Pipe based on iterable elements.
func Pipe[T any](elements c.Iterator[T]) c.Pipe[T, []T, c.Iterator[T]] {
	return impl.NewPipe(elements)
}

//WrapMap returns the KVIterator based on map elements.
func WrapMap[k comparable, v any](elements map[k]v) c.KVIterator[k, v] {
	return impl.NewKV(elements)
}

//Map creates the Iterator that converts elements with a converter and returns them.
func Map[From, To any](elements c.Iterator[From], by c.Converter[From, To]) c.Iterator[To] {
	return impl.Map(elements, by)
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any](elements c.Iterator[From], fit c.Predicate[From], by c.Converter[From, To]) c.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates the Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](elements c.Iterator[From], by c.Flatter[From, To]) c.Iterator[To] {
	return impl.Flatt(elements, by)
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any](elements c.Iterator[From], fit c.Predicate[From], flatt c.Flatter[From, To]) c.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates the Iterator that checks elements by a filter and returns successful ones.
func Filter[T any](elements c.Iterator[T], filter c.Predicate[T]) c.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates the Iterator that filters nullable elements.
func NotNil[T any](elements c.Iterator[*T]) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

//Reduce reduces elements to an one.
func Reduce[T any](elements c.Iterator[T], by op.Binary[T]) T {
	return impl.Reduce(elements, by)
}

func ReduceKV[k, v any](elements c.KVIterator[k, v], by op.Quaternary[k, v]) (k, v) {
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

//ForEach applies func on elements.
func ForEach[T any](elements c.Iterator[T], apply func(T)) {
	impl.ForEach(elements, apply)
}

func ForEachFit[T any](elements c.Iterator[T], apply func(T), fit c.Predicate[T]) {
	impl.ForEachFit(elements, apply, fit)
}
