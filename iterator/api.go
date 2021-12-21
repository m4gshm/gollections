package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

type Iterator[T any] interface {
	HasNext() bool
	Get() T
	Next() (T, bool)
}

func ForEach[T any](items Iterator[T], apply func(T), filters ...check.Predicate[T]) {
	if len(filters) == 0 {
		for v, ok := items.Next(); ok; v, ok = items.Next() {
			apply(v)
		}
		return
	}
	for v, ok := items.Next(); ok; v, ok = items.Next() {
		if check.IsFit(v, filters...) {
			apply(v)
		}
	}
}

func Wrap[T any](elements []T) Iterator[T] {
	return &Slice[T]{elements: elements}
}

func WrapMap[K comparable, V any](values map[K]V) Iterator[*KV[K, V]] {
	return NewMap(values)
}

//Map applies Converter to items and accumulate to result slice
func Map[From, To any](items Iterator[From], by conv.Converter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	var (
		result = make([]To, 0)
		add    = func(v From) { result = append(result, by(v)) }
	)
	if len(filters) == 0 {
		for v, ok := items.Next(); ok; v, ok = items.Next() {
			add(v)
		}
	} else {
		for v, ok := items.Next(); ok; v, ok = items.Next() {
			if check.IsFit(v, filters...) {
				add(v)
			}
		}
	}
	return Wrap(result)
}

//Flatt extracts embedded slices of items by Flatter and accumulate to result slice
func Flatt[From, To any](items Iterator[From], by conv.Flatter[From, To], filters ...check.Predicate[From]) Iterator[To] {
	var (
		result = make([]To, 0)
		add    = func(v From) { result = append(result, by(v)...) }
	)
	if len(filters) == 0 {
		for v, ok := items.Next(); ok; v, ok = items.Next() {
			add(v)
		}
	} else {
		for v, ok := items.Next(); ok; v, ok = items.Next() {
			if check.IsFit(v, filters...) {
				add(v)
			}
		}
	}
	return Wrap(result)
}

//Filter tests items and adds that fit to result
func Filter[T any](items Iterator[T], filters ...check.Predicate[T]) Iterator[T] {
	var (
		result = make([]T, 0)
		add    = func(v T) { result = append(result, v) }
	)
	if len(filters) == 0 {
		for v, ok := items.Next(); ok; v, ok = items.Next() {
			add(v)
		}

	} else {
		for v, ok := items.Next(); ok; v, ok = items.Next() {
			if check.IsFit(v, filters...) {
				add(v)
			}
		}
	}
	return Wrap(result)
}

//NotNil filter nullable values
func NotNil[T any](items Iterator[T]) Iterator[T] {
	var (
		result = make([]T, 0)
		add    = func(v T) { result = append(result, v) }
	)
	for v, ok := items.Next(); ok; v, ok = items.Next() {
		if check.NotNil(v) {
			add(v)
		}
	}
	return Wrap(result)
}
