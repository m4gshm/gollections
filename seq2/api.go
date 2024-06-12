// Package seq2 provides helpers for  “range-over-func” feature introduced in go 1.22.
package seq2

import "iter"

// Filtered creates a rangefunc that iterates only those elements for which the 'filter' function returns true.
func Filter[I, T any](all iter.Seq2[I, T], filter func(I, T) bool) iter.Seq2[I, T] {
	return func(consumer func(I, T) bool) {
		all(func(i I, e T) bool {
			if filter(i, e) {
				return consumer(i, e)
			}
			return true
		})
	}
}

// Convert creates a rangefunc that applies the 'converter' function to each iterable element.
func Convert[I, From, To any](all iter.Seq2[I, From], converter func(I, From) To) iter.Seq2[I, To] {
	return func(consumer func(I, To) bool) {
		all(func(i I, from From) bool {
			return consumer(i, converter(i, from))
		})
	}
}

func ToSeq[I, T any](all iter.Seq2[I, T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		all(func(i I, e T) bool {
			return yield(e)
		})
	}
}
