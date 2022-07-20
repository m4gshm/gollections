//Package slice provides generic functions for slice types
package slice

import (
	"bytes"
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/it/slice"
	"github.com/m4gshm/gollections/op"
)

//ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = it.ErrBreak

//Of is generic sclie constructor
func Of[T any](elements ...T) []T { return elements }

//Clone makes new slice instance with copied elements.
func Clone[T any, TS ~[]T](elements TS) []T {
	copied := make([]T, len(elements))
	copy(copied, elements)
	return copied
}

//Delete removes an element by index from the slice
func Delete[T any, TS ~[]T](index int, elements TS) []T {
	return append(elements[0:index], elements[index+1:]...)
}

//Group converts the slice into a map with keys computeable by a Converter.
func Group[T any, K comparable, TS ~[]T](elements TS, by c.Converter[T, K]) map[K][]T {
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

//Filter creates a slice containing only the filtered elements
func Filter[T any, TS ~[]T](elements TS, filter c.Predicate[T]) []T {
	return it.Slice[T](slice.Filter(elements, filter))
}

//Flatt unfolds the n-dimensional slice into a n-1 dimensional slice
func Flatt[From, To any, FS ~[]From](elements FS, by c.Flatter[From, To]) []To {
	return it.Slice[To](slice.Flatt(elements, by))
}

//Range generates a sclice of integers in the range defined by from and to inclusive.
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

//Reverse inverts elements
func Reverse[T any, TS ~[]T](elements TS) []T {
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

//Sort returns sorted elements
func Sort[T any, TS ~[]T](elements TS, less func(e1, e2 T) bool) []T {
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	return elements
}

//SortByOrdered sorts elements by converting them to Ordered values and applying the operator <
func SortByOrdered[T any, o constraints.Ordered, TS ~[]T](elements TS, by c.Converter[T, o]) []T {
	return Sort(elements, func(e1, e2 T) bool { return by(e1) < by(e2) })
}

//Reduce reduces elements to an one.
func Reduce[T any, TS ~[]T](elements TS, by op.Binary[T]) T {
	var result T
	for i, v := range elements {
		if i == 0 {
			result = v
		} else {
			result = by(result, v)
		}
	}
	return result
}

//Get returns an element from elements by index, otherwise, if the provided index is ouf of the elements, returns zero T and false in the second result
func Get[T any, TS ~[]T](elements TS, index int) (T, bool) {
	l := len(elements)
	if l > 0 && (index >= 0 || index < l) {
		return elements[index], true
	}
	var no T
	return no, false
}

//Track applies tracker to elements with error checking. To stop traking just return the ErrBreak.
func Track[T any, TS ~[]T](elements TS, tracker func(int, T) error) error {
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
func TrackEach[T any, TS ~[]T](elements TS, tracker func(int, T)) {
	for i, e := range elements {
		tracker(i, e)
	}
}

//For applies walker to elements. To stop walking just return the ErrBreak.
func For[T any, TS ~[]T](elements TS, walker func(T) error) error {
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
func ForEach[T any, TS ~[]T](elements TS, walker func(T)) {
	for _, e := range elements {
		walker(e)
	}
}

//ForEachRef applies walker to references without error checking
func ForEachRef[T any, TS ~[]*T](references TS, walker func(T)) {
	for _, e := range references {
		walker(*e)
	}
}

//ToString converts elements to their default string representation
func ToString[T any, TS ~[]T](elements TS) string {
	return ToStringf(elements, "%+v", " ")
}

//ToStringf converts elements to a string representation defined by a custom element format and a delimiter
func ToStringf[T any, TS ~[]T](elements TS, elementFormat, delimeter string) string {
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
func ToStringRefs[T any, TS ~[]*T](references TS) string {
	return ToStringRefsf(references, "%+v", "nil", " ")
}

//ToStringRefsf converts references to a string representation defined by a delimiter and a nil value representation
func ToStringRefsf[T any, TS ~[]*T](references TS, elementFormat, nilValue, delimeter string) string {
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
