package slice

import (
	"constraints"
	"fmt"

	"github.com/m4gshm/gollections/check"
	impl "github.com/m4gshm/gollections/slice/impl/slice"
	"github.com/m4gshm/gollections/typ"
)

//Of - constructor.
func Of[T any](elements ...T) []T {
	return elements
}

//Copy - makes a new slice with copied elements.
func Copy[T any](elements []T) []T {
	copied := make([]T, len(elements))
	copy(copied, elements)
	return copied
}

func Delete[T any](index int, elements []T) []T {
	return append(elements[0:index], elements[index+1:]...)
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

func Track[T any](elements []T, tracker func(int, T) error) error {
	for i, e := range elements {
		if err := tracker(i, e); err != nil {
			return err
		}
	}
	return nil
}

func TrackEach[T any](elements []T, tracker func(int, T)) error {
	for i, e := range elements {
		tracker(i, e)
	}
	return nil
}

func For[T any](elements []T, walker func(T) error) error {
	for _, e := range elements {
		if err := walker(e); err != nil {
			return err
		}
	}
	return nil
}

func ForEach[T any](elements []T, walker func(T)) error {
	for _, e := range elements {
		walker(e)
	}
	return nil
}

func ForRefs[T any](elements []*T, walker func(T) error) error {
	for _, e := range elements {
		if err := walker(*e); err != nil {
			return err
		}
	}
	return nil
}

func ForEachRef[T any](elements []*T, walker func(T)) error {
	for _, e := range elements {
		walker(*e)
	}
	return nil
}

//Map creates a lazy Iterator that converts elements with a converter and returns them.
func Map[From, To any](elements []From, by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements, by)
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any](elements []From, fit typ.Predicate[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](elements []From, by typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements, by)
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any](elements []From, fit typ.Predicate[From], flatt typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones.
func Filter[T any](elements []T, filter typ.Predicate[T]) typ.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates a lazy Iterator that filters nullable elements.
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
