package loop

import "github.com/m4gshm/gollections/c"

// NewIter creates an Iter instance that loops over elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func NewIter[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) Iter[S, T] {
	return Iter[S, T]{source: source, hasNext: hasNext, next: getNext}
}

// Iter - universal c.Iterator implementation
type Iter[S, T any] struct {
	source  S
	hasNext func(S) bool
	next    func(S) (T, error)
	abort   error
}

var (
	_ c.Iterator[any]          = (*Iter[any, any])(nil)
	_ c.IteratorBreakable[any] = (*Iter[any, any])(nil)
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *Iter[S, T]) For(walker func(element T) error) error {
	return For(i.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i *Iter[S, T]) ForEach(walker func(element T)) {
	ForEach(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *Iter[S, T]) Next() (next T, ok bool) {
	if i != nil {
		abort := i.abort
		if abort == nil && i.hasNext(i.source) {
			next, abort = i.next(i.source)
			i.abort = abort
			return next, abort == nil
		}
	}
	return next, false
}

// Error implements c.IteratorBreakable
func (i *Iter[S, T]) Error() error {
	if i == nil {
		return nil
	}
	return i.abort
}
