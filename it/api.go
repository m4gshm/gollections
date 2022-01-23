package it

import (
	"github.com/m4gshm/gollections/check"
	impl "github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

const NoStarted = impl.NoStarted

var Exhausted = impl.Exhausted

//Of - the Iterator constructor.
func Of[T any](elements ...T) typ.Iterator[T] {
	return impl.NewReseteable(elements)
}

//Wrap the same as 'Of' but for a slice.
func Wrap[T any](elements []T) typ.Iterator[T] {
	return impl.NewReseteable(elements)
}

//Pipe returns the Pipe based on iterable elements.
func Pipe[T any](elements typ.Iterator[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return impl.NewPipe(elements)
}

//WrapMap returns the KVIterator based on map elements.
func WrapMap[k comparable, v any](elements map[k]v) typ.KVIterator[k, v] {
	return impl.NewKV(elements)
}

//Map creates the Iterator that converts elements with a converter and returns them.
func Map[From, To any](elements typ.Iterator[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements, by)
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any](elements typ.Iterator[From], fit typ.Predicate[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates the Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](elements typ.Iterator[From], by typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements, by)
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any](elements typ.Iterator[From], fit typ.Predicate[From], flatt typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates the Iterator that checks elements by a filter and returns successful ones.
func Filter[T any](elements typ.Iterator[T], filter typ.Predicate[T]) typ.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates the Iterator that filters nullable elements.
func NotNil[T any](elements typ.Iterator[*T]) typ.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

//Reduce reduces elements to an one.
func Reduce[T any](elements typ.Iterator[T], by op.Binary[T]) T {
	return impl.Reduce(elements, by)
}

func ReduceKV[k, v any](elements typ.KVIterator[k, v], by op.Quaternary[k, v]) (k, v) {
	return impl.ReduceKV(elements, by)
}

//Slice converts an Iterator to a slice.
func Slice[T any](elements typ.Iterator[T]) []T {
	return impl.Slice[T](elements)
}

//Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable](elements typ.Iterator[T], by typ.Converter[T, K]) typ.MapPipe[K, T, map[K][]T] {
	return impl.Group(elements, by)
}

//For applies func on elements.
func For[T any](elements typ.Iterator[T], apply func(T)) error {
	return impl.For(elements, apply)
}

func ForFit[T any](elements typ.Iterator[T], apply func(T), fit typ.Predicate[T]) error {
	return impl.ForFit(elements, apply, fit)
}
