// Package flat provides short aliases for loop functions
package flat

import "github.com/m4gshm/gollections/loop"

// AndConvert - flattener.AndConvert flattens and converts elements retrieved by the 'next' function
func AndConvert[From, I, To any](next func() (From, bool), flattener func(From) []I, convert func(I) To) loop.Loop[To] {
	return loop.Convert(loop.Flat(next, flattener), convert)
}
