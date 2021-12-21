package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

type Iterator[T any] interface {
	Next() bool
	Get() T
}

func ForEach[T any](items Iterator[T], apply func(T), filters ...check.Predicate[T]) {
	if len(filters) == 0 {
		for items.Next() {
			apply(items.Get())
		}
		return
	}
	for items.Next() {
		if v := items.Get(); check.IsFit(v, filters...) {
			apply(items.Get())
		}
	}
}

func Wrap[T any](values []T) Iterator[T] {
	return &SliceIter[T]{values: values}
}

func New[T any](values ...T) Iterator[T] {
	return &SliceIter[T]{values: values}
}

func WrapMap[K comparable, V any](values map[K]V) Iterator[*KV[K, V]] {
	return wrapMap(values)
}

//Map applies Converter to items and accumulate to result slice
func Map[From, To any](items Iterator[From], by conv.Converter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	var (
		result = make([]To, 0)
		add    = func(v From) { result = append(result, by(v)) }
	)
	ForEach(items, add, filters...)
	return Wrap(result)
}

//Flatt extracts embedded slices of items by Flatter and accumulate to result slice
func Flatt[From, To any](items Iterator[From], by conv.Flatter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	result := make([]To, 0)
	ForEach(items, func(v From) { result = append(result, by(v)...) }, filters...)
	return Wrap(result)
}

//Filter tests items and adds that fit to result
func Filter[T any](items Iterator[T], filters ...check.Predicate[T]) Iterator[T] {
	result := make([]T, 0)
	ForEach(items, func(v T) { result = append(result, v) }, filters...)
	return Wrap(result)
}

//NotNil filter nullable values
func NotNil[T any](items Iterator[T]) Iterator[T] {
	result := make([]T, 0)
	ForEach(items, func(v T) { result = append(result, v) }, check.NotNil[T])
	return Wrap(result)
}
