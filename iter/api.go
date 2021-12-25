package iter

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	in "github.com/m4gshm/container/iter/impl/iter"
	impl "github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/typ"
)

//Of Iterator constructor
func Of[T any](elements ...T) typ.Iterator[T] {
	return impl.NewReseteable(elements)
}

//Wrap slice to Iterator
func Wrap[T any](elements []T) typ.Iterator[T] {
	return impl.NewReseteable(elements)
}

// WrapMap Key, value Iterator constructor.
func WrapMap[K comparable, V any](values map[K]V) typ.Iterator[*KV[K, V]] {
	return NewMap(values)
}

//Map creates a lazy Iterator that converts elements with a converter and returns them
func Map[From, To any](elements typ.Iterator[From], by conv.Converter[From, To]) typ.Iterator[To] {
	return in.Map(elements, by)
}

// additionally filters 'From' elements by filters
func MapFit[From, To any](elements typ.Iterator[From], fit check.Predicate[From], by conv.Converter[From, To]) typ.Iterator[To] {
	return in.MapFit(elements, fit, by)
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any](elements typ.Iterator[From], by conv.Flatter[From, To]) typ.Iterator[To] {
	return in.Flatt(elements, by)
}

// additionally checks 'From' elements by fit
func FlattFit[From, To any](elements typ.Iterator[From], fit check.Predicate[From], flatt conv.Flatter[From, To]) typ.Iterator[To] {
	return in.FlattFit(elements, fit, flatt)
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any](elements typ.Iterator[T], filter check.Predicate[T]) typ.Iterator[T] {
	return in.Filter(elements, filter)
}

//NotNil creates a lazy Iterator that filters nullable elements
func NotNil[T any](elements typ.Iterator[*T]) typ.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

//ForEach applies func on elements
func ForEach[T any](elements typ.Iterator[T], apply func(T)) {
	in.ForEach(elements, apply)
}

func ForEachFit[T any](elements typ.Iterator[T], apply func(T), fit check.Predicate[T]) {
	in.ForEachFit(elements, apply, fit)
}

//Reduce reduces elements to an one
func Reduce[T any](elements typ.Iterator[T], by op.Binary[T]) T {
	return in.Reduce(elements, by)
}

//ToSlice converts Iterator to slice
func ToSlice[T any](elements typ.Iterator[T]) []T {
	return in.ToSlice[T](elements)
}
