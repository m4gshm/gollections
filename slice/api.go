//Package slice provides generic functions for slice types
package slice

import (
	"constraints"
	"fmt"

	"github.com/m4gshm/gollections/typ"
)

//Of transforms elements to the slice of them
func Of[T any](elements ...T) []T {
	return elements
}

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
func Group[T any, K comparable](elements []T, by typ.Converter[T, K]) map[K][]T {
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

//Track applies tracker to elements with error checking
func Track[T any](elements []T, tracker func(int, T) error) error {
	for i, e := range elements {
		if err := tracker(i, e); err != nil {
			return err
		}
	}
	return nil
}

//Track applies tracker to elements without error checking
func TrackEach[T any](elements []T, tracker func(int, T)) {
	for i, e := range elements {
		tracker(i, e)
	}
}

//For applies walker to elements with error checking
func For[T any](elements []T, walker func(T) error) error {
	for _, e := range elements {
		if err := walker(e); err != nil {
			return err
		}
	}
	return nil
}

//ForEach applies walker to elements without error checking
func ForEach[T any](elements []T, walker func(T)) {
	for _, e := range elements {
		walker(e)
	}
}

//ForRefs applies walker to references with error checking
func ForRefs[T any](references []*T, walker func(T) error) error {
	for _, e := range references {
		if err := walker(*e); err != nil {
			return err
		}
	}
	return nil
}

//ForRefs applies walker to references without error checking
func ForEachRef[T any](references []*T, walker func(T)) {
	for _, e := range references {
		walker(*e)
	}
}

//ToString converts elements to the string representation
func ToString[T any](elements []T) string {
	str := ""
	for _, v := range elements {
		if len(str) > 0 {
			str += ", "
		}
		str += fmt.Sprintf("%+v", v)
	}
	return "[" + str + "]"
}

//ToString converts references to the string representation
func ToStringRefs[T any](references []*T) string {
	str := ""
	for _, ref := range references {
		if len(str) > 0 {
			str += ", "
		}
		if ref == nil {
			str += "nil"
		} else {
			str += fmt.Sprintf("%+v", *ref)
		}
	}
	return "[" + str + "]"
}
