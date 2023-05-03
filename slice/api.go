// Package slice provides generic functions for slice types
package slice

import (
	"bytes"
	"fmt"
	"unsafe"

	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/op"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = loop.ErrBreak

// Of is generic slice constructor
func Of[T any](elements ...T) []T { return elements }

// Len return the length of the 'elements' slice
func Len[TS ~[]T, T any](elements TS) int {
	return len(elements)
}

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
func Generate[T any](next func() (T, bool)) []T {
	return loop.Slice(next)
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
	if elements == nil {
		return nil
	}
	return append(elements[0:index], elements[index+1:]...)
}

// Group converts the 'elements' slice into a map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to an value.
func Group[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) K, valExtractor func(T) V) map[K][]V {
	if elements == nil {
		return nil
	}
	return ToMapResolv(elements, keyExtractor, valExtractor, resolv.Append[K, V])
}

// GroupInMultiple converts the 'elements' slice into a map, extracting multiple keys, values per each element applying the 'keysExtractor' and 'valsExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valsExtractor retrieves one or more values per element.
func GroupInMultiple[TS ~[]T, T any, K comparable, V any](elements TS, keysExtractor func(T) []K, valsExtractor func(T) []V) map[K][]V {
	if elements == nil {
		return nil
	}
	groups := map[K][]V{}
	for _, e := range elements {
		if keys, vals := keysExtractor(e), valsExtractor(e); len(keys) == 0 {
			var key K
			for _, v := range vals {
				initGroup(key, v, groups)
			}
		} else {
			for _, key := range keys {
				if len(vals) == 0 {
					var v V
					initGroup(key, v, groups)
				} else {
					for _, v := range vals {
						initGroup(key, v, groups)
					}
				}
			}
		}
	}
	return groups
}

// GroupInMultipleKeys converts the 'elements' slice into a map, extracting multiple keys, one value per each element applying the 'keysExtractor' and 'valExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valExtractor converts an element to a value.
func GroupInMultipleKeys[TS ~[]T, T any, K comparable, V any](elements TS, keysExtractor func(T) []K, valExtractor func(T) V) map[K][]V {
	if elements == nil {
		return nil
	}
	groups := map[K][]V{}
	for _, e := range elements {
		if keys, v := keysExtractor(e), valExtractor(e); len(keys) == 0 {
			var key K
			initGroup(key, v, groups)
		} else {
			for _, key := range keys {
				initGroup(key, v, groups)
			}
		}
	}
	return groups
}

// GroupInMultipleValues converts the 'elements' slice into a map, extracting one key, multiple values per each element applying the 'keyExtractor' and 'valsExtractor' functions.
// The keyExtractor converts an element to a key.
// The valsExtractor retrieves one or more values per element.
func GroupInMultipleValues[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) K, valsExtractor func(T) []V) map[K][]V {
	if elements == nil {
		return nil
	}
	groups := map[K][]V{}
	for _, e := range elements {
		if key, vals := keyExtractor(e), valsExtractor(e); len(vals) == 0 {
			var v V
			initGroup(key, v, groups)
		} else {
			for _, v := range vals {
				initGroup(key, v, groups)
			}
		}
	}
	return groups
}

func initGroup[K comparable, T any, TS ~[]T](key K, e T, groups map[K]TS) {
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

// Conv creates a slice consisting of the transformed elements using the converter 'by'
func Conv[FS ~[]From, From, To any](elements FS, by func(From) (To, error)) ([]To, error) {
	if elements == nil {
		return nil, nil
	}
	result := make([]To, len(elements))
	for i, e := range elements {
		c, err := by(e)
		if err != nil {
			return result[:i], err
		}
		result[i] = c
	}
	return result, nil
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, by func(From) To) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, len(elements)/2)
	for _, e := range elements {
		if filter(e) {
			result = append(result, by(e))
		}
	}
	return result
}

// ConvertAndFilter additionally filters 'To' elements
func ConvertAndFilter[FS ~[]From, From, To any](elements FS, by func(From) To, filter func(To) bool) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, len(elements)/2)
	for _, e := range elements {
		if r := by(e); filter(r) {
			result = append(result, r)
		}
	}
	return result
}

// FilterConvertFilter filters source, converts, and filters converted elements
func FilterConvertFilter[FS ~[]From, From, To any](elements FS, filter func(From) bool, by func(From) To, filterConverted func(To) bool) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, len(elements)/4)
	for _, e := range elements {
		if filter(e) {
			if r := by(e); filterConverted(r) {
				result = append(result, r)
			}
		}
	}
	return result
}

// ConvertIndexed creates a slice consisting of the transformed elements using the 'converter' function which additionally applies the index of the element being converted
func ConvertIndexed[FS ~[]From, From, To any](elements FS, converter func(index int, from From) To) []To {
	if elements == nil {
		return nil
	}
	result := make([]To, len(elements))
	for i, e := range elements {
		result[i] = converter(i, e)
	}
	return result
}

// FilterAndConvertIndexed additionally filters 'From' elements
func FilterAndConvertIndexed[FS ~[]From, From, To any](elements FS, filter func(index int, from From) bool, converter func(index int, from From) To) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, len(elements)/2)
	for i, e := range elements {
		if filter(i, e) {
			result = append(result, converter(i, e))
		}
	}
	return result
}

// ConvertCheck is similar to ConvertFit, but it checks and transforms elements together
func ConvertCheck[FS ~[]From, From, To any](elements FS, by func(from From) (To, bool)) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, len(elements))
	for _, e := range elements {
		if to, ok := by(e); ok {
			result = append(result, to)
		}
	}
	return result
}

// ConvertCheckIndexed additionally filters 'From' elements
func ConvertCheckIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) (To, bool)) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, len(elements))
	for i, e := range elements {
		if to, ok := by(i, e); ok {
			result = append(result, to)
		}
	}
	return result
}

// Flatt unfolds the n-dimensional slice into a n-1 dimensional slice
func Flatt[FS ~[]From, From, To any](elements FS, flattener func(From) []To) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, int(float32(len(elements))*1.618))
	for _, e := range elements {
		result = append(result, flattener(e)...)

	}
	return result
}

// Flat unfolds the n-dimensional slice into a n-1 dimensional slice
func Flat[FS ~[]From, From, To any](elements FS, flattener func(From) ([]To, error)) ([]To, error) {
	if elements == nil {
		return nil, nil
	}
	var result = make([]To, 0, int(float32(len(elements))*1.618))
	for _, e := range elements {
		f, err := flattener(e)
		if err != nil {
			return result, err
		}
		result = append(result, f...)

	}
	return result, nil
}

// FlattAndConvert unfolds the n-dimensional slice into a n-1 dimensional slice and converts the elements
func FlattAndConvert[FS ~[]From, From, I, To any](elements FS, flattener func(From) []I, convert func(I) To) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, int(float32(len(elements))*1.618))
	for _, e := range elements {
		for _, f := range flattener(e) {
			result = append(result, convert(f))
		}
	}
	return result
}

// FilterAndFlatt filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlatt[FS ~[]From, From, To any](elements FS, filter func(From) bool, flattener func(From) []To) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, int(float32(len(elements)/2)*1.618))
	for _, e := range elements {
		if filter(e) {
			result = append(result, flattener(e)...)
		}
	}
	return result
}

// FlattAndFiler unfolds the n-dimensional slice into a n-1 dimensional slice with additinal filtering of 'To' elements.
func FlattAndFiler[FS ~[]From, From, To any](elements FS, by func(From) []To, filter func(To) bool) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, int(float32(len(elements))*1.618)/2)
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
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, int(float32(len(elements)/2)*1.618)/2)
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

// NotNil returns only not nil elements
func NotNil[TS ~[]*T, T any](elements TS) TS {
	return Filter(elements, check.NotNil[T])
}

// Filter creates a slice containing only the filtered elements
func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) []T {
	if elements == nil {
		return nil
	}
	var result = make([]T, 0, len(elements)/2)
	for _, e := range elements {
		if filter(e) {
			result = append(result, e)
		}
	}
	return result
}

// Filt creates a slice containing only the filtered elements
func Filt[TS ~[]T, T any](elements TS, filter func(T) (bool, error)) ([]T, error) {
	if elements == nil {
		return nil, nil
	}
	var result = make([]T, 0, len(elements)/2)
	for _, e := range elements {
		if ok, err := filter(e); err != nil {
			return result, err
		} else if ok {
			result = append(result, e)
		}
	}
	return result, nil
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
	if elements == nil {
		return nil
	}
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

// Less less element qualifier alias.
// Is a function that must return true it the first element is less the second
type Less[T any] func(first, second T) bool

// Sorter is alias for sort.Slice or SliceStable functions
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

// Reduce reduces the elements into an one using the 'merge' function
func Reduce[TS ~[]T, T any](elements TS, merge func(T, T) T) T {
	var result T
	for i, v := range elements {
		if i == 0 {
			result = v
		} else {
			result = merge(result, v)
		}
	}
	return result
}

// Sum returns the sum of all elements
func Sum[TS ~[]T, T c.Summable](elements TS) T {
	return Reduce(elements, op.Sum[T])
}

// First returns the first element that satisfies requirements of the predicate 'by'
func First[TS ~[]T, T any](elements TS, by func(T) bool) (no T, ok bool) {
	for _, e := range elements {
		if by(e) {
			return e, true
		}
	}
	return no, false
}

// Firstt returns the first element that satisfies the condition of the 'by' function
func Firstt[TS ~[]T, T any](elements TS, by func(T) (bool, error)) (no T, ok bool, err error) {
	for _, e := range elements {
		if ok, err = by(e); err != nil || ok {
			return e, ok, err
		}
	}
	return no, false, nil
}

// Last returns the latest element that satisfies requirements of the predicate 'filter'
func Last[TS ~[]T, T any](elements TS, by func(T) bool) (no T, ok bool) {
	for i := len(elements) - 1; i >= 0; i-- {
		e := elements[i]
		if by(e) {
			return e, true
		}
	}
	return no, false
}

// Lastt returns the latest element that satisfies requirements of the predicate 'filter'
func Lastt[TS ~[]T, T any](elements TS, by func(T) (bool, error)) (no T, ok bool, err error) {
	for i := len(elements) - 1; i >= 0; i-- {
		e := elements[i]
		if ok, err = by(e); err != nil {
			return no, false, err
		} else if ok {
			return e, true, nil
		}
	}
	return no, false, nil
}

// Track applies the 'tracker' function to the elements. Return the c.ErrBreak to stop tracking
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

// TrackEach applies the 'tracker' function to the elements
func TrackEach[TS ~[]T, T any](elements TS, tracker func(int, T)) {
	for i, e := range elements {
		tracker(i, e)
	}
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop
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

// ForEach applies the 'walker' function for the elements
func ForEach[TS ~[]T, T any](elements TS, walker func(T)) {
	for _, e := range elements {
		walker(e)
	}
}

// ForEachRef applies the 'walker' function for the references
func ForEachRef[T any, TS ~[]*T](references TS, walker func(T)) {
	for _, e := range references {
		walker(*e)
	}
}

// ToString converts the elements to their default string representation
func ToString[TS ~[]T, T any](elements TS) string {
	return ToStringf(elements, "%+v", " ")
}

// ToStringf converts the elements to a string representation defined by the elementFormat and a delimiter
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

// ToStringRefsf converts references to a string representation defined by the delimiter and the nilValue representation
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

// Upcast transforms a user-defined type slice to a underlying slice type
func Upcast[TS ~[]T, T any](elements TS) []T {
	return *UpcastRef(&elements)
}

// UpcastRef transforms a user-defined reference type slice to a underlying slice reference type
func UpcastRef[TS ~[]T, T any](elements *TS) *[]T {
	ptr := unsafe.Pointer(elements)
	s := (*[]T)(ptr)
	return s
}

// Downcast transforms a slice type to an user-defined type slice based on that type
func Downcast[TS ~[]T, T any](elements []T) TS {
	return *DowncastRef[TS](&elements)
}

// DowncastRef  transforms a slice typeref to an user-defined type ref slice based on that type
func DowncastRef[TS ~[]T, T any](elements *[]T) *TS {
	ptr := unsafe.Pointer(elements)
	s := (*TS)(ptr)
	return s
}

// Filled returns the 'ifEmpty' if the 'elements' slise is empty
func Filled[TS ~[]T, T any](elements TS, ifEmpty []T) TS {
	if Empty(elements) {
		return elements
	}
	return ifEmpty
}

// GetFilled returns the 'notEmpty' if the 'elementsFactory' return an empty slice
func GetFilled[TS ~[]T, T any](elementsFactory func() TS, ifEmpty []T) TS {
	return Filled(elementsFactory(), ifEmpty)
}

// HasAny tests if the 'elements' slice contains an element that satisfies the "predicate" condition
func HasAny[TS ~[]T, T any](elements TS, predicate func(T) bool) bool {
	_, ok := First(elements, predicate)
	return ok
}

// Contains checks is the 'elements' slice contains the example
func Contains[TS ~[]T, T comparable](elements TS, example T) bool {
	for _, e := range elements {
		if e == example {
			return true
		}
	}
	return false
}

// Has checks is the 'elements' slice contains a value that satisfies the specified condition
func Has[TS ~[]T, T any](elements TS, condition func(T) bool) bool {
	for _, e := range elements {
		if condition(e) {
			return true
		}
	}
	return false
}

// Empty checks whether the specified slice is empty
func Empty[TS ~[]T, T any](elements TS) bool {
	return len(elements) == 0
}

// NotEmpty checks whether the specified slice is not empty
func NotEmpty[TS ~[]T, T any](elements TS) bool {
	return !Empty(elements)
}

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[TS ~[]T, T any, K comparable, V, VR any](elements TS, keyExtractor func(T) K, valExtractor func(T) V, resolver func(bool, K, VR, V) VR) map[K]VR {
	m := make(map[K]VR, len(elements))
	for _, e := range elements {
		k, v := keyExtractor(e), valExtractor(e)
		exists, ok := m[k]
		m[k] = resolver(ok, k, exists, v)
	}
	return m
}

// ToKV transforms slice elements to key/value pairs slice. One pair per one element
func ToKV[TS ~[]T, T, K, V any](elements TS, keyExtractor func(T) K, valExtractor func(T) V) []c.KV[K, V] {
	return Convert(elements, func(e T) c.KV[K, V] { return convert.ToKV(e, keyExtractor, valExtractor) })
}

// ToKVs transforms slice elements to key/value pairs slice. Multiple pairs per one element
func ToKVs[TS ~[]T, T, K, V any](elements TS, keysExtractor func(T) []K, valsExtractor func(T) []V) []c.KV[K, V] {
	return Flatt(elements, func(e T) []c.KV[K, V] { return convert.ToKVs(e, keysExtractor, valsExtractor) })
}

// FlattValues transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlattValues[TS ~[]T, T, V any](elements TS, valsExtractor func(T) []V) []c.KV[T, V] {
	return Flatt(elements, func(e T) []c.KV[T, V] { return convert.FlattValues(e, valsExtractor) })
}

// FlattKeys transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlattKeys[TS ~[]T, T, K any](elements TS, keysExtractor func(T) []K) []c.KV[K, T] {
	return Flatt(elements, func(e T) []c.KV[K, T] { return convert.FlattKeys(e, keysExtractor) })
}
