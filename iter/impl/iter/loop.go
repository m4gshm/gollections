package iter

import "github.com/m4gshm/gollections/c"

// NewLoop creates an LoopIter instance that loops over elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func NewLoop[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) LoopIter[S, T] {
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

// Next implements c.Iterator
func (i *LoopIter[S, T]) Next() (t T, ok bool) {
	if i != nil && i.abort == nil && i.hasNext(i.source) {
		next, err := i.getNext(i.source)
		if err == nil {
			return next, true
		}
		i.abort = err
	}
	return
}

// Error implements c.IteratorBreakable
func (i *LoopIter[S, T]) Error() error {
	if i == nil {
		return nil
	}
	return i.abort
}