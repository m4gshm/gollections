package split

import "github.com/m4gshm/gollections/slice"

func Of[TS ~[]T, T, F, S any](elements TS, firstConverter func(T) F, secondConverter func(T) S) ([]F, []S) {
	return slice.SplitTwo(elements, func(t T) (F, S) { return firstConverter(t), secondConverter(t) })
}

func AndReduce[TS ~[]T, T, F, S any](elements TS, firstConverter func(T) F, secondConverter func(T) S, firstMerge func(F, F) F, secondMerger func(S, S) S) (F, S) {
	return slice.SplitAndReduceTwo(elements, func(t T) (F, S) { return firstConverter(t), secondConverter(t) }, firstMerge, secondMerger)
}
