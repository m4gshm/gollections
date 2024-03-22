package stream

import "github.com/m4gshm/gollections/break/loop"

// Stream is collection or stream of elements in transformation state.
// It supports interrupting on an error that may occur in intermediate or final executor functions.
type Stream[T any] interface {
	Loop() loop.Loop[T]

	Slice() ([]T, error)

	Reduce(merger func(T, T) (T, error)) (T, error)
	HasAny(predicate func(T) (bool, error)) (bool, error)
}
