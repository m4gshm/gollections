// Package convert provides loop converation helpers
package convert

import (
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op/check/not"
)

// AndConvert - convert.AndConvert makes double converts From->Intermediate->To of the elements
func AndConvert[From, I, To any](next func() (From, bool), firsConverter func(From) I, secondConverter func(I) To) loop.Loop[To] {
	return loop.Convert(next, func(from From) To { return secondConverter(firsConverter(from)) })
}

// AndFilter - convert.AndFilter converts only filtered elements and returns them
func AndFilter[From, To any](next func() (From, bool), converter func(From) To, filter func(To) bool) loop.Loop[To] {
	return loop.ConvertAndFilter(next, converter, filter)
}

// NotNil - convert.NotNil converts only not nil elements and returns them
func NotNil[From, To any](next func() (*From, bool), converter func(*From) To) loop.Loop[To] {
	return loop.FilterAndConvert(next, not.Nil[From], converter)
}

// ToNotNil - convert.ToNotNil converts elements and returns only not nil converted elements
func ToNotNil[From, To any](next func() (From, bool), converter func(From) *To) loop.Loop[*To] {
	return loop.ConvertOK(next, func(f From) (*To, bool) {
		if t := converter(f); t != nil {
			return t, true
		}
		return nil, false
	})
}

// NilSafe - convert.NilSafe filters not nil next, converts that ones, filters not nils after converting and returns them
func NilSafe[From, To any](next func() (*From, bool), converter func(*From) *To) loop.Loop[*To] {
	return loop.ConvertOK(next, func(f *From) (*To, bool) {
		if f != nil {
			if t := converter(f); t != nil {
				return t, true
			}
		}
		return nil, false
	})
}

// Check - convert.Check is a short alias of loop.ConvertOK
func Check[From, To any](next func() (From, bool), converter func(from From) (To, bool)) loop.Loop[To] {
	return loop.ConvertOK(next, converter)
}

// FromIndexed - convert.FromIndexed retrieves elements from a indexed source and converts them
func FromIndexed[From, To any](amount int, next func(int) From, converter func(from From) To) loop.Loop[To] {
	return loop.Convert(loop.OfIndexed(amount, next), converter)
}

// AndReduce - convert.AndReduce converts elements and merge them into one
func AndReduce[From, To any](next func() (From, bool), converter func(From) To, merge func(To, To) To) (out To) {
	return loop.ConvertAndReduce(next, converter, merge)
}
