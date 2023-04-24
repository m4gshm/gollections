package flatt

import "github.com/m4gshm/gollections/slice"

func AndConvert[FS ~[]From, From, I, To any](elements FS, flattener func(From) []I, convert func(I) To) []To {
	return slice.FlattAndConvert(elements, flattener, convert)
}