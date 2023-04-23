package loop

import (
	"github.com/m4gshm/gollections/c"
)

// ConvertFitIter is the Converter with elements filtering.
type ConvertFitIter[From, To any] struct {
	next   func() (From, bool)
	by     func(From) To
	filter func(From) bool
}

var (
	_ c.Iterator[any] = (*ConvertFitIter[any, any])(nil)
	_ c.Iterator[any] = ConvertFitIter[any, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertFitIter[From, To]) Next() (t To, ok bool) {
	if next, filter := c.next, c.filter; next != nil && filter != nil {
		if f, ok := nextFiltered(next, filter); ok {
			return c.by(f), true
		}
	}
	return t, false
}

// ConvertIter is the iterator wrapper implementation applying a converter to all iterable elements.
type ConvertIter[From, To any] struct {
	next      func() (From, bool)
	converter func(From) To
}

var (
	_ c.Iterator[any] = (*ConvertIter[any, any])(nil)
	_ c.Iterator[any] = ConvertIter[any, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertIter[From, To]) Next() (t To, ok bool) {
	if next := c.next; next != nil {
		if v, ok := next(); ok {
			return c.converter(v), true
		}
	}
	return t, false
}
