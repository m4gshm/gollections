package stream

import "github.com/m4gshm/gollections/break/c"

// Stream is collection or stream of elements in transformation state.
// It supports interrupting on an error that may occur in intermediate or final executor functions.
type Stream[T any, I c.Iterator[T]] interface {
	c.Iterator[T]
	Iter() I

	Slice() ([]T, error)

	// Filt(predicate func(T) (bool, error)) StreamBreakable[T]
	// Filter(predicate func(T) bool) StreamBreakable[T]
	// Conv(converter func(T) (T, error)) StreamBreakable[T]
	// Convert(converter func(T) T) StreamBreakable[T]

	Reduce(merger func(T, T) (T, error)) (T, error)
	HasAny(predicate func(T) (bool, error)) (bool, error)
}
