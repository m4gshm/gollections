package filter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/slice"
)

func AndConvert[FS ~[]From, From, To any](elements FS, filter predicate.Predicate[From], converter c.Converter[From, To]) []To {
	return slice.FilterAndConvert(elements, filter, converter)
}

func ConvertFilter[FS ~[]From, From, To any](elements FS, filter predicate.Predicate[From], converter c.Converter[From, To], filterTo predicate.Predicate[To]) []To {
	return slice.FilterConvertFilter(elements, filter, converter, filterTo)
}
