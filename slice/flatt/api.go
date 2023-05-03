// Package flatt provides short aliases for slice functions
package flatt

import "github.com/m4gshm/gollections/slice"

// AndConvert - flatt.AndConvert alias
func AndConvert[FS ~[]From, From, I, To any](elements FS, flattener func(From) []I, convert func(I) To) []To {
	return slice.FlattAndConvert(elements, flattener, convert)
}
