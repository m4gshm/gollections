package it

import (
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

//Map creates the Iterator that converts elements with a converter and returns them.
func Map[From, To any, IT typ.Iterator[From]](elements IT, by typ.Converter[From, To]) *Convert[From, To, IT, typ.Converter[From, To]] {
	return &Convert[From, To, IT, typ.Converter[From, To]]{Iter: elements, By: by}
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any, IT typ.Iterator[From]](elements IT, fit typ.Predicate[From], by typ.Converter[From, To]) *ConvertFit[From, To] {
	return &ConvertFit[From, To]{Iter: elements, By: by, Fit: fit}
}

//Flatt creates the Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, IT typ.Iterator[From]](elements IT, by typ.Flatter[From, To]) *Flatten[From, To] {
	return &Flatten[From, To]{Iter: elements, Flatt: by}
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any, IT typ.Iterator[From]](elements IT, fit typ.Predicate[From], flatt typ.Flatter[From, To]) *FlattenFit[From, To] {
	return &FlattenFit[From, To]{Iter: elements, Flatt: flatt, Fit: fit}
}

//Filter creates the Iterator that checks elements by filters and returns successful ones.
func Filter[T any, IT typ.Iterator[T]](elements IT, filter typ.Predicate[T]) *Fit[T] {
	return &Fit[T]{Iter: elements, By: filter}
}

//NotNil creates the Iterator that filters nullable elements.
func NotNil[T any, IT typ.Iterator[*T]](elements IT) *Fit[*T] {
	return Filter(elements, check.NotNil[T])
}

//MapKV creates the Iterator that converts elements with a converter and returns them.
func MapKV[k, v any, IT typ.KVIterator[k, v], k2, v2 any](elements IT, by typ.BiConverter[k, v, k2, v2]) *ConvertKV[k, v, IT, k2, v2, typ.BiConverter[k, v, k2, v2]] {
	return &ConvertKV[k, v, IT, k2, v2, typ.BiConverter[k, v, k2, v2]]{Iter: elements, By: by}
}

//FilterKV creates the Iterator that checks elements by filters and returns successful ones.
func FilterKV[k, v any, IT typ.KVIterator[k, v]](elements IT, filter typ.BiPredicate[k, v]) *FitKV[k, v] {
	return &FitKV[k, v]{Iter: elements, By: filter}
}

//Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable, IT typ.Iterator[T]](elements IT, by typ.Converter[T, K]) typ.MapPipe[K, T, map[K][]T] {
	return NewKVPipe(NewKeyValuer(elements, by), collect.Groups[K, T])
}

//For applies func on elements.
func For[T any, IT typ.Iterator[T]](elements IT, apply func(T)) error {
	for elements.HasNext() {
		v, err := elements.Get()
		if err != nil {
			return err
		}
		apply(v)
	}
	return nil
}

func ForFit[T any, IT typ.Iterator[T]](elements IT, apply func(T), fit typ.Predicate[T]) error {
	for elements.HasNext() {
		if v, err := elements.Get(); err != nil {
			return err
		} else if fit(v) {
			apply(v)
		}
	}
	return nil
}

//Reduce reduces elements to an one.
func Reduce[T any, IT typ.Iterator[T]](elements IT, by op.Binary[T]) T {
	if !elements.HasNext() {
		var nothing T
		return nothing
	}
	result := Next[T](elements)
	for elements.HasNext() {
		result = by(result, Next[T](elements))
	}
	return result
}

func ReduceKV[k, v any, IT typ.KVIterator[k, v]](elements IT, by op.Quaternary[k, v]) (k, v) {
	if !elements.HasNext() {
		var key k
		var val v
		return key, val
	}
	key, val := NextKV[k, v](elements)
	for elements.HasNext() {
		key2, val2 := NextKV[k, v](elements)
		key, val = by(key, val, key2, val2)
	}
	return key, val
}

//Slice converts an Iterator to a slice.
func Slice[T any, IT typ.Iterator[T]](elements IT) []T {
	s := make([]T, 0)

	for elements.HasNext() {
		s = append(s, Next[T](elements))
	}

	return s
}

func Next[T any, IT typ.Iterator[T]](elements IT) T {
	n, err := elements.Get()
	if err != nil {
		panic(err)
	}
	return n
}

func NextKV[k, v any, IT typ.KVIterator[k, v]](elements IT) (k, v) {
	key, val, err := elements.Get()
	if err != nil {
		panic(err)
	}
	return key, val
}
