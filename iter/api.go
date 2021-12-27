package iter

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
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
func WrapMap[k comparable, v any](values map[k]v) typ.Iterator[*typ.KV[k, v]] {
	return impl.NewKV(values)
}

//Map creates a lazy Iterator that converts elements with a converter and returns them
func Map[From, To any](elements typ.Iterator[From], by conv.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements, by)
}

// additionally filters 'From' elements by filters
func MapFit[From, To any](elements typ.Iterator[From], fit check.Predicate[From], by conv.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any](elements typ.Iterator[From], by conv.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements, by)
}

// additionally checks 'From' elements by fit
func FlattFit[From, To any](elements typ.Iterator[From], fit check.Predicate[From], flatt conv.Flatter[From, To]) typ.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any](elements typ.Iterator[T], filter check.Predicate[T]) typ.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates a lazy Iterator that filters nullable elements
func NotNil[T any](elements typ.Iterator[*T]) typ.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

//ForEach applies func on elements
func ForEach[T any](elements typ.Iterator[T], apply func(T)) {
	impl.ForEach(elements, apply)
}

func ForEachFit[T any](elements typ.Iterator[T], apply func(T), fit check.Predicate[T]) {
	impl.ForEachFit(elements, apply, fit)
}

//Reduce reduces elements to an one
func Reduce[T any](elements typ.Iterator[T], by op.Binary[T]) T {
	return impl.Reduce(elements, by)
}

//ToSlice converts Iterator to slice
func ToSlice[T any](elements typ.Iterator[T]) []T {
	return impl.ToSlice[T](elements)
}
