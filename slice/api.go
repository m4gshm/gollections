// Package slice provides generic functions for slice types
package slice

import (
	"bytes"
	"fmt"
	"slices"
	"unsafe"

	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/comparer"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
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

// OfIndexed builds a slice by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](len int, getAt func(int) T) []T {
	r := make([]T, len)
	for i := 0; i < len; i++ {
		r[i] = getAt(i)
	}
	return r
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
func Delete[TS ~[]T, T any](elements TS, index int) TS {
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
	return ToMapResolv(elements, keyExtractor, valExtractor, resolv.Slice[K, V])
}

// GroupByMultiple converts the 'elements' slice into a map, extracting multiple keys, values per each element applying the 'keysExtractor' and 'valsExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valsExtractor retrieves one or more values per element.
func GroupByMultiple[TS ~[]T, T any, K comparable, V any](elements TS, keysExtractor func(T) []K, valsExtractor func(T) []V) map[K][]V {
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

// GroupByMultipleKeys converts the 'elements' slice into a map, extracting multiple keys, one value per each element applying the 'keysExtractor' and 'valExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valExtractor converts an element to a value.
func GroupByMultipleKeys[TS ~[]T, T any, K comparable, V any](elements TS, keysExtractor func(T) []K, valExtractor func(T) V) map[K][]V {
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

// GroupByMultipleValues converts the 'elements' slice into a map, extracting one key, multiple values per each element applying the 'keyExtractor' and 'valsExtractor' functions.
// The keyExtractor converts an element to a key.
// The valsExtractor retrieves one or more values per element.
func GroupByMultipleValues[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) K, valsExtractor func(T) []V) map[K][]V {
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

// Convert creates a slice consisting of the transformed elements using the converter.
func Convert[FS ~[]From, From, To any](elements FS, converter func(From) To) []To {
	if elements == nil {
		return nil
	}
	result := make([]To, len(elements))
	for i, e := range elements {
		result[i] = converter(e)
	}
	return result
}

// Conv creates a slice consisting of the transformed elements using the converter.
func Conv[FS ~[]From, From, To any](elements FS, converter func(From) (To, error)) ([]To, error) {
	if elements == nil {
		return nil, nil
	}
	result := make([]To, len(elements))
	for i, e := range elements {
		c, err := converter(e)
		if err != nil {
			return result[:i], err
		}
		result[i] = c
	}
	return result, nil
}

// FilterAndConvert selects elements that match the filter, converts and places them into a new slice.
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To) []To {
	if elements == nil {
		return nil
	}
	return AppendFilterAndConvert(elements, make([]To, 0, len(elements)/2), filter, converter)
}

// AppendFilterAndConvert selects elements that match the filter, converts and appends them to the dest.
func AppendFilterAndConvert[FS ~[]From, DS ~[]To, From, To any](src FS, dest DS, filter func(From) bool, converter func(From) To) DS {
	for _, e := range src {
		if filter(e) {
			dest = append(dest, converter(e))
		}
	}
	return dest
}

// FiltAndConv selects elements that match the filter, converts and returns them.
func FiltAndConv[FS ~[]From, From, To any](elements FS, filter func(From) (bool, error), by func(From) (To, error)) ([]To, error) {
	if elements == nil {
		return nil, nil
	}
	return AppendFiltAndConv(elements, make([]To, 0, len(elements)/2), filter, by)
}

// AppendFiltAndConv selects elements that match the filter, converts and appends them to the dest.
func AppendFiltAndConv[FS ~[]From, DS ~[]To, From, To any](src FS, dest DS, filter func(From) (bool, error), by func(From) (To, error)) (DS, error) {
	for _, e := range src {
		if ok, err := filter(e); err != nil {
			return dest, err
		} else if ok {
			c, err := by(e)
			if err != nil {
				return dest, err
			}
			dest = append(dest, c)
		}
	}
	return dest, nil
}

// ConvertAndFilter converts elements, filters and returns them.
func ConvertAndFilter[FS ~[]From, From, To any](elements FS, converter func(From) To, filter func(To) bool) []To {
	if elements == nil {
		return nil
	}
	return AppendConvertAndFilter(elements, make([]To, 0, len(elements)/2), converter, filter)
}

// AppendConvertAndFilter converts elements, filters and append mached ones to the dest slice.
func AppendConvertAndFilter[FS ~[]From, DS ~[]To, From, To any](src FS, dest DS, converter func(From) To, filter func(To) bool) DS {
	for _, e := range src {
		if r := converter(e); filter(r) {
			dest = append(dest, r)
		}
	}
	return dest
}

// FilterConvertFilter applies operations chain: filter, convert, filter the converted elemens of the src slice and returns them.
func FilterConvertFilter[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To, filterConverted func(To) bool) []To {
	if elements == nil {
		return nil
	}
	return AppendFilterConvertFilter(elements, make([]To, 0, len(elements)/2), filter, converter, filterConverted)
}

// AppendFilterConvertFilter applies operations chain filter->convert->filter the src elemens and addends to the dest slice.
func AppendFilterConvertFilter[FS ~[]From, DS ~[]To, From, To any](src FS, dest DS, filter func(From) bool, converter func(From) To, filterConverted func(To) bool) DS {
	for _, e := range src {
		if filter(e) {
			if r := converter(e); filterConverted(r) {
				dest = append(dest, r)
			}
		}
	}
	return dest
}

// ConvertIndexed converets the elements using the converter that takes the index and value of each element from the elements slice.
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

// FilterAndConvertIndexed filter elements that match the filter condition, converts and returns a slice of result elements.
func FilterAndConvertIndexed[FS ~[]From, From, To any](elements FS, filter func(index int, from From) bool, converter func(index int, from From) To) []To {
	if elements == nil {
		return nil
	}
	return AppendFilterAndConvertIndexed(elements, make([]To, 0, len(elements)/2), filter, converter)
}

// AppendFilterAndConvertIndexed filters elements that match the filter condition, converts them and appends to the dest.
func AppendFilterAndConvertIndexed[FS ~[]From, DS ~[]To, From, To any](src FS, dest DS, filter func(index int, from From) bool, converter func(index int, from From) To) DS {
	for i, e := range src {
		if filter(i, e) {
			dest = append(dest, converter(i, e))
		}
	}
	return dest
}

// ConvertCheck is similar to ConvertFilt, but it checks and transforms elements together
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

// Flat unfolds the n-dimensional slice into a n-1 dimensional slice
func Flat[FS ~[]From, From any, TS ~[]To, To any](elements FS, flattener func(From) TS) []To {
	if elements == nil {
		return nil
	}
	var result = make([]To, 0, int(float32(len(elements))*1.618))
	for _, e := range elements {
		result = append(result, flattener(e)...)

	}
	return result
}

// Flatt unfolds the n-dimensional slice into a n-1 dimensional slice
func Flatt[FS ~[]From, From, To any](elements FS, flattener func(From) ([]To, error)) ([]To, error) {
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

// FlatAndConvert unfolds the n-dimensional slice into a n-1 dimensional slice and converts the elements
func FlatAndConvert[FS ~[]From, From any, IS ~[]I, I, To any](elements FS, flattener func(From) IS, convert func(I) To) []To {
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

// FilterAndFlat retrieves src elements that match the filter condition, extracts 'To' type slices from them and joins into a new slice.
func FilterAndFlat[FS ~[]From, From, To any](elements FS, filter func(From) bool, flattener func(From) []To) []To {
	if elements == nil {
		return nil
	}
	return AppendFilterAndFlat(elements, make([]To, 0, int(float32(len(elements)/2)*1.618)), filter, flattener)
}

// AppendFilterAndFlat retrieves src elements that match the filter condition, extracts 'To' type slices from them and appends the dest.
func AppendFilterAndFlat[FS ~[]From, DS ~[]To, From, To any](src FS, dest DS, filter func(From) bool, flattener func(From) []To) DS {
	for _, e := range src {
		if filter(e) {
			dest = append(dest, flattener(e)...)
		}
	}
	return dest
}

// FlatAndFiler extracts a slice of type "To" from each src element, filters and joins into a new slice.
func FlatAndFiler[FS ~[]From, From, To any](elements FS, flattener func(From) []To, filter func(To) bool) []To {
	if elements == nil {
		return nil
	}
	return AppendFlatAndFiler(elements, make([]To, 0, int(float32(len(elements))*1.618)/2), flattener, filter)
}

// AppendFlatAndFiler extracts a slice of type "To" from each src element, filters and appends to dest.
func AppendFlatAndFiler[FS ~[]From, DS ~[]To, From, To any](src FS, dest DS, flattener func(From) []To, filter func(To) bool) DS {
	for _, e := range src {
		for _, to := range flattener(e) {
			if filter(to) {
				dest = append(dest, to)
			}
		}
	}
	return dest
}

// FilterFlatFilter applies operations chain filter->flat->filter the src elemens and returns that result.
func FilterFlatFilter[FS ~[]From, From, To any](elements FS, filterFrom func(From) bool, flat func(From) []To, filterTo func(To) bool) []To {
	if elements == nil {
		return nil
	}
	return AppendFilterFlatFilter(elements, make([]To, 0, int(float32(len(elements)/2)*1.618)/2), filterFrom, flat, filterTo)
}

// AppendFilterFlatFilter applies operations chain filter->flat->filter the src elemens and addends to the dest.
func AppendFilterFlatFilter[FS ~[]From, DS ~[]To, From, To any](src FS, dest DS, filterFrom func(From) bool, flat func(From) []To, filterTo func(To) bool) DS {
	for _, e := range src {
		if filterFrom(e) {
			for _, to := range flat(e) {
				if filterTo(to) {
					dest = append(dest, to)
				}
			}
		}
	}
	return dest
}

// NotNil returns only not nil elements
func NotNil[TS ~[]*T, T any](elements TS) TS {
	return Filter(elements, not.Nil[T])
}

// ToValues returns values referenced by the pointers.
// If a pointer is nil then it is replaced by the zero value.
func ToValues[TS ~[]*T, T any](pointers TS) []T {
	return Convert(pointers, convert.PtrVal[T])
}

// GetValues returns values referenced by the pointers.
// All nil pointers are excluded from the final result.
func GetValues[TS ~[]*T, T any](elements TS) []T {
	return ConvertCheck(elements, convert.NoNilPtrVal[T])
}

// Filter filters elements that match the filter condition and returns them.
func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) []T {
	if elements == nil {
		return nil
	}
	return AppendFilter(elements, make([]T, 0, len(elements)/2), filter)
}

// AppendFilter filters elements that match the filter condition and adds them to the dest.
func AppendFilter[TS ~[]T, DS ~[]T, T any](src TS, dest DS, filter func(T) bool) DS {
	for _, e := range src {
		if filter(e) {
			dest = append(dest, e)
		}
	}
	return dest
}

// Filt filters elements that match the filter condition and returns them.
func Filt[TS ~[]T, T any](elements TS, filter func(T) (bool, error)) ([]T, error) {
	if elements == nil {
		return nil, nil
	}
	return AppendFilt(elements, make([]T, 0, len(elements)/2), filter)
}

// AppendFilt filters elements that match the filter condition and adds them to the dest.
func AppendFilt[TS ~[]T, DS ~[]T, T any](src TS, dest DS, filter func(T) (bool, error)) (DS, error) {
	for _, e := range src {
		ok, err := filter(e)
		if ok {
			dest = append(dest, e)
		}
		if err != nil {
			return dest, err
		}
	}
	return dest, nil
}

// RangeClosed generates a slice of integers in the range defined by from and to inclusive
func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) []T {
	if toInclusive == from {
		return []T{from}
	}
	amount := toInclusive - from
	delta := 1
	if amount < 0 {
		amount = -amount
		delta = -1
	}
	amount++

	elements := make([]T, amount)
	e := from
	for i := 0; i < int(amount); i++ {
		elements[i] = e
		e = e + T(delta)
	}
	return elements
}

// Range generates a slice of integers in the range defined by from and to exclusive
func Range[T constraints.Integer | rune](from T, toExclusive T) []T {
	if toExclusive == from {
		return nil
	}
	amount := toExclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	elements := make([]T, amount)
	e := from
	for i := T(0); i < T(amount); i++ {
		elements[i] = e
		e = e + delta
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

// Comparer aims to compare two values and must return a positive num if the first value is more then the second, a negative if less, and 0 if they equal.
type Comparer[T any] func(T, T) int

// Sorter is alias for slices.SortFunc or slices.SortStableFunc functions
type Sorter[TS ~[]T, T any] func(TS, func(T, T) int)

// Sort sorts elements in place using the comparer function
func Sort[TS ~[]T, T any](elements TS, comparer Comparer[T]) TS {
	return sort(elements, slices.SortFunc, comparer)
}

// StableSort sorts elements in place using the comparer function
func StableSort[TS ~[]T, T any](elements TS, comparer Comparer[T]) TS {
	return sort(elements, slices.SortStableFunc, comparer)
}

// SortAsc sorts elements in ascending order, using the orderConverner function to retrieve a value of type Ordered.
func SortAsc[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return Sort(elements, comparer.Of(orderConverner))
}

// StableSortAsc sorts elements in ascending order, using the orderConverner function to retrieve a value of type Ordered.
func StableSortAsc[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return StableSort(elements, comparer.Of(orderConverner))
}

// SortDesc sorts elements in descending order, using the orderConverner function to retrieve a value of type Ordered.
func SortDesc[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return Sort(elements, comparer.Reverse(orderConverner))
}

// StableSortDesc sorts elements in descending order, using the orderConverner function to retrieve a value of type Ordered.
func StableSortDesc[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return StableSort(elements, comparer.Reverse(orderConverner))
}

func sort[TS ~[]T, T any](elements TS, sorter Sorter[TS, T], comparer Comparer[T]) TS {
	sorter(elements, comparer)
	return elements
}

// Reduce reduces the elements into an one using the 'merge' function
func Reduce[TS ~[]T, T any](elements TS, merge func(T, T) T) (out T) {
	l := len(elements)
	if l >= 1 {
		out = elements[0]
	}
	if l > 1 {
		for _, v := range elements[1:] {
			out = merge(out, v)
		}
	}
	return out
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

func All[TS ~[]T, T any](elements TS, yield func(T) bool) {
	for i := 0; i < len(elements) && yield(elements[i]); i++ {

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

// ToMap collects key\value elements to a map by iterating over the elements
func ToMap[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) K, valExtractor func(T) V) map[K]V {
	return ToMapResolv(elements, keyExtractor, valExtractor, resolv.First[K, V])
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

// KeyValue transforms slice elements to key/value pairs slice. One pair per one element
func KeyValue[TS ~[]T, T, K, V any](elements TS, keyExtractor func(T) K, valExtractor func(T) V) []c.KV[K, V] {
	return Convert(elements, func(e T) c.KV[K, V] { return convert.KeyValue(e, keyExtractor, valExtractor) })
}

// KeysValues transforms slice elements to key/value pairs slice. Multiple pairs per one element
func KeysValues[TS ~[]T, T, K, V any](elements TS, keysExtractor func(T) []K, valsExtractor func(T) []V) []c.KV[K, V] {
	return Flat(elements, func(e T) []c.KV[K, V] { return convert.KeysValues(e, keysExtractor, valsExtractor) })
}

// ExtraVals transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func ExtraVals[TS ~[]T, T, V any](elements TS, valsExtractor func(T) []V) []c.KV[T, V] {
	return Flat(elements, func(e T) []c.KV[T, V] { return convert.ExtraVals(e, valsExtractor) })
}

// ExtraKeys transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func ExtraKeys[TS ~[]T, T, K any](elements TS, keysExtractor func(T) []K) []c.KV[K, T] {
	return Flat(elements, func(e T) []c.KV[K, T] { return convert.ExtraKeys(e, keysExtractor) })
}

// SplitTwo splits the elements into two slices
func SplitTwo[TS ~[]T, T, F, S any](elements TS, splitter func(T) (F, S)) ([]F, []S) {
	var (
		l      = len(elements)
		first  = make([]F, l)
		second = make([]S, l)
	)
	for i, e := range elements {
		first[i], second[i] = splitter(e)
	}
	return first, second
}

// SplitThree splits the elements into three slices
func SplitThree[TS ~[]T, T, F, S, TH any](elements TS, splitter func(T) (F, S, TH)) ([]F, []S, []TH) {
	var (
		l      = len(elements)
		first  = make([]F, l)
		second = make([]S, l)
		third  = make([]TH, l)
	)
	for i, e := range elements {
		first[i], second[i], third[i] = splitter(e)
	}
	return first, second, third
}

// SplitAndReduceTwo splits each element of the specified slice into two values and then reduces that ones
func SplitAndReduceTwo[TS ~[]T, T, F, S any](elements TS, splitter func(T) (F, S), firstMerge func(F, F) F, secondMerger func(S, S) S) (first F, second S) {
	l := len(elements)
	if l >= 1 {
		first, second = splitter(elements[0])
	}
	if l > 1 {
		for _, e := range elements[1:] {
			f, s := splitter(e)
			first = firstMerge(first, f)
			second = secondMerger(second, s)
		}
	}
	return first, second
}

// ConvertAndReduce converts each elements and merges them into one
func ConvertAndReduce[FS ~[]From, From, To any](elements FS, converter func(From) To, merger func(To, To) To) (out To) {
	l := len(elements)
	if l >= 1 {
		out = converter(elements[0])
	}
	if l > 1 {
		for _, e := range elements[1:] {
			out = merger(out, converter(e))
		}
	}
	return out
}

// ConvAndReduce converts each elements and merges them into one
func ConvAndReduce[FS ~[]From, From, To any](elements FS, converter func(From) (To, error), merger func(To, To) To) (out To, err error) {
	l := len(elements)
	if l >= 1 {
		if out, err = converter(elements[0]); err != nil {
			return out, err
		}
	}
	if l > 1 {
		for _, e := range elements[1:] {
			c, err := converter(e)
			if err != nil {
				return out, err
			}
			out = merger(out, c)
		}
	}
	return out, nil
}
