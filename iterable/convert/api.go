package convert

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/not"
)

func AndConvert[From, To, Too any](elements c.Iterable[From], firsConverter func(From) To, secondConverter func(To) Too) c.Pipe[Too] {
	return iterable.Convert(iterable.Convert(elements, firsConverter), secondConverter)
}

func AndFilter[From, To any](elements c.Iterable[From], converter func(From) To, filter func(To) bool) c.Pipe[To] {
	return iterable.Filter(iterable.Convert(elements, converter), filter)
}

func NotNil[From, To any](elements c.Iterable[*From], converter func(*From) To) c.Pipe[To] {
	return iterable.FilterAndConvert(elements, not.Nil[From], converter)
}
