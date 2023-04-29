package stream

import "github.com/m4gshm/gollections/c"

// Stream is collection or stream of elements in transformation state
type Stream[T, I any] interface {
	c.Iterator[T]
	c.Collection[T, I]
}