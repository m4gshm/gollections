// Package slice provides generic functions for slice types
package slice

import (
	"bytes"
	"fmt"
	"unsafe"

	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = loop.ErrBreak

// Of is generic slice constructor
func Of[T any](elements ...T) []T { return elements }

// OfLoop builds a slice by iterating elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfLoop[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) ([]T, error) {
	var r []T
	for hasNext(source) {
		o, err := getNext(source)
		if err != nil {
			return r, err
		}
		r = append(r, o)
	}
	return r, nil
}

// Generate builds a slice by an generator function.
// The generator returns an element, or false if the generation is over, or an error.
func Generate[T any](next func() (T, bool, error)) ([]T, error) {
	var r []T
	for {
		e, ok, err := next()
		if err != nil || !ok {
			return r, err
		}
		r = append(r, e)
	}
}

// Clone makes new slice instance with copied elements
func Clone[TS ~[]T, T any](elements TS) TS {
	if elements == nil {
		return nil
	}
	copied := make(TS, len(elements))
	copy(copied, elements)
	return copied
}

// DeepClone copies slice elements using a copier function and returns them as a new slice
func DeepClone[TS ~[]T, T any](elements TS, copier func(T) T) TS {
	return Convert(elements, copier)
}

// Delete removes an element by index from the slice 'elements'
func Delete[TS ~[]T, T any](index int, elements TS) TS {
	return append(elements[0:index], elements[index+1:]...)
}

// Group converts a slice into a map, extracting a key for each element of the slice applying the converter 'keyProducer'
func Group[T any, K comparable, TS ~[]T](elements TS, keyProducer func(T) K) map[K]TS {
	groups := map[K]TS{}
	for _, e := range elements {
		initGroup(keyProducer(e), e, groups)
	}
	return groups
}

// GroupInMultiple converts a slice into a map, extracting multiple keys per each element of the slice applying the converter 'keyProducer'
func GroupInMultiple[T any, K comparable, TS ~[]T](elements TS, keysProducer func(T) []K) map[K]TS {
	groups := map[K]TS{}
	for _, e := range elements {
		if keys := keysProducer(e); len(keys) == 0 {
			var key K
			initGroup(key, e, groups)
		} else {
			for _, key := range keys {
				initGroup(key, e, groups)
			}
		}
	}
	return groups
}

func initGroup[T any, K comparable, TS ~[]T](key K, e T, groups map[K]TS) {
	groups[key] = append(groups[key], e)
}

// Convert creates a slice consisting of the transformed elements using the converter 'by'
func Convert[FS ~[]From, From, To any](elements FS, by func(From) To) []To {
	if elements == nil {
		return nil
	}
	result := make([]To, len(elements))
	for i, e := range elements {
		result[i] = by(e)
	}
	return result
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, by func(From) To) []To {
	var result []To
	for _, e := range elements {
		if filter(e) {
			result = append(result, by(e))
		}
	}
	return result
}

// ConvertAndFilter additionally filters 'To' elements
func ConvertAndFilter[FS ~[]From, From, To any](elements FS, by func(From) To, filter func(To) bool) []To {
	var result []To
	for _, e := range elements {
		if r := by(e); filter(r) {
			result = append(result, r)
		}
	}
	return result
}

// FilterConvertFilter filters source, converts, and filters converted elements
func FilterConvertFilter[FS ~[]From, From, To any](elements FS, filter func(From) bool, by func(From) To, filterConverted func(To) bool) []To {
	var result []To
	for _, e := range elements {
		if filter(e) {
			if r := by(e); filterConverted(r) {
				result = append(result, r)
			}
		}
	}
	return result
}

// ConvertIndexed creates a slice consisting of the transformed elements using the converter 'by' which additionally applies the index of the element being converted
func ConvertIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) To) []To {
	if elements == nil {
		return nil
	}
	result := make([]To, len(elements))
	for i, e := range elements {
		result[i] = by(i, e)
	}
	return result
}

// FilterAndConvertIndexed additionally filters 'From' elements
func FilterAndConvertIndexed[FS ~[]From, From, To any](elements FS, filter func(index int, from From) bool, converter func(index int, from From) To) []To {
	var result []To
	for i, e := range elements {
		if filter(i, e) {
			result = append(result, converter(i, e))
		}
	}
	return result
}

// ConvertCheck is similar to ConvertFit, but it checks and transforms elements together
func ConvertCheck[FS ~[]From, From, To any](elements FS, by func(from From) (To, bool)) []To {
	var result []To
	for _, e := range elements {
		if to, ok := by(e); ok {
			result = append(result, to)
		}
	}
	return result
}

// ConvertCheckIndexed additionally filters 'From' elements
func ConvertCheckIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) (To, bool)) []To {
	var result []To
	for i, e := range elements {
		if to, ok := by(i, e); ok {
			result = append(result, to)
		}
	}
	return result
}

// Flatt unfolds the n-dimensional slice into a n-1 dimensional slice
func Flatt[FS ~[]From, From, To any](elements FS, by func(From) []To) []To {
	var result []To
	for _, e := range elements {
		result = append(result, by(e)...)

	}
	return result
}

// FilerAndFlatt additionally filters 'From' elements.
func FilerAndFlatt[FS ~[]From, From, To any](elements FS, filter func(From) bool, by func(From) []To) []To {
	var result []To
	for _, e := range elements {
		if filter(e) {
			result = append(result, by(e)...)
		}
	}
	return result
}

// FlattAndFiler unfolds the n-dimensional slice into a n-1 dimensional slice with additinal filtering of 'To' elements.
func FlattAndFiler[FS ~[]From, From, To any](elements FS, by func(From) []To, filter func(To) bool) []To {
	var result []To
	for _, e := range elements {
		for _, to := range by(e) {
			if filter(to) {
				result = append(result, to)
			}
		}
	}
	return result
}

// FilterFlattFilter unfolds the n-dimensional slice 'elements' into a n-1 dimensional slice with additinal filtering of 'From' and 'To' elements.
func FilterFlattFilter[FS ~[]From, From, To any](elements FS, filterFrom func(From) bool, by func(From) []To, filterTo func(To) bool) []To {
	var result []To
	for _, e := range elements {
		if filterFrom(e) {
			for _, to := range by(e) {
				if filterTo(to) {
					result = append(result, to)
				}
			}
		}
	}
	return result
}

func NotNil[TS ~[]*T, T any](elements TS) TS {
	return Filter(elements, check.NotNil[T])
}

// Filter creates a slice containing only the filtered elements
func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) []T {
	var result []T
	for _, e := range elements {
		if filter(e) {
			result = append(result, e)
		}
	}
	return result
}

// Range generates a slice of integers in the range defined by from and to inclusive.
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
func Reverse[TS ~[]T, T any](elements TS) []T {
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

type Less[T any] func(e1, e2 T) bool
type Sorter func(x any, less func(i, j int) bool)

// Sort sorts elements in place using a function that checks if an element is smaller than the others
func Sort[TS ~[]T, T any](elements TS, sorter Sorter, less Less[T]) TS {
	sorter(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	return elements
}

// SortByOrdered sorts elements in place by converting them to constraints.Ordered values and applying the operator <
func SortByOrdered[T any, o constraints.Ordered, TS ~[]T](elements TS, sorter Sorter, by func(T) o) TS {
	return Sort(elements, sorter, func(e1, e2 T) bool { return by(e1) < by(e2) })
}

// Reduce reduces elements to an one
func Reduce[TS ~[]T, T any](elements TS, by func(T, T) T) T {
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

// First returns the first element that satisfies requirements of the predicate 'filter'
func First[TS ~[]T, T any](elements TS, by func(T) bool) (T, bool) {
	for _, e := range elements {
		if by(e) {
			return e, true
		}
	}
	var no T
	return no, false
}

// Last returns the latest element that satisfies requirements of the predicate 'filter'
func Last[TS ~[]T, T any](elements TS, by func(T) bool) (T, bool) {
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
func Get[TS ~[]T, T any](elements TS, index int) (T, bool) {
	l := len(elements)
	if l > 0 && (index >= 0 || index < l) {
		return elements[index], true
	}
	var no T
	return no, false
}

// Track applies tracker to elements with error checking. To stop traking just return the ErrBreak
func Track[TS ~[]T, T any](elements TS, tracker func(int, T) error) error {
	for i, e := range elements {
		if err := tracker(i, e); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// TrackEach applies tracker to elements without error checking
func TrackEach[TS ~[]T, T any](elements TS, tracker func(int, T)) {
	for i, e := range elements {
		tracker(i, e)
	}
}

// For applies walker to elements. To stop walking just return the ErrBreak
func For[TS ~[]T, T any](elements TS, walker func(T) error) error {
	for _, e := range elements {
		if err := walker(e); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEach applies walker to elements without error checking
func ForEach[TS ~[]T, T any](elements TS, walker func(T)) {
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
func ToString[TS ~[]T, T any](elements TS) string {
	return ToStringf(elements, "%+v", " ")
}

// ToStringf converts elements to a string representation defined by a custom element format and a delimiter
func ToStringf[TS ~[]T, T any](elements TS, elementFormat, delimeter string) string {
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

// BehaveAsStrings draws a string inherited type slice as the slice of strings
func BehaveAsStrings[T ~string, TS ~[]T](elements TS) []string {
	ptr := unsafe.Pointer(&elements)
	s := *(*[]string)(ptr)
	return s
}

// StringsBehaveAs draws a string slice as the slice of a string inherited type
func StringsBehaveAs[TS ~[]T, T ~string](elements []string) TS {
	ptr := unsafe.Pointer(&elements)
	s := *(*TS)(ptr)
	return s
}

func NoEmpty[TS ~[]T, T any](elements TS, def []T) TS {
	if len(elements) > 0 {
		return elements
	}
	return def
}

func GetNoEmpty[TS ~[]T, T any](elementsProducer func() TS, def []T) TS {
	return NoEmpty(elementsProducer(), def)
}
