package mutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
)

// WrapVector instantiates Vector using a slise as internal storage
func WrapVector[T any](elements []T) *Vector[T] {
	v := Vector[T](elements)
	return &v
}

// Vector is a collection implementation that provides elements order and index access
type Vector[T any] []T

var (
	_ c.Addable[any]             = (*Vector[any])(nil)
	_ c.AddableAll[seq.Seq[any]] = (*Vector[any])(nil)
	_ c.Deleteable[int]          = (*Vector[any])(nil)
	_ c.DeleteableVerify[int]    = (*Vector[any])(nil)
	_ c.Settable[int, any]       = (*Vector[any])(nil)
	_ c.SettableNew[int, any]    = (*Vector[any])(nil)
	_ c.OrderedRange[any]        = (*Vector[any])(nil)
	_ collection.Vector[any]     = (*Vector[any])(nil)
	_ fmt.Stringer               = (*Vector[any])(nil)
)

// All is used to iterate through the collection using `for e := range`.
func (v *Vector[T]) All(consumer func(T) bool) {
	if v != nil {
		slice.WalkWhile(*v, consumer)
	}
}

// IAll is used to iterate through the collection using `for i, e := range`.
func (v *Vector[T]) IAll(consumer func(int, T) bool) {
	if v != nil {
		slice.TrackWhile(*v, consumer)
	}
}

// Head returns the first element.
func (v *Vector[T]) Head() (t T, ok bool) {
	if v == nil {
		return t, false
	}
	return collection.Head(v)
}

// Tail returns the latest element
func (v *Vector[T]) Tail() (t T, ok bool) {
	if v == nil {
		return t, false
	}
	return slice.Tail(*v)
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
	return collection.IsEmpty(v)
}

// Len returns amount of elements
func (v *Vector[T]) Len() int {
	if v == nil {
		return 0
	}
	return notsafe.GetLen(*v)
}

// TrackEach applies consumer to elements without error checking
func (v *Vector[T]) TrackEach(consumer func(int, T)) {
	if v != nil {
		slice.TrackEach(*v, consumer)
	}
}

// ForEach applies consumer to elements without error checking
func (v *Vector[T]) ForEach(consumer func(T)) {
	if !(v == nil) {
		slice.ForEach(*v, consumer)
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
func (v *Vector[T]) AddAll(other seq.Seq[T]) {
	if v != nil {
		seq.ForEach(other, func(element T) { *v = append(*v, element) })
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
	switch l {
	case 0:
		return false
	case 1:
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

// Filter returns a seq consisting of vector elements matching the filter
func (v *Vector[T]) Filter(filter func(T) bool) seq.Seq[T] {
	return collection.Filter(v, filter)
}

// Filt returns an errorable seq consisting of elements that satisfy the condition of the 'filter' function
func (v *Vector[T]) Filt(filter func(T) (bool, error)) seq.SeqE[T] {
	return collection.Filt(v, filter)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (v *Vector[T]) Convert(converter func(T) T) seq.Seq[T] {
	return collection.Convert(v, converter)
}

// Conv returns an errorable seq that applies the 'converter' function to the collection elements
func (v *Vector[T]) Conv(converter func(T) (T, error)) seq.SeqE[T] {
	return collection.Conv(v, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (v *Vector[T]) Reduce(merge func(T, T) T) (out T) {
	if v != nil {
		out = slice.Reduce(*v, merge)
	}
	return out
}

// HasAny checks whether the vector contains an element that satisfies the condition.
func (v *Vector[T]) HasAny(condition func(T) bool) (ok bool) {
	if v != nil {
		ok = slice.HasAny(*v, condition)
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
