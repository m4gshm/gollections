package split

import "github.com/m4gshm/gollections/slice"

func Of[TS ~[]T, T, F, S any](elements TS, splitter func(T) (F, S)) ([]F, []S) {
	return slice.SplitTwo(elements, splitter)
}
