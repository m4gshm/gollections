package mutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
)

// NewVectorCap instantiates Vector with a predefined capacity
func NewVectorCap[T any](capacity int) *Vector[T] {
	return WrapVector(make([]T, 0, capacity))
}

// NewVector instantiates Vector based on copy of elements slice
func NewVector[T any](elements []T) *Vector[T] {
	return WrapVector(slice.Clone(elements))
}

// WrapVector instantiates Vector using a slise as internal storage
func WrapVector[T any](elements []T) *Vector[T] {
	v := Vector[T](elements)
	return &v
}

// Vector is the Collection implementation that provides elements order and index access
type Vector[T any] []T

var (
	_ c.Addable[any]          = (*Vector[any])(nil)
	_ c.AddableAll[any]       = (*Vector[any])(nil)
	_ c.Deleteable[int]       = (*Vector[any])(nil)
	_ c.DeleteableVerify[int] = (*Vector[any])(nil)
	_ c.Settable[int, any]    = (*Vector[any])(nil)
	_ c.SettableNew[int, any] = (*Vector[any])(nil)
	_ c.Vector[any]           = (*Vector[any])(nil)
	_ fmt.Stringer            = (*Vector[any])(nil)
)

// Begin creates an iterator of the vector
func (v *Vector[T]) Begin() c.Iterator[T] {
	h := v.Head()
	return &h
}

// BeginEdit creates an iterator with deleting elements
func (v *Vector[T]) BeginEdit() c.DelIterator[T] {
	h := v.Head()
	return &h
}

// Head creates an iterator impl instace of the vector
func (v *Vector[T]) Head() Iter[Vector[T], T] {
	return NewHead(v, v.DeleteActualOne)
}

// Tail creates an iterator pointing to the end of the vector
func (v *Vector[T]) Tail() Iter[Vector[T], T] {
	return NewTail(v, v.DeleteActualOne)
}

// First returns the first element of the vector, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then returns false in the last value.
func (v *Vector[T]) First() (Iter[Vector[T], T], T, bool) {
	var (
		iterator  = NewHead(v, v.DeleteActualOne)
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Last returns the last element of the vector, an iterator to iterate over the remaining elements, and true\false marker of availability prev elements.
// If no more elements then returns false in the last value.
func (v *Vector[T]) Last() (Iter[Vector[T], T], T, bool) {
	var (
		iterator  = NewTail(v, v.DeleteActualOne)
		first, ok = iterator.Prev()
	)
	return iterator, first, ok
}

// Slice transforms the vector to a slice
func (v *Vector[T]) Slice() []T {
	return slice.Clone(*v)
}

// Copy just makes a copy of the vector instance
func (v *Vector[T]) Copy() *Vector[T] {
	return WrapVector(slice.Clone(*v))
}

// IsEmpty checks if there are elements in the vector
func (v *Vector[T]) IsEmpty() bool {
	return v.Len() == 0
}

// Len returns amount of elements
func (v *Vector[T]) Len() int {
	return notsafe.GetLen(*v)
}

// Track applies tracker to elements with error checking. To stop traking just return the ErrBreak
func (v *Vector[T]) Track(tracker func(int, T) error) error {
	return slice.Track(*v, tracker)
}

// TrackEach applies tracker to elements without error checking
func (v *Vector[T]) TrackEach(tracker func(int, T)) {
	slice.TrackEach(*v, tracker)
}

// For applies walker to elements. To stop walking just return the ErrBreak
func (v *Vector[T]) For(walker func(T) error) error {
	return slice.For(*v, walker)
}

// ForEach applies walker to elements without error checking
func (v *Vector[T]) ForEach(walker func(T)) {
	slice.ForEach(*v, walker)
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the vector len, returns zero T and false in the second result
func (v *Vector[T]) Get(index int) (T, bool) {
	return slice.Get(*v, index)
}

// Add adds elements to the end of the vector
func (v *Vector[T]) Add(elements ...T) {
	*v = append(*v, elements...)
}

// AddOne adds a element to the end of the vector
func (v *Vector[T]) AddOne(element T) {
	*v = append(*v, element)
}

func (v *Vector[T]) AddAll(elements c.Iterable[T]) {
	*v = append(*v, loop.ToSlice(elements.Begin().Next)...)
}

func (v *Vector[T]) AddAllNew(elements c.Iterator[T]) {
	*v = append(*v, loop.ToSlice(elements.Next)...)
}

// Delete removes a element by the index
func (v *Vector[T]) DeleteOne(index int) {
	_ = v.DeleteActualOne(index)
}

// DeleteActualOne removes a element by the index
func (v *Vector[T]) DeleteActualOne(index int) bool {
	_, ok := v.Remove(index)
	return ok
}

// Remove removes and returns a element by the index
func (v *Vector[T]) Remove(index int) (T, bool) {
	if e := *v; index >= 0 && index < len(e) {
		de := e[index]
		*v = slice.Delete(index, e)
		return de, true
	}
	var no T
	return no, false
}

// Delete drops elements by indexes
func (v *Vector[T]) Delete(indexes ...int) {
	v.DeleteActual(indexes...)
}

// DeleteActual drops elements by indexes with verification of no-op
func (v *Vector[T]) DeleteActual(indexes ...int) bool {
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

// Set puts a element into the vector at the index
func (v *Vector[T]) Set(index int, value T) {
	v.SetNew(index, value)
}

// SetNew puts a element into the vector at the index
func (v *Vector[T]) SetNew(index int, value T) bool {
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

// Filter returns a pipe consisting of vector elements matching the filter
func (v *Vector[T]) Filter(filter func(T) bool) c.Pipe[T] {
	h := v.Head()
	return iter.NewPipe[T](iter.Filter(h, h.Next, filter))
}

// Map returns a pipe of converted vector elements by the converter 'by'
func (v *Vector[T]) Convert(by func(T) T) c.Pipe[T] {
	h := v.Head()
	return iter.NewPipe[T](iter.Convert(h, h.Next, by))
}

// Reduce reduces elements to an one
func (v *Vector[T]) Reduce(by func(T, T) T) T {
	h := v.Head()
	return loop.Reduce(h.Next, by)
}

// Sort sorts the Vector in-place and returns it
func (v *Vector[T]) Sort(less slice.Less[T]) *Vector[T] {
	return v.sortBy(sort.Slice, less)
}

func (v *Vector[T]) StableSort(less slice.Less[T]) *Vector[T] {
	return v.sortBy(sort.SliceStable, less)
}

func (v *Vector[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) *Vector[T] {
	slice.Sort(*v, sorter, less)
	return v
}

// String returns then string representation
func (v *Vector[T]) String() string {
	return slice.ToString(*v)
}
