// Package slice provides generic functions for slice types
package slice

import (
	"bytes"
	"fmt"
	"sort"
	"unsafe"

	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = it.ErrBreak

// Of is generic slice constructor
func Of[T any](elements ...T) []T { return elements }

// OfLoop builds a slice by iterating elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfLoop[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) ([]T, error) {
	r := []T{}
	for hasNext(source) {
		o, err := getNext(source)
		if err != nil {
			return r, err
		}
		r = append(r, o)
	}
	return r, nil
}

// OfLoop builds a slice by an generator function.
// The generator returns an element, or false if the generation is over, or an error.
func Generate[T any](next func() (T, bool, error)) ([]T, error) {
	r := []T{}
	for  {
		e, ok, err := next()
		if err != nil || !ok {
			return r, err
		}
		r = append(r, e)
	}
	return r, nil
}

// Clone makes new slice instance with copied elements.
func Clone[T any, TS ~[]T](elements TS) []T {
	copied := make([]T, len(elements))
	copy(copied, elements)
	return copied
}

// Delete removes an element by index from the slice 'elements'
func Delete[T any, TS ~[]T](index int, elements TS) []T {
	return append(elements[0:index], elements[index+1:]...)
}

// Group converts the slice into a map with keys computeable by the converter 'by'
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

// Map creates a slice consisting of the transformed elements using the converter 'by'
func Map[From, To any, FS ~[]From](elements FS, by c.Converter[From, To]) []To {
	result := make([]To, len(elements))
	for i, e := range elements {
		result[i] = by(e)
	}
	return result
}

// MapFit additionally filters 'From' elements.
func MapFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], by c.Converter[From, To]) []To {
	result := make([]To, 0)
	for _, e := range elements {
		if fit(e) {
			result = append(result, by(e))
		}
	}
	return result
}

// Flatt unfolds the n-dimensional slice into a n-1 dimensional slice
func Flatt[From, To any, FS ~[]From](elements FS, by c.Flatter[From, To]) []To {
	result := make([]To, 0)
	for _, e := range elements {
		result = append(result, by(e)...)

	}
	return result
}

// FlattFit additionally filters 'From' elements.
func FlattFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], by c.Flatter[From, To]) []To {
	result := make([]To, 0)
	for _, e := range elements {
		if fit(e) {
			result = append(result, by(e)...)
		}
	}
	return result
}

// FlattElemFit unfolds the n-dimensional slice into a n-1 dimensional slice with additinal filtering of 'To' elements.
func FlattElemFit[From, To any, FS ~[]From](elements FS, by c.Flatter[From, To], fit c.Predicate[To]) []To {
	result := make([]To, 0)
	for _, e := range elements {
		for _, to := range by(e) {
			if fit(to) {
				result = append(result, to)
			}
		}
	}
	return result
}

// FlattFitFit unfolds the n-dimensional slice 'elements' into a n-1 dimensional slice with additinal filtering of 'From' and 'To' elements.
func FlattFitFit[From, To any, FS ~[]From](elements FS, fitFrom c.Predicate[From], by c.Flatter[From, To], fitTo c.Predicate[To]) []To {
	result := make([]To, 0)
	for _, e := range elements {
		if fitFrom(e) {
			for _, to := range by(e) {
				if fitTo(to) {
					result = append(result, to)
				}
			}
		}
	}
	return result
}

// Filter creates a slice containing only the filtered elements
func Filter[T any, TS ~[]T](elements TS, filter c.Predicate[T]) []T {
	result := make([]T, 0)
	for _, e := range elements {
		if filter(e) {
			result = append(result, e)
		}
	}
	return result
}

// Range generates a sclice of integers in the range defined by from and to inclusive.
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

// Reverse inverts elements order
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

// Sort sorts elements in place by applying the function 'less'
func Sort[T any, TS ~[]T](elements TS, less func(e1, e2 T) bool) []T {
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	return elements
}

// SortByOrdered sorts elements in place by converting them to Ordered values and applying the operator <
func SortByOrdered[T any, o constraints.Ordered, TS ~[]T](elements TS, by c.Converter[T, o]) []T {
	return Sort(elements, func(e1, e2 T) bool { return by(e1) < by(e2) })
}

// Reduce reduces elements to an one
func Reduce[T any, TS ~[]T](elements TS, by c.Binary[T]) T {
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

// Sum returns the sum of all elements
func Sum[T c.Summable, TS ~[]T](elements TS) T {
	return Reduce(elements, op.Sum[T])
}

// First returns the first element that satisfies requirements of the predicate 'fit'
func First[T any, TS ~[]T](elements TS, by c.Predicate[T]) (T, bool) {
	for _, e := range elements {
		if by(e) {
			return e, true
		}
	}
	var no T
	return no, false
}

// Last returns the latest element that satisfies requirements of the predicate 'fit'
func Last[T any, TS ~[]T](elements TS, by c.Predicate[T]) (T, bool) {
	for i := len(elements) - 1; i >= 0; i-- {
		e := elements[i]
		if by(e) {
			return e, true
		}
	}
	var no T
	return no, false
}

// Get returns an element from the elements by index, otherwise, if the provided index is ouf of the elements, returns zero T and false in the second result
func Get[T any, TS ~[]T](elements TS, index int) (T, bool) {
	l := len(elements)
	if l > 0 && (index >= 0 || index < l) {
		return elements[index], true
	}
	var no T
	return no, false
}

// Track applies tracker to elements with error checking. To stop traking just return the ErrBreak
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

// TrackEach applies tracker to elements without error checking
func TrackEach[T any, TS ~[]T](elements TS, tracker func(int, T)) {
	for i, e := range elements {
		tracker(i, e)
	}
}

// For applies walker to elements. To stop walking just return the ErrBreak
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

// ForEach applies walker to elements without error checking
func ForEach[T any, TS ~[]T](elements TS, walker func(T)) {
	for _, e := range elements {
		walker(e)
	}
}

// ForEachRef applies walker to references without error checking
func ForEachRef[T any, TS ~[]*T](references TS, walker func(T)) {
	for _, e := range references {
		walker(*e)
	}
}

// ToString converts elements to their default string representation
func ToString[T any, TS ~[]T](elements TS) string {
	return ToStringf(elements, "%+v", " ")
}

// ToStringf converts elements to a string representation defined by a custom element format and a delimiter
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

// ToStringRefs converts references to the default string representation
func ToStringRefs[T any, TS ~[]*T](references TS) string {
	return ToStringRefsf(references, "%+v", "nil", " ")
}

// ToStringRefsf converts references to a string representation defined by a delimiter and a nil value representation
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

func BehaveAsStrings[T ~string, TS ~[]T](elements TS) []string {
	ptr := unsafe.Pointer(&elements)
	s := *(*[]string)(ptr)
	return s
}

func StringsBehaveAs[TS ~[]T, T ~string](elements []string) TS {
	ptr := unsafe.Pointer(&elements)
	s := *(*TS)(ptr)
	return s
}
