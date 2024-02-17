package loop

import (
	"github.com/m4gshm/gollections/c"
)

// ConvertFiltIter iterator implementation that retrieves an element by the 'next' function, converts by the 'converter' and addition checks by the 'filter'.
// If the filter returns true then the converted element is returned as next.
type ConvertFiltIter[From, To any] struct {
	next       func() (From, bool)
	converter  func(From) To
	filterFrom func(From) bool
	filterTo   func(To) bool
}

var (
	_ c.Iterator[any] = (*ConvertFiltIter[any, any])(nil)
	_ c.Iterator[any] = ConvertFiltIter[any, any]{}
)

var _ c.IterFor[any, ConvertFiltIter[any, any]] = ConvertFiltIter[any, any]{}

func (c ConvertFiltIter[From, To]) All(consumer func(element To) bool) {
	All(c.Next, consumer)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (c ConvertFiltIter[From, To]) For(walker func(element To) error) error {
	return For(c.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (c ConvertFiltIter[From, To]) ForEach(walker func(element To)) {
	ForEach(c.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertFiltIter[From, To]) Next() (t To, ok bool) {
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

// Start is used with for loop construct like 'for i, val, ok := i.Start(); ok; val, ok = i.Next() { }'
func (c ConvertFiltIter[From, To]) Start() (ConvertFiltIter[From, To], To, bool) {
	return startIt[To](c)
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

var _ c.IterFor[any, ConvertIter[any, any]] = ConvertIter[any, any]{}

func (c ConvertIter[From, To]) All(consumer func(element To) bool) {
	All(c.Next, consumer)
}

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

// Start is used with for loop construct like 'for i, val, ok := i.Start(); ok; val, ok = i.Next() { }'
func (c ConvertIter[From, To]) Start() (ConvertIter[From, To], To, bool) {
	return startIt[To](c)
}

// ConvertCheckIter converts and filters elements at the same time
type ConvertCheckIter[From, To any] struct {
	next      func() (From, bool)
	converter func(From) (To, bool)
}

var (
	_ c.Iterator[any] = (*ConvertCheckIter[any, any])(nil)
	_ c.Iterator[any] = ConvertCheckIter[any, any]{}
)

var _ c.IterFor[any, ConvertCheckIter[any, any]] = ConvertCheckIter[any, any]{}

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

func (c ConvertCheckIter[From, To]) All(consumer func(element To) bool) {
	All(c.Next, consumer)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (c ConvertCheckIter[From, To]) For(walker func(element To) error) error {
	return For(c.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (c ConvertCheckIter[From, To]) ForEach(walker func(element To)) {
	ForEach(c.Next, walker)
}

// Start is used with for loop construct like 'for i, val, ok := i.Start(); ok; val, ok = i.Next() { }'
func (c ConvertCheckIter[From, To]) Start() (ConvertCheckIter[From, To], To, bool) {
	return startIt[To](c)
}
