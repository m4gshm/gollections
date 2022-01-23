package map_

import (
	"github.com/m4gshm/gollections/K"
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

func TrackEach[m map[k]v, k comparable, v any](elements m, tracker func(k, v)) {
	for key, val := range elements {
		tracker(key, val)
	}
}

func For[m map[k]v, k comparable, v any](elements m, walker func(*typ.KV[k, v]) error) error {
	for key, val := range elements {
		if err := walker(K.V(key, val)); err != nil {
			return err
		}
	}
	return nil
}

func ForEach[m map[k]v, k comparable, v any](elements m, walker func(*typ.KV[k, v])) {
	for key, val := range elements {
		walker(K.V(key, val))
	}
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

func TrackEachOrdered[m map[k]v, k comparable, v any](elements []*k, uniques m, tracker func(k, v)) {
	for _, ref := range elements {
		key := *ref
		tracker(key, uniques[key])
	}
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

func ForEachOrdered[m map[k]v, k comparable, v any](elements []*k, uniques m, walker func(*typ.KV[k, v])) {
	for _, ref := range elements {
		key := *ref
		walker(K.V(key, uniques[key]))
	}
}

func ForKeys[m map[k]v, k comparable, v any](elements m, walker func(k) error) error {
	for key := range elements {
		if err := walker(key); err != nil {
			return err
		}
	}
	return nil
}

func ForEachKey[m map[k]v, k comparable, v any](elements m, walker func(k)) {
	for key := range elements {
		walker(key)
	}
}

func ForValues[m map[k]v, k comparable, v any](elements m, walker func(v) error) error {
	for _, val := range elements {
		if err := walker(val); err != nil {
			return err
		}
	}
	return nil
}

func ForEachValue[m map[k]v, k comparable, v any](elements m, walker func(v)) {
	for _, val := range elements {
		walker(val)
	}
}

func ForOrderedValues[m map[k]v, k comparable, v any](elements []*k, uniques m, walker func(v) error) error {
	for _, r := range elements {
		key := *r
		val := uniques[key]
		if err := walker(val); err != nil {
			return err
		}
	}
	return nil
}

func ForEachOrderedValues[m map[k]v, k comparable, v any](elements []*k, uniques m, walker func(v)) {
	for _, r := range elements {
		key := *r
		val := uniques[key]
		walker(val)
	}
}
