package stream

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// Stream is collection or stream of elements in transformation state
type Stream[T any] interface {
	c.Iterator[T]
	c.Iterable[T, loop.Loop[T]]
	c.Collection[T]
}
