package stream

import "github.com/m4gshm/gollections/c"

// Stream is collection or stream of elements in transformation state
type Stream[T any, I c.Iterator[T]] interface {
	c.Iterator[T]
	c.Collection[T, I]

	HasAny(func(T) bool) bool
}
