package conv

import "github.com/m4gshm/gollections/loop"
import breakLoop "github.com/m4gshm/gollections/break/loop"

func Indexed[From, To any](len int, next func(int) From, converter func(from From) (To, error)) breakLoop.ConvertIter[From, To] {
	return loop.Conv(loop.OfIndexed(len, next), converter)
}
