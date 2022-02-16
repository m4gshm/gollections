package map_

import (
	"errors"
	"fmt"
	"strings"
)

//Break is For, Track breaker
var Break = errors.New("Break")

//Copy makes a map copy.
func Copy[m map[k]v, k comparable, v any](elements m) m {
	var copied m
	for key, val := range elements {
		copied[key] = val
	}
	return copied
}

//Track applies a tracker for every key/value pairs from a map. To stop traking just return the Break.
func Track[m map[k]v, k comparable, v any](elements m, tracker func(k, v) error) error {
	for key, val := range elements {
		if err := tracker(key, val); err == Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//TrackEach applies a tracker for every key/value pairs from a map.
func TrackEach[m map[k]v, k comparable, v any](elements m, tracker func(k, v)) {
	for key, val := range elements {
		tracker(key, val)
	}
}

//For applies a walker for every key/value pairs from a map. Key/value pair is boxed to the KV. To stop walking just return the Break.
func For[m map[k]v, k comparable, v any](elements m, walker func(*KV[k, v]) error) error {
	for key, val := range elements {
		if err := walker(NewKV(key, val)); err == Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//ForEach applies a walker for every key/value pairs from a map. Key/value pair is boxed to the KV.
func ForEach[m map[k]v, k comparable, v any](elements m, walker func(*KV[k, v])) {
	for key, val := range elements {
		walker(NewKV(key, val))
	}
}

//TrackOrdered applies a tracker for every key/value pairs from a map in order. To stop traking just return the Break.
func TrackOrdered[m map[k]v, k comparable, v any](order []k, uniques m, tracker func(k, v) error) error {
	for _, key := range order {
		if err := tracker(key, uniques[key]); err == Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//TrackEachOrdered applies a tracker for every key/value pairs from a map in order.
func TrackEachOrdered[m map[k]v, k comparable, v any](elements []k, uniques m, tracker func(k, v)) {
	for _, key := range elements {
		tracker(key, uniques[key])
	}
}

//ForOrdered applies a walker for every key/value pairs from a map in order. Key/value pair is boxed to the KV. To stop walking just return the Break.
func ForOrdered[m map[k]v, k comparable, v any](elements []k, uniques m, walker func(*KV[k, v]) error) error {
	for _, key := range elements {
		if err := walker(NewKV(key, uniques[key])); err == Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//ForEachOrdered applies a walker for every key/value pairs from a map in order. Key/value pair is boxed to the KV.
func ForEachOrdered[m map[k]v, k comparable, v any](elements []k, uniques m, walker func(*KV[k, v])) {
	for _, key := range elements {
		walker(NewKV(key, uniques[key]))
	}
}

//For applies a walker for every key from a map. To stop walking just return the Break.
func ForKeys[m map[k]v, k comparable, v any](elements m, walker func(k) error) error {
	for key := range elements {
		if err := walker(key); err == Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//ForEachKey applies a walker for every key from a map.
func ForEachKey[m map[k]v, k comparable, v any](elements m, walker func(k)) {
	for key := range elements {
		walker(key)
	}
}

//ForValues applies a walker for every value from a map. To stop walking just return the Break.
func ForValues[m map[k]v, k comparable, v any](elements m, walker func(v) error) error {
	for _, val := range elements {
		if err := walker(val); err == Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//ForEachValue applies a walker for every value from a map.
func ForEachValue[m map[k]v, k comparable, v any](elements m, walker func(v)) {
	for _, val := range elements {
		walker(val)
	}
}

//ForOrderedValues applies a walker for every value from a map in order. To stop walking just return the Break.
func ForOrderedValues[m map[k]v, k comparable, v any](elements []k, uniques m, walker func(v) error) error {
	for _, key := range elements {
		val := uniques[key]
		if err := walker(val); err == Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

//ForEachOrderedValues applies a walker for every value from a map in order.
func ForEachOrderedValues[m map[k]v, k comparable, v any](elements []k, uniques m, walker func(v)) {
	for _, key := range elements {
		val := uniques[key]
		walker(val)
	}
}

//ToStringOrdered converts elements to the string representation according to the order.
func ToStringOrdered[k comparable, v any](order []k, elements map[k]v) string {
	return ToStringOrderedf(order, elements, "%+v:%+v", " ")
}

//ToStringOrderedf converts elements to a string representation using a key/value pair format and a delimeter. In order.
func ToStringOrderedf[k comparable, v any](order []k, elements map[k]v, kvFormat, delim string) string {
	str := strings.Builder{}
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

//ToString converts elements to the string representation.
func ToString[k comparable, v any](elements map[k]v) string {
	return ToStringf(elements, "%+v:%+v", " ")
}

//ToStringf converts elements to a string representation using a key/value pair format and a delimeter.
func ToStringf[k comparable, v any](elements map[k]v, kvFormat, delim string) string {
	str := strings.Builder{}
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
