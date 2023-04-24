package flatt

import "github.com/m4gshm/gollections/loop"

func AndConvert[From, I, To any](next func() (From, bool), flattener func(From) []I, convert func(I) To) loop.ConvertIter[I, To] {
	f := loop.Flatt(next, flattener)

	return loop.Convert((&f).Next, convert)
}
