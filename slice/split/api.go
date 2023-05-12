// Package splitprovides utils for splitting slices
package split

import "github.com/m4gshm/gollections/slice"

// Of splits the elements into two slices
func Of[TS ~[]T, T, F, S any](elements TS, firstConverter func(T) F, secondConverter func(T) S) ([]F, []S) {
	return slice.SplitTwo(elements, func(t T) (F, S) { return firstConverter(t), secondConverter(t) })
}

// AndReduce - split.AndReduce splits each element of the specified slice into two values and then reduces that ones
func AndReduce[TS ~[]T, T, F, S any](elements TS, firstConverter func(T) F, secondConverter func(T) S, firstMerge func(F, F) F, secondMerger func(S, S) S) (F, S) {
	return slice.SplitAndReduceTwo(elements, func(t T) (F, S) { return firstConverter(t), secondConverter(t) }, firstMerge, secondMerger)
}
