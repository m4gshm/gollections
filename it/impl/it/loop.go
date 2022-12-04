package it

import "github.com/m4gshm/gollections/c"

// NewLoop creates an LoopIter instance that loops over elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func NewLoop[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) LoopIter[S, T] {
	return LoopIter[S, T]{source: source, hasNext: hasNext, getNext: getNext}
}

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
func (i *LoopIter[S, T]) Next() (T, bool) {
	if i.abort == nil && i.hasNext(i.source) {
		if next, err := i.getNext(i.source); err == nil {
			return next, true
		} else {
			i.abort = err
		}
	}
	var no T
	return no, false
}

// Error implements c.IteratorBreakable
func (i *LoopIter[S, T]) Error() error {
	return i.abort
}
