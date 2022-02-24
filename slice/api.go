//Package slice provides generic functions for slice types
package slice

import (
	"bytes"
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

//ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = it.ErrBreak

//Of is generic sclie constructor
func Of[T any](elements ...T) []T { return elements }

//Copy makes the new slice with copied elements.
func Copy[T any](elements []T) []T {
	copied := make([]T, len(elements))
	copy(copied, elements)
	return copied
}

//Delete removes an element by index from a slice
func Delete[T any](index int, elements []T) []T {
	return append(elements[0:index], elements[index+1:]...)
}

//Group converts elements into the map containing slices of the elements separated by keys, which are retrieved using a Converter object.
func Group[T any, K comparable](elements []T, by c.Converter[T, K]) map[K][]T {
	groups := map[K][]T{}
	for _, e := range elements {
		key := by(e)
		group := groups[key]
		if group == nil {
			group = make([]T, 0)
		}
		groups[key] = append(group, e)

	}
	return groups
}

//Range generates the sclie of integers in the range defined by from and to inclusive.
func Range[T constraints.Integer](from T, to T) []T {
	if to == from {
		return []T{to}
	}
	amount := to - from
	delta := 1
	if amount < 0 {
		amount = -amount
		delta = -1
	}
	amount = amount + 1
	elements := make([]T, amount)
	e := from
	for i := 0; i < int(amount); i++ {
		elements[i] = e
		e = e + T(delta)
	}
	return elements
}

//Reverse inverts slice elements
func Reverse[T any](elements []T) []T {
	l := 0
	h := len(elements) - 1
	for l < h {
		le, he := elements[l], elements[h]
		elements[l], elements[h] = he, le
		l++
		h--
	}
	return elements
}

//Sort sorts a slice and returns the one.
func Sort[T any](elements []T, less func(e1, e2 T) bool) []T {
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	return elements
}

//SortByOrdered sorts elements by a converter that thransforms a element to an Ordered (int, string and so on).
func SortByOrdered[O any, o constraints.Ordered](elements []O, by c.Converter[O, o]) []O {
	return Sort(elements, func(e1, e2 O) bool { return by(e1) < by(e2) })
}

//SortCopy sorts copied slices
func SortCopy[T any](elements []T, less func(e1, e2 T) bool) []T {
	var dest = make([]T, len(elements))
	copy(dest, elements)
	return Sort(dest, less)
}

//Get returns the element by its index in elements, otherwise, if the provided index is ouf of the elements, returns zero T and false in the second result
func Get[T any](elements []T, index int) (T, bool) {
	l := len(elements)
	if l > 0 && (index >= 0 || index < l) {
		return elements[index], true
	}
	var no T
	return no, false
}

//Track applies tracker to elements with error checking. To stop traking just return the ErrBreak.
func Track[T any](elements []T, tracker func(int, T) error) error {
	for i, e := range elements {
		if err := tracker(i, e); err != ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//TrackEach applies tracker to elements without error checking.
func TrackEach[T any](elements []T, tracker func(int, T)) {
	for i, e := range elements {
		tracker(i, e)
	}
}

//For applies a walker to elements of an slice. To stop walking just return the ErrBreak.
func For[T any](elements []T, walker func(T) error) error {
	for _, e := range elements {
		if err := walker(e); err != ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//ForEach applies walker to elements without error checking.
func ForEach[T any](elements []T, walker func(T)) {
	for _, e := range elements {
		walker(e)
	}
}

//ForEachRef applies walker to references without error checking
func ForEachRef[T any](references []*T, walker func(T)) {
	for _, e := range references {
		walker(*e)
	}
}

//ToString converts elements to their default string representation
func ToString[T any](elements []T) string {
	return ToStringf(elements, "%+v", " ")
}

//ToStringf converts elements to a string representation defined by a custom element format and a delimiter
func ToStringf[T any](elements []T, elementFormat, delimeter string) string {
	str := bytes.Buffer{}
	str.WriteString("[")
	for i, v := range elements {
		if i > 0 {
			_, _ = str.WriteString(delimeter)
		}
		str.WriteString(fmt.Sprintf(elementFormat, v))
	}
	str.WriteString("]")
	return str.String()
}

//ToStringRefs converts references to the default string representation
func ToStringRefs[T any](references []*T) string {
	return ToStringRefsf(references, "%+v", "nil", " ")
}

//ToStringRefsf converts references to a string representation defined by a custom delimiter and a nil value representation
func ToStringRefsf[T any](references []*T, elementFormat, nilValue, delimeter string) string {
	str := bytes.Buffer{}
	str.WriteString("[")
	for i, ref := range references {
		if i > 0 {
			_, _ = str.WriteString(delimeter)
		}
		if ref == nil {
			str.WriteString(nilValue)
		} else {
			str.WriteString(fmt.Sprintf(elementFormat, *ref))
		}
	}
	str.WriteString("]")
	return str.String()
}
