// Package seq provides helpers for “range-over-func” feature introduced in go 1.22.
package seq

import "iter"

// Filtered creates an iterator that iterates only those elements for which the 'filter' function returns true.
func Filter[T any](all iter.Seq[T], filter func(T) bool) iter.Seq[T] {

	return func(consumer func(T) bool) {
		all(func(e T) bool {
			if filter(e) {
				return consumer(e)
			}
			return true
		})
	}
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func Convert[From, To any](all iter.Seq[From], converter func(From) To) iter.Seq[To] {
	return func(consumer func(To) bool) {
		all(func(from From) bool {
			return consumer(converter(from))
		})
	}
}

// Conv creates an iterator that applies the 'converter' function to each iterable element and returns value-error pairs.
// The error should be checked at every iteration step, like:
//  var integers []int 
//  ...
//  for s, err := range seq.Conv(integers,  strconv.Itoa) {
//      if err != nil {
//          break
//      }
//      ...
//  }
func Conv[From, To any](all iter.Seq[From], converter func(From) (To, error)) iter.Seq2[To, error] {
	return func(consumer func(To, error) bool) {
		all(func(from From) bool {
			return consumer(converter(from))
		})
	}
}
