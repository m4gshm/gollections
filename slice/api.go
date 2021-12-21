package slice

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

//Of slice constructor
func Of[T any](values ...T) []T { return values }

func ForEach[T any](items []T, apply func(T), filters ...check.Predicate[T]) {
	if len(filters) == 0 {
		for _, v := range items {
			apply(v)
		}
		return
	}
	for _, v := range items {
		if check.IsFit(v, filters...) {
			apply(v)
		}
	}
}

//Map applies Converter to items and accumulate to result slice
func Map[From, To any](items []From, by conv.Converter[From, To], filters ...check.Predicate[From]) []To {
	result := make([]To, 0)
	ForEach(items, func(v From) { result = append(result, by(v)) }, filters...)
	return result
}

//Flatt extracts embedded slices of items by Flatter and accumulate to result slice
func Flatt[From, To any](items []From, by conv.Flatter[From, To], filters ...check.Predicate[From]) []To {
	result := make([]To, 0)
	ForEach(items, func(v From) { result = append(result, by(v)...) }, filters...)
	return result
}

//Filter tests items and adds that fit to result
func Filter[T any](items []T, filters ...check.Predicate[T]) []T {
	result := make([]T, 0)
	ForEach(items, func(v T) { result = append(result, v) }, filters...)
	return result
}

//NotNil filter nullable values
func NotNil[T any](items []T) []T {
	result := make([]T, 0)
	ForEach(items, func(v T) { result = append(result, v) }, check.NotNil[T])
	return result
}
