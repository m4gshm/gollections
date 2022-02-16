package mutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewVector[T any](capacity int) *Vector[T] {
	return WrapVector(make([]T, 0, capacity))
}

func ToVector[T any](elements []T) *Vector[T] {
	return WrapVector(slice.Copy(elements))
}

func WrapVector[T any](elements []T) *Vector[T] {
	return &Vector[T]{elements: elements}
}

//Vector stores ordered elements, provides index access.
type Vector[T any] struct {
	elements []T
}

var (
	_ Addable[any]       = (*Vector[any])(nil)
	_ Deleteable[int]    = (*Vector[any])(nil)
	_ Settable[int, any] = (*Vector[any])(nil)
	_ c.Vector[any]      = (*Vector[any])(nil)
	_ fmt.Stringer       = (*Vector[any])(nil)
)

func (v *Vector[T]) Begin() c.Iterator[T] {
	return v.Head()
}

func (v *Vector[T]) BeginEdit() Iterator[T] {
	return v.Head()
}

func (v *Vector[T]) Head() *Iter[T] {
	return NewHead(&v.elements, v.DeleteOne)
}

func (v *Vector[T]) Tail() *Iter[T] {
	return NewTail(&v.elements, v.DeleteOne)
}

func (v *Vector[T]) Collect() []T {
	return slice.Copy(v.elements)
}

func (v *Vector[T]) Copy() *Vector[T] {
	return WrapVector(slice.Copy(v.elements))
}

func (v *Vector[T]) IsEmpty() bool {
	return v.Len() == 0
}

func (v *Vector[T]) Len() int {
	return it.GetLen(v.elements)
}

func (v *Vector[T]) Track(tracker func(int, T) error) error {
	return slice.Track(v.elements, tracker)
}

func (v *Vector[T]) TrackEach(tracker func(int, T)) {
	slice.TrackEach(v.elements, tracker)
}

func (v *Vector[T]) For(walker func(T) error) error {
	return slice.For(v.elements, walker)
}

func (v *Vector[T]) ForEach(walker func(T)) {
	slice.ForEach(v.elements, walker)
}

func (v *Vector[T]) Get(index int) (T, bool) {
	return slice.Get(v.elements, index)
}

func (v *Vector[T]) Add(elements ...T) bool {
	return v.AddAll(elements)
}

func (v *Vector[T]) AddAll(elements []T) bool {
	v.elements = append(v.elements, elements...)
	return true
}

func (v *Vector[T]) AddOne(element T) bool {
	v.elements = append(v.elements, element)
	return true
}

func (v *Vector[T]) DeleteOne(index int) bool {
	_, ok := v.Remove(index)
	return ok
}

func (v *Vector[T]) Remove(index int) (T, bool) {
	if e := v.elements; index >= 0 && index < len(e) {
		de := e[index]
		v.elements = slice.Delete(index, e)
		return de, true
	}
	var no T
	return no, false
}

func (v *Vector[T]) Delete(indexes ...int) bool {
	l := len(indexes)
	if l == 0 {
		return false
	} else if l == 1 {
		return v.DeleteOne(indexes[0])
	}

	e := v.elements
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
		v.elements = e
		return true
	}
	return false
}

func (v *Vector[T]) Set(index int, value T) bool {
	e := v.elements
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
		v.elements = e
	}
	e[index] = value
	return true
}

func (v *Vector[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Filter(v.Head(), filter))
}

func (v *Vector[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Map(v.Head(), by))
}

func (v *Vector[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(v.Head(), by)
}

//Sotr sorts the Vector in-place and returns it.
func (v *Vector[t]) Sort(less func(e1, e2 t) bool) *Vector[t] {
	v.elements = slice.Sort(v.elements, less)
	return v
}

//String returns then string representation.
func (v *Vector[T]) String() string {
	return slice.ToString(v.elements)
}
