// Package flat provides short aliases for slice functions
package flat

import "github.com/m4gshm/gollections/slice"

// AndConvert - flattener.AndConvert alias
func AndConvert[FS ~[]From, From any, IS ~[]I, I, To any](elements FS, flattener func(From) IS, convert func(I) To) []To {
	return slice.FlatAndConvert(elements, flattener, convert)
}
