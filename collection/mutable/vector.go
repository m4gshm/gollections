package mutable

import (
	"fmt"
	"sort"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
)

// WrapVector instantiates Vector using a slise as internal storage
func WrapVector[T any](elements []T) *Vector[T] {
	v := Vector[T](elements)
	return &v
}

// Vector is a collection implementation that provides elements order and index access
type Vector[T any] []T

var (
	_ c.Addable[any]                    = (*Vector[any])(nil)
	_ c.AddableAll[c.ForEachLoop[any]]  = (*Vector[any])(nil)
	_ c.Deleteable[int]                 = (*Vector[any])(nil)
	_ c.DeleteableVerify[int]           = (*Vector[any])(nil)
	_ c.Settable[int, any]              = (*Vector[any])(nil)
	_ c.SettableNew[int, any]           = (*Vector[any])(nil)
	_ collection.Vector[any]            = (*Vector[any])(nil)
	_ loop.Looper[any, *SliceIter[any]] = (*Vector[any])(nil)
	_ fmt.Stringer                      = (*Vector[any])(nil)
)

func (v *Vector[T]) All(yield func(T) bool) {
	if v != nil {
		slice.All(*v, yield)
	}
}

// Iter creates an iterator and returns as interface
func (v *Vector[T]) Iter() c.Iterator[T] {
	h := v.Head()
	return &h
}

// Loop creates an iterator and returns as implementation type reference
func (v *Vector[T]) Loop() *SliceIter[T] {
	h := v.Head()
	return &h
}

// IterEdit creates iterator that can delete iterable elements
func (v *Vector[T]) IterEdit() c.DelIterator[T] {
	h := v.Head()
	return &h
}

// Head creates an iterator and returns as implementation type value
func (v *Vector[T]) Head() SliceIter[T] {
	return NewHead(v, v.DeleteActualOne)
}

// Tail creates an iterator pointing to the end of the collection
func (v *Vector[T]) Tail() SliceIter[T] {
	return NewTail(v, v.DeleteActualOne)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (v *Vector[T]) First() (SliceIter[T], T, bool) {
	var (
		iterator  = NewHead(v, v.DeleteActualOne)
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Last returns the latest element of the collection, an iterator to reverse iterate over the remaining elements, and true\false marker of availability previous elements.
// If no more elements then ok==false.
func (v *Vector[T]) Last() (SliceIter[T], T, bool) {
	var (
		iterator  = NewTail(v, v.DeleteActualOne)
		first, ok = iterator.Prev()
	)
	return iterator, first, ok
}

// Slice collects the elements to a slice
func (v *Vector[T]) Slice() (out []T) {
	if v == nil {
		return out
	}
	return slice.Clone(*v)
}

// Append collects the values to the specified 'out' slice
func (v *Vector[T]) Append(out []T) []T {
	if v == nil {
		return out
	}
	return append(out, (*v)...)
}

// Clone just makes a copy of the vector instance
func (v *Vector[T]) Clone() *Vector[T] {
	return WrapVector(slice.Clone(*v))
}

// IsEmpty returns true if the collection is empty
func (v *Vector[T]) IsEmpty() bool {
	return v.Len() == 0
}

// Len returns amount of elements
func (v *Vector[T]) Len() int {
	if v == nil {
		return 0
	}
	return notsafe.GetLen(*v)
}

// Track applies tracker to elements with error checking. Return the c.ErrBreak to stop tracking.
func (v *Vector[T]) Track(tracker func(int, T) error) error {
	if v == nil {
		return nil
	}
	return slice.Track(*v, tracker)
}

// TrackEach applies tracker to elements without error checking
func (v *Vector[T]) TrackEach(tracker func(int, T)) {
	if v != nil {
		slice.TrackEach(*v, tracker)
	}
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (v *Vector[T]) For(walker func(T) error) error {
	if v == nil {
		return nil
	}
	return slice.For(*v, walker)
}

// ForEach applies walker to elements without error checking
func (v *Vector[T]) ForEach(walker func(T)) {
	if !(v == nil) {
		slice.ForEach(*v, walker)
	}
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the vector len, returns zero T and false in the second result
func (v *Vector[T]) Get(index int) (t T, ok bool) {
	if v == nil {
		return
	}
	return slice.Gett(*v, index)
}

// Add adds elements to the end of the vector
func (v *Vector[T]) Add(elements ...T) {
	if v != nil {
		*v = append(*v, elements...)
	}
}

// AddOne adds an element to the end of the vector
func (v *Vector[T]) AddOne(element T) {
	if v != nil {
		*v = append(*v, element)
	}
}

// AddAll inserts all elements from the "other" collection
func (v *Vector[T]) AddAll(other c.ForEachLoop[T]) {
	if v != nil {
		other.ForEach(func(element T) { *v = append(*v, element) })
	}
}

// DeleteOne removes an element by the index
func (v *Vector[T]) DeleteOne(index int) {
	_ = v.DeleteActualOne(index)
}

// DeleteActualOne removes an element by the index
func (v *Vector[T]) DeleteActualOne(index int) bool {
	_, ok := v.Remove(index)
	return ok
}

// Remove removes and returns an element by the index
func (v *Vector[T]) Remove(index int) (t T, ok bool) {
	if v == nil {
		return t, ok
	}
	if e := *v; index >= 0 && index < len(e) {
		de := e[index]
		*v = slice.Delete(e, index)
		return de, true
	}
	return t, ok
}

// Delete drops elements by indexes
func (v *Vector[T]) Delete(indexes ...int) {
	v.DeleteActual(indexes...)
}

// DeleteActual drops elements by indexes with verification of no-op
func (v *Vector[T]) DeleteActual(indexes ...int) bool {
	if v == nil {
		return false
	}
	l := len(indexes)
	if l == 0 {
		return false
	} else if l == 1 {
		return v.DeleteActualOne(indexes[0])
	}

	e := *v
	el := len(e)

	sort.Ints(indexes)

	shift := 0
	for i := 0; i < l; i++ {
		index := indexes[i] - shift
		delAmount := 1
		if index >= 0 && index < el {
			curIndex := index
			for i < l-1 {
				nextIndex := indexes[i+1]
				if nextIndex-curIndex == 1 {
					delAmount++
					i++
					curIndex = nextIndex
				} else {
					break
				}
			}

			e = append(e[0:index], e[index+delAmount:]...)
			shift += delAmount
		}
	}
	if shift > 0 {
		*v = e
		return true
	}
	return false
}

// Set puts an element into the vector at the index
func (v *Vector[T]) Set(index int, value T) {
	v.SetNew(index, value)
}

// SetNew puts an element into the vector at the index
func (v *Vector[T]) SetNew(index int, value T) bool {
	if v == nil {
		return false
	}
	e := *v
	if index < 0 {
		return false
	}
	l := len(e)
	if index >= l {
		c := l * 2
		l := index + 1
		if l > c {
			c = l
		}
		ne := make([]T, l, c)
		copy(ne, e)
		e = ne
		*v = e
	}
	e[index] = value
	return true
}

// Filter returns a stream consisting of vector elements matching the filter
func (v *Vector[T]) Filter(filter func(T) bool) stream.Iter[T] {
	h := v.Head()
	return stream.New(loop.Filter(h.Next, filter).Next)
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (v *Vector[T]) Filt(predicate func(T) (bool, error)) breakStream.Iter[T] {
	h := v.Head()
	return breakStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (v *Vector[T]) Convert(converter func(T) T) stream.Iter[T] {
	return collection.Convert(v, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (v *Vector[T]) Conv(converter func(T) (T, error)) breakStream.Iter[T] {
	return collection.Conv(v, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (v *Vector[T]) Reduce(merge func(T, T) T) (out T) {
	if v != nil {
		out = slice.Reduce(*v, merge)
	}
	return out
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (v *Vector[T]) HasAny(predicate func(T) bool) (ok bool) {
	if v != nil {
		ok = slice.HasAny(*v, predicate)
	}
	return ok
}

// Sort sorts the Vector in-place and returns it
func (v *Vector[T]) Sort(comparer slice.Comparer[T]) *Vector[T] {
	return v.sortBy(slice.Sort, comparer)
}

// StableSort stable sorts the Vector in-place and returns it
func (v *Vector[T]) StableSort(comparer slice.Comparer[T]) *Vector[T] {
	return v.sortBy(slice.StableSort, comparer)
}

func (v *Vector[T]) sortBy(sorter func([]T, slice.Comparer[T]) []T, comparer slice.Comparer[T]) *Vector[T] {
	if v != nil {
		sorter(*v, comparer)
	}
	return v
}

// String returns then string representation
func (v *Vector[T]) String() string {
	if v == nil {
		return ""
	}
	return slice.ToString(*v)
}
