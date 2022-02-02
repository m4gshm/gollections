package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/op"
)

//Map creates the Iterator that converts elements with a converter and returns them.
func Map[From, To any, IT c.Iterator[From]](elements IT, by c.Converter[From, To]) *Convert[From, To, IT, c.Converter[From, To]] {
	return &Convert[From, To, IT, c.Converter[From, To]]{Iter: elements, By: by}
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any, IT c.Iterator[From]](elements IT, fit c.Predicate[From], by c.Converter[From, To]) *ConvertFit[From, To] {
	return &ConvertFit[From, To]{Iter: elements, By: by, Fit: fit}
}

//Flatt creates the Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, IT c.Iterator[From]](elements IT, by c.Flatter[From, To]) *Flatten[From, To] {
	return &Flatten[From, To]{Iter: elements, Flatt: by}
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any, IT c.Iterator[From]](elements IT, fit c.Predicate[From], flatt c.Flatter[From, To]) *FlattenFit[From, To] {
	return &FlattenFit[From, To]{Iter: elements, Flatt: flatt, Fit: fit}
}

//Filter creates the Iterator that checks elements by filters and returns successful ones.
func Filter[T any, IT c.Iterator[T]](elements IT, filter c.Predicate[T]) *Fit[T] {
	return &Fit[T]{Iter: elements, By: filter}
}

//NotNil creates the Iterator that filters nullable elements.
func NotNil[T any, IT c.Iterator[*T]](elements IT) *Fit[*T] {
	return Filter(elements, check.NotNil[T])
}

//MapKV creates the Iterator that converts elements with a converter and returns them.
func MapKV[k, v any, IT c.KVIterator[k, v], k2, v2 any](elements IT, by c.BiConverter[k, v, k2, v2]) *ConvertKV[k, v, IT, k2, v2, c.BiConverter[k, v, k2, v2]] {
	return &ConvertKV[k, v, IT, k2, v2, c.BiConverter[k, v, k2, v2]]{Iter: elements, By: by}
}

//FilterKV creates the Iterator that checks elements by filters and returns successful ones.
func FilterKV[k, v any, IT c.KVIterator[k, v]](elements IT, filter c.BiPredicate[k, v]) *FitKV[k, v] {
	return &FitKV[k, v]{Iter: elements, By: filter}
}

//Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable, IT c.Iterator[T]](elements IT, by c.Converter[T, K]) c.MapPipe[K, T, map[K][]T] {
	return NewKVPipe(NewKeyValuer(elements, by), collect.Groups[K, T])
}

//For applies func on elements.
func For[T any, IT c.Iterator[T]](elements IT, apply func(T) error) error {
	for elements.HasNext() {
		err:= apply(elements.Next())
		if err != nil {
			return err
		}
	}
	return nil
}

//ForEach applies func on elements.
func ForEach[T any, IT c.Iterator[T]](elements IT, apply func(T)) {
	for elements.HasNext() {
		apply(elements.Next())
	}
}

func ForEachFit[T any, IT c.Iterator[T]](elements IT, apply func(T), fit c.Predicate[T]) {
	for elements.HasNext() {
		if v := elements.Next(); fit(v) {
			apply(v)
		}
	}
}

//Reduce reduces elements to an one.
func Reduce[T any, IT c.Iterator[T]](elements IT, by op.Binary[T]) T {
	var result T
	for elements.HasNext() {
		result = by(result, Next[T](elements))
	}
	return result
}

func ReduceKV[k, v any, IT c.KVIterator[k, v]](elements IT, by op.Quaternary[k, v]) (k, v) {
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
func Slice[T any, IT c.Iterator[T]](elements IT) []T {
	s := make([]T, 0)
	for elements.HasNext() {
		s = append(s, Next[T](elements))
	}
	return s
}

func Next[T any, IT c.Iterator[T]](elements IT) T {
	return elements.Next()

}

func NextKV[k, v any, IT c.KVIterator[k, v]](elements IT) (k, v) {
	return elements.Next()
}
