package loop

import (
	"github.com/m4gshm/gollections/break/c"
)

// ConvFiltIter iterator implementation that retrieves an element by the 'next' function, converts by the 'converter' and addition checks by the 'filter'.
// If the filter returns true then the converted element is returned as next.
type ConvFiltIter[From, To any] struct {
	next       func() (From, bool, error)
	converter  func(From) (To, error)
	filterFrom func(From) (bool, error)
	filterTo   func(To) (bool, error)
}

var (
	_ c.Iterator[any] = (*ConvFiltIter[any, any])(nil)
	_ c.Iterator[any] = ConvFiltIter[any, any]{}
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (c ConvFiltIter[From, To]) For(walker func(element To) error) error {
	return For(c.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvFiltIter[From, To]) Next() (t To, ok bool, err error) {
	next, filterFrom, filterTo := c.next, c.filterFrom, c.filterTo
	if next == nil || filterFrom == nil || filterTo == nil {
		return t, false, nil
	}
	for {
		if f, ok, err := nextFiltered(next, filterFrom); err != nil || !ok {
			return t, false, err
		} else if cf, err := c.converter(f); err != nil {
			return t, false, err
		} else if ok, err := filterTo(cf); err != nil || !ok {
			return t, false, err
		} else {
			return cf, true, nil
		}
	}
}

// ConvertIter iterator implementation that retrieves an element by the 'next' function and converts by the 'converter'
type ConvertIter[From, To any] struct {
	next      func() (From, bool, error)
	converter func(From) (To, error)
}

var (
	_ c.Iterator[any] = (*ConvertIter[any, any])(nil)
	_ c.Iterator[any] = ConvertIter[any, any]{}
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (c ConvertIter[From, To]) For(walker func(element To) error) error {
	return For(c.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertIter[From, To]) Next() (t To, ok bool, err error) {
	if next := c.next; next == nil {
		return t, false, nil
	} else if v, ok, err := next(); err != nil || !ok {
		return t, false, err
	} else {
		vc, err := c.converter(v)
		return vc, err == nil, err
	}
}

// ConvertCheckIter converts and filters elements at the same time
type ConvertCheckIter[From, To any] struct {
	next      func() (From, bool, error)
	converter func(From) (To, bool, error)
}

var (
	_ c.Iterator[any] = (*ConvertIter[any, any])(nil)
	_ c.Iterator[any] = ConvertIter[any, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertCheckIter[From, To]) Next() (t To, ok bool, err error) {
	next, converter := c.next, c.converter
	if next == nil || converter == nil {
		return t, false, nil
	}
	for {
		if e, ok, err := next(); err != nil || !ok {
			return t, false, err
		} else if t, ok, err := converter(e); err != nil || ok {
			return t, ok, err
		}
	}
}
