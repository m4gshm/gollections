package iter

import "github.com/m4gshm/gollections/c"

// New creates an LoopIter instance that loops over elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func New[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) LoopIter[S, T] {
	return LoopIter[S, T]{source: source, hasNext: hasNext, getNext: getNext}
}

// LoopIter - universal c.Iterator implementation
type LoopIter[S, T any] struct {
	source  S
	hasNext func(S) bool
	getNext func(S) (T, error)
	abort   error
}

var (
	_ c.Iterator[any]          = (*LoopIter[any, any])(nil)
	_ c.IteratorBreakable[any] = (*LoopIter[any, any])(nil)
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *LoopIter[S, T]) Next() (next T, ok bool) {
	if i != nil {
		abort := i.abort
		if abort == nil && i.hasNext(i.source) {
			next, abort = i.getNext(i.source)
			i.abort = abort
			return next, abort == nil
		}
	}
	return next, false
}

// Error implements c.IteratorBreakable
func (i *LoopIter[S, T]) Error() error {
	if i == nil {
		return nil
	}
	return i.abort
}
