package map_

import (
	"fmt"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/check"
	impl "github.com/m4gshm/gollections/slice/impl/slice"
	"github.com/m4gshm/gollections/typ"
)

//Copy - makes a new slice with copied elements.
func Copy[m map[k]v, k comparable, v any](elements m) m {
	var copied m
	for key, val := range elements {
		copied[key] = val
	}
	return copied
}

func Track[m map[k]v, k comparable, v any](elements m, tracker func(k, v) error) error {
	for key, val := range elements {
		if err := tracker(key, val); err != nil {
			return err
		}
	}
	return nil
}

func TrackEach[m map[k]v, k comparable, v any](elements m, tracker func(k, v)) error {
	for key, val := range elements {
		tracker(key, val)
	}
	return nil
}

func For[m map[k]v, k comparable, v any](elements m, walker func(*typ.KV[k, v]) error) error {
	for key, val := range elements {
		if err := walker(K.V(key, val)); err != nil {
			return err
		}
	}
	return nil
}

func ForEach[m map[k]v, k comparable, v any](elements m, walker func(*typ.KV[k, v])) error {
	for key, val := range elements {
		walker(K.V(key, val))
	}
	return nil
}

func TrackOrdered[m map[k]v, k comparable, v any](elements []*k, uniques m, tracker func(k, v) error) error {
	for _, ref := range elements {
		key := *ref
		if err := tracker(key, uniques[key]); err != nil {
			return err
		}
	}
	return nil
}

func TrackEachOrdered[m map[k]v, k comparable, v any](elements []*k, uniques m, tracker func(k, v)) error {
	for _, ref := range elements {
		key := *ref
		tracker(key, uniques[key])
	}
	return nil
}

func ForOrdered[m map[k]v, k comparable, v any](elements []*k, uniques m, walker func(*typ.KV[k, v]) error) error {
	for _, ref := range elements {
		key := *ref
		if err := walker(K.V(key, uniques[key])); err != nil {
			return err
		}
	}
	return nil
}

func ForEachOrdered[m map[k]v, k comparable, v any](elements []*k, uniques m, walker func(*typ.KV[k, v])) error {
	for _, ref := range elements {
		key := *ref
		walker(K.V(key, uniques[key]))
	}
	return nil
}

func ForKeys[m map[k]v, k comparable, v any](elements m, walker func(k) error) error {
	for key := range elements {
		if err := walker(key); err != nil {
			return err
		}
	}
	return nil
}

func ForEachKey[m map[k]v, k comparable, v any](elements m, walker func(k)) error {
	for key := range elements {
		walker(key)
	}
	return nil
}

//Map creates a lazy Iterator that converts elements with a converter and returns them.
func Map[From, To any](elements []From, by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements, by)
}

// additionally filters 'From' elements by filters.
func MapFit[From, To any](elements []From, fit typ.Predicate[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](elements []From, by typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements, by)
}

// additionally checks 'From' elements by fit.
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
