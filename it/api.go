package it

import (
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/it/impl/it"
	impl "github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

const NoStarted = it.NoStarted

var Exhausted = it.Exhausted

//Of Iterator constructor
func Of[T any](elements ...T) typ.Iterator[T] {
	return impl.NewReseteable(elements)
}

//Wrap slice to Iterator
func Wrap[T any](elements []T) typ.Iterator[T] {
	return impl.NewReseteable(elements)
}

//Pipe new stream
func Pipe[T any](elements typ.Iterator[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return impl.NewPipe(elements)
}

// WrapMap Key, value Iterator constructor.
func WrapMap[k comparable, v any](elements map[k]v) typ.Iterator[*typ.KV[k, v]] {
	return impl.NewKV(elements)
}

//Map creates a lazy Iterator that converts elements with a converter and returns them
func Map[From, To any](elements typ.Iterator[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements, by)
}

// additionally filters 'From' elements by filters
func MapFit[From, To any](elements typ.Iterator[From], fit typ.Predicate[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any](elements typ.Iterator[From], by typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements, by)
}

// additionally checks 'From' elements by fit
func FlattFit[From, To any](elements typ.Iterator[From], fit typ.Predicate[From], flatt typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any](elements typ.Iterator[T], filter typ.Predicate[T]) typ.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates a lazy Iterator that filters nullable elements
func NotNil[T any](elements typ.Iterator[*T]) typ.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

//ForEach applies func on elements
func ForEach[T any](elements typ.Iterator[T], apply func(T)) error {
	for elements.HasNext() {
		n, err := elements.Next()
		if err != nil {
			return err
		}
		apply(n)
	}
	return nil
}

func ForEachFit[T any](elements typ.Iterator[T], apply func(T), fit typ.Predicate[T]) {
	impl.ForEachFit(elements, apply, fit)
}

//Reduce reduces elements to an one
func Reduce[T any](elements typ.Iterator[T], by op.Binary[T]) T {
	return impl.Reduce(elements, by)
}

//Slice converts Iterator to slice
func Slice[T any](elements typ.Iterator[T]) []T {
	return impl.Slice[T](elements)
}

//Group groups elements to slices by a converter and returns map
func Group[T any, K comparable, IT typ.Iterator[T]](elements IT, by typ.Converter[T, K]) typ.MapPipe[K, T, map[K][]T] {
	return impl.NewKVPipe(impl.NewKeyValuer(elements, by), collect.Groups[K, T])
}
