package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/break/c"
	"github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/notsafe"
)

// ConvFitIter is the array based Iterator thath provides converting of elements by a Converter with addition filtering of the elements by a Predicate.
type ConvFitIter[From, To any] struct {
	array      unsafe.Pointer
	elemSize   uintptr
	size, i    int
	converter  func(From) (To, error)
	filterFrom func(From) (bool, error)
	filterTo   func(To) (bool, error)
}

var _ c.Iterator[any] = (*ConvFitIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *ConvFitIter[From, To]) For(walker func(element To) error) error {
	return loop.For(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *ConvFitIter[From, To]) Next() (t To, ok bool, err error) {
	if s == nil || s.array == nil {
		return t, false, nil
	}
	next := func() (out From, ok bool, err error) {
		return nextFiltered(s.array, s.size, s.elemSize, s.filterFrom, &s.i)
	}
	for {
		if v, ok, err := next(); err != nil || !ok {
			return t, false, err
		} else if t, err = s.converter(v); err != nil {
			return t, false, err
		} else if ok, err := s.filterTo(t); err != nil {
			return t, false, err
		} else if ok {
			return t, true, nil
		}
	}
	return t, false, nil
}

// Cap returns the iterator capacity
func (s *ConvFitIter[From, To]) Cap() int {
	return s.size
}

// ConvertIter is the array based Iterator thath provides converting of elements by a ConvertIter.
type ConvertIter[From, To any] struct {
	array     unsafe.Pointer
	elemSize  uintptr
	size, i   int
	converter func(From) (To, error)
}

var _ c.Iterator[any] = (*ConvertIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *ConvertIter[From, To]) For(walker func(element To) error) error {
	return loop.For(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *ConvertIter[From, To]) Next() (t To, ok bool, err error) {
	if s.i < s.size {
		v := *(*From)(notsafe.GetArrayElemRef(s.array, s.i, s.elemSize))
		s.i++
		t, err = s.converter(v)
		return t, err == nil, err
	}
	return t, false, nil
}

// Cap returns the iterator capacity
func (s *ConvertIter[From, To]) Cap() int {
	return s.size
}

func nextFiltered[T any](array unsafe.Pointer, size int, elemSize uintptr, filter func(T) (bool, error), index *int) (T, bool, error) {
	for i := *index; i < size; i++ {
		v := *(*T)(notsafe.GetArrayElemRef(array, i, elemSize))
		if ok, err := filter(v); err != nil || !ok {
			return v, false, err
		} else {
			*index = i + 1
			return v, true, nil
		}
	}
	var v T
	return v, false, nil
}
