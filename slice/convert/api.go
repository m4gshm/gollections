// Package convert provides slice converation helpers
package convert

import (
	"github.com/m4gshm/gollections/not"
	"github.com/m4gshm/gollections/slice"
)

// AndConvert - convert.AndConvert makes double converts From->Intermediate->To of the elements
func AndConvert[FS ~[]From, From, I, To any](elements FS, firsConverter func(From) I, secondConverter func(I) To) []To {
	return slice.Convert(elements, func(from From) To { return secondConverter(firsConverter(from)) })
}

// AndFilter - convert.AndFilter converts only filtered elements and returns them
func AndFilter[FS ~[]From, From, To any](elements FS, converter func(From) To, filter func(To) bool) []To {
	return slice.ConvertAndFilter(elements, converter, filter)
}

// NotNil - convert.NotNil converts only not nil elements and returns them
func NotNil[FS ~[]*From, From, To any](elements FS, converter func(*From) To) []To {
	return slice.FilterAndConvert(elements, not.Nil[From], converter)
}

// ToNotNil - convert.ToNotNil converts elements and returns only not nil converted elements
func ToNotNil[FS ~[]From, From, To any](elements FS, converter func(From) *To) []*To {
	return slice.ConvertCheck(elements, func(f From) (*To, bool) {
		if t := converter(f); t != nil {
			return t, true
		}
		return nil, false
	})
}

// NilSafe - convert.NilSafe filters not nil elements, converts that ones, filters not nils after converting and returns them
func NilSafe[FS ~[]*From, From, To any](elements FS, converter func(*From) *To) []*To {
	return slice.ConvertCheck(elements, func(f *From) (*To, bool) {
		if f != nil {
			if t := converter(f); t != nil {
				return t, true
			}
		}
		return nil, false
	})
}

// Check - convert.Check is a short alias of slice.ConvertCheck
func Check[FS ~[]From, From, To any](elements FS, converter func(from From) (To, bool)) []To {
	return slice.ConvertCheck(elements, converter)
}

// CheckIndexed - convert.CheckIndexed is a short alias of slice.ConvertCheckIndexed
func CheckIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) (To, bool)) []To {
	return slice.ConvertCheckIndexed(elements, by)
}

// Indexed - convert.Indexed is a short alias of slice.ConvertIndexed
func Indexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) To) []To {
	return slice.ConvertIndexed(elements, by)
}
