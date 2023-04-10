package convert

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/not"
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/slice"
)

func Fit[FS ~[]From, From, To any](elements FS, by c.Converter[From, To], fit predicate.Predicate[From]) []To {
	return slice.ConvertFit(elements, fit, by)
}

func NotNil[FS ~[]*From, From, To any](elements FS, by c.Converter[*From, To]) []To {
	return Fit(elements, by, not.Nil[From])
}

func Check[FS ~[]From, From, To any](elements FS, by func(from From) (To, bool)) []To {
	return slice.ConvertCheck(elements, by)
}

func CheckIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) (To, bool)) []To {
	return slice.ConvertCheckIndexed(elements, by)
}

func Indexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) To) []To {
	return slice.ConvertIndexed(elements, by)
}
