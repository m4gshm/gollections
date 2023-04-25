// Package flatt provides short aliases for loop functions
package flatt

import "github.com/m4gshm/gollections/loop"

// AndConvert - flatt.AndConvert flattens and converts elements retrieved by the 'next' function
func AndConvert[From, I, To any](next func() (From, bool), flattener func(From) []I, convert func(I) To) loop.ConvertIter[I, To] {
	f := loop.Flatt(next, flattener)

	return loop.Convert((&f).Next, convert)
}
