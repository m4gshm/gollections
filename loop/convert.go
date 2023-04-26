package loop

import (
	"github.com/m4gshm/gollections/c"
)

// ConvertFitIter iterator implementation that retrieves an element by the 'next' function, converts by the 'converter' and addition checks by the 'filter'.
// If the filter returns true then the converted element is returned as next.
type ConvertFitIter[From, To any] struct {
	next       func() (From, bool)
	converter  func(From) To
	filterFrom func(From) bool
	filterTo   func(To) bool
}

var (
	_ c.Iterator[any] = (*ConvertFitIter[any, any])(nil)
	_ c.Iterator[any] = ConvertFitIter[any, any]{}
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (c ConvertFitIter[From, To]) For(walker func(element To) error) error {
	return For(c.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (c ConvertFitIter[From, To]) ForEach(walker func(element To)) {
	ForEach(c.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertFitIter[From, To]) Next() (t To, ok bool) {
	if next, filterFrom := c.next, c.filterFrom; next != nil && filterFrom != nil {
		for f, ok := nextFiltered(next, filterFrom); ok; f, ok = nextFiltered(next, filterFrom) {
			if filterTo := c.filterTo; filterTo != nil {
				if t = c.converter(f); filterTo(t) {
					return t, ok
				}
			}
		}
	}
	return t, false
}

// ConvertIter iterator implementation that retrieves an element by the 'next' function and converts by the 'converter'
type ConvertIter[From, To any] struct {
	next      func() (From, bool)
	converter func(From) To
}

var (
	_ c.Iterator[any] = (*ConvertIter[any, any])(nil)
	_ c.Iterator[any] = ConvertIter[any, any]{}
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (c ConvertIter[From, To]) For(walker func(element To) error) error {
	return For(c.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (c ConvertIter[From, To]) ForEach(walker func(element To)) {
	ForEach(c.Next, walker)
}

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

// ConvertCheckIter converts and filters elements at the same time
type ConvertCheckIter[From, To any] struct {
	next      func() (From, bool)
	converter func(From) (To, bool)
}

var (
	_ c.Iterator[any] = (*ConvertIter[any, any])(nil)
	_ c.Iterator[any] = ConvertIter[any, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertCheckIter[From, To]) Next() (t To, ok bool) {
	if next, converter := c.next, c.converter; next != nil && converter != nil {
		for e, ok := next(); ok; e, ok = next() {
			if t, ok := converter(e); ok {
				return t, true
			}
		}
	}
	return t, false
}
