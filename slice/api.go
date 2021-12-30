package slice

import (
	"fmt"

	"github.com/m4gshm/container/check"
	impl "github.com/m4gshm/container/slice/impl/slice"
	"github.com/m4gshm/container/typ"
)

func Of[T any](elements ...T) []T {
	return elements
}

//Map creates a lazy Iterator that converts elements with a converter and returns them
func Map[From, To any](elements []From, by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements, by)
}

// additionally filters 'From' elements by filters
func MapFit[From, To any](elements []From, fit typ.Predicate[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any](elements []From, by typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements, by)
}

// additionally checks 'From' elements by fit
func FlattFit[From, To any](elements []From, fit typ.Predicate[From], flatt typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any](elements []T, filter typ.Predicate[T]) typ.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates a lazy Iterator that filters nullable elements
func NotNil[T any](elements []*T) typ.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

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

func ToStringRefs[T any](elements []*T) string {
	str := ""
	for _, ref := range elements {
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
