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
