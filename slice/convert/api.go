package convert

import (
	"github.com/m4gshm/gollections/not"
	"github.com/m4gshm/gollections/slice"
)

func AndConvert[FS ~[]From, From, To, Too any](elements FS, firsConverter func(From) To, secondConverter func(To) Too) []Too {
	return slice.Convert(slice.Convert(elements, firsConverter), secondConverter)
}

func AndFilter[FS ~[]From, From, To any](elements FS, converter func(From) To, filter func(To) bool) []To {
	return slice.ConvertAndFilter(elements, converter, filter)
}

func NotNil[FS ~[]*From, From, To any](elements FS, converter func(*From) To) []To {
	return slice.FilterAndConvert(elements, not.Nil[From], converter)
}

func ToNotNil[FS ~[]From, From, To any](elements FS, converter func(From) *To) []*To {
	return slice.ConvertCheck(elements, func(f From) (*To, bool) {
		if t := converter(f); t != nil {
			return t, true
		}
		return nil, false
	})
}

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

func Check[FS ~[]From, From, To any](elements FS, converter func(from From) (To, bool)) []To {
	return slice.ConvertCheck(elements, converter)
}

func CheckIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) (To, bool)) []To {
	return slice.ConvertCheckIndexed(elements, by)
}

func Indexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) To) []To {
	return slice.ConvertIndexed(elements, by)
}
