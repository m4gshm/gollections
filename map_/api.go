package map_

import (
	"bytes"
	"fmt"
)

//Copy makes a map copy.
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

func For[m map[k]v, k comparable, v any](elements m, walker func(*KV[k, v]) error) error {
	for key, val := range elements {
		if err := walker(NewKV(key, val)); err != nil {
			return err
		}
	}
	return nil
}

func ForEach[m map[k]v, k comparable, v any](elements m, walker func(*KV[k, v])) {
	for key, val := range elements {
		walker(NewKV(key, val))
	}
}

func TrackOrdered[m map[k]v, k comparable, v any](elements []k, uniques m, tracker func(k, v) error) error {
	for _, key := range elements {
		if err := tracker(key, uniques[key]); err != nil {
			return err
		}
	}
	return nil
}

func TrackEachOrdered[m map[k]v, k comparable, v any](elements []k, uniques m, tracker func(k, v)) {
	for _, key := range elements {
		tracker(key, uniques[key])
	}
}

func ForOrdered[m map[k]v, k comparable, v any](elements []k, uniques m, walker func(*KV[k, v]) error) error {
	for _, key := range elements {
		if err := walker(NewKV(key, uniques[key])); err != nil {
			return err
		}
	}
	return nil
}

func ForEachOrdered[m map[k]v, k comparable, v any](elements []k, uniques m, walker func(*KV[k, v])) {
	for _, key := range elements {
		walker(NewKV(key, uniques[key]))
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

func ForOrderedValues[m map[k]v, k comparable, v any](elements []k, uniques m, walker func(v) error) error {
	for _, key := range elements {
		val := uniques[key]
		if err := walker(val); err != nil {
			return err
		}
	}
	return nil
}

func ForEachOrderedValues[m map[k]v, k comparable, v any](elements []k, uniques m, walker func(v)) {
	for _, key := range elements {
		val := uniques[key]
		walker(val)
	}
}

//ToStringOrdered converts elements to the string representation according to the order
func ToStringOrdered[k comparable, v any](order []k, elements map[k]v) string {
	return ToStringOrderedf(order, elements, "%+v:%+v", " ")
}

func ToStringOrderedf[k comparable, v any](order []k, elements map[k]v, kvFormat, delim string) string {
	str := bytes.Buffer{}
	str.WriteString("[")
	for i, k := range order {
		if i > 0 {
			_, _ = str.WriteString(delim)
		}
		str.WriteString(fmt.Sprintf(kvFormat, k, elements[k]))
	}
	str.WriteString("]")
	return str.String()
}

//ToString converts elements to the string representation
func ToString[k comparable, v any](elements map[k]v) string {
	return ToStringf(elements, "%+v:%+v", " ")
}

func ToStringf[k comparable, v any](elements map[k]v, kvFormat, delim string) string {
	str := bytes.Buffer{}
	str.WriteString("[")
	i := 0
	for k, v := range elements {
		if i > 0 {
			_, _ = str.WriteString(delim)
		}
		str.WriteString(fmt.Sprintf(kvFormat, k, v))
		i++
	}
	str.WriteString("]")
	return str.String()
}
