// Package flat provides short aliases for slice functions
package flat

import "github.com/m4gshm/gollections/slice"

// AndConvert - flattener.AndConvert alias
func AndConvert[FS ~[]From, From, I, To any](elements FS, flattener func(From) []I, convert func(I) To) []To {
	return slice.FlattAndConvert(elements, flattener, convert)
}
