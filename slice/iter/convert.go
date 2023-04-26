package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
)

// ConvertFit is the array based Iterator thath provides converting of elements by a Converter with addition filtering of the elements by a Predicate.
type ConvertFit[From, To any] struct {
	array      unsafe.Pointer
	elemSize   uintptr
	size, i    int
	converter  func(From) To
	filterFrom func(From) bool
	filterTo   func(To) bool
}

var _ c.Iterator[any] = (*ConvertFit[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *ConvertFit[From, To]) For(walker func(element To) error) error {
	return loop.For(f.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (f *ConvertFit[From, To]) ForEach(walker func(element To)) {
	loop.ForEach(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *ConvertFit[From, To]) Next() (t To, ok bool) {
	if s == nil || s.array == nil {
		return t, false
	}
	next := func() (From, bool) { return nextFiltered(s.array, s.size, s.elemSize, s.filterFrom, &s.i) }
	for v, ok := next(); ok; v, ok = next() {
		if t = s.converter(v); s.filterTo(t) {
			return t, true
		}
	}
	return t, false
}

// Cap returns the iterator capacity
func (s *ConvertFit[From, To]) Cap() int {
	return s.size
}

// ConvertIter is the array based Iterator thath provides converting of elements by a ConvertIter.
type ConvertIter[From, To any] struct {
	array     unsafe.Pointer
	elemSize  uintptr
	size, i   int
	converter func(From) To
}

var _ c.Iterator[any] = (*ConvertIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *ConvertIter[From, To]) For(walker func(element To) error) error {
	return loop.For(f.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (f *ConvertIter[From, To]) ForEach(walker func(element To)) {
	loop.ForEach(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *ConvertIter[From, To]) Next() (To, bool) {
	if s.i < s.size {
		v := *(*From)(notsafe.GetArrayElemRef(s.array, s.i, s.elemSize))
		s.i++
		return s.converter(v), true
	}
	var no To
	return no, false
}

// Cap returns the iterator capacity
func (s *ConvertIter[From, To]) Cap() int {
	return s.size
}

func nextFiltered[T any](array unsafe.Pointer, size int, elemSize uintptr, filter func(T) bool, index *int) (T, bool) {
	for i := *index; i < size; i++ {
		if v := *(*T)(notsafe.GetArrayElemRef(array, i, elemSize)); filter(v) {
			*index = i + 1
			return v, true
		}
	}
	var v T
	return v, false
}
