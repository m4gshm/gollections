package map_

import (
	"fmt"
	"strings"

	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kvit"
	"github.com/m4gshm/gollections/loop"
)

// ErrBreak is For, Track breaker
var ErrBreak = loop.ErrBreak

// Of creates a map from a slice of key/value pairs.
func Of[K comparable, V any](elements ...c.KV[K, V]) map[K]V {
	var (
		uniques = make(map[K]V, len(elements))
	)
	for _, kv := range elements {
		key := kv.Key()
		val := kv.Value()
		uniques[key] = val
	}
	return uniques
}

// OfLoop builds a map by iterating key\value pairs of a source.
// The hasNext specifies a predicate that tests existing of a next pair in the source.
// The getNext extracts the pair.
func OfLoop[S any, K comparable, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) (map[K]V, error) {
	return OfLoopResolv(source, hasNext, getNext, kvit.FirstVal[K, V])
}

// OfLoopResolv builds a map by iterating elements of a source.
// The hasNext specifies a predicate that tests existing of a next pair in the source.
// The getNext extracts the element.
// The resolv values for duplicated keys.
func OfLoopResolv[S any, K comparable, E, V any](source S, hasNext func(S) bool, getNext func(S) (K, E, error), resolv func(bool, K, V, E) V) (map[K]V, error) {
	r := map[K]V{}
	for hasNext(source) {
		k, elem, err := getNext(source)
		if err != nil {
			return r, err
		}
		existVal, ok := r[k]
		r[k] = resolv(ok, k, existVal, elem)
	}
	return r, nil
}

// GroupOfLoop builds a map of slices by iterating over elements, extracting key\value pairs and grouping the values for each key in the slices.
// The hasNext specifies a predicate that tests existing of a next pair in the source.
// The getNext extracts the pair.
func GroupOfLoop[S any, K comparable, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) (map[K][]V, error) {
	return OfLoopResolv(source, hasNext, getNext, func(exists bool, key K, elements []V, val V) []V {
		return append(elements, val)
	})
}

// Generate builds a map by an generator function.
// The next returns an key\value pair, or false if the generation is over, or an error.
func Generate[K comparable, V any](next func() (K, V, bool, error)) (map[K]V, error) {
	return GenerateResolv(next, kvit.FirstVal[K, V])
}

// GenerateResolv builds a map by an generator function.
// The next returns an key\value pair, or false if the generation is over, or an error.
// The resolv selects value for duplicated keys.
func GenerateResolv[K comparable, V any](next func() (K, V, bool, error), resolv func(bool, K, V, V) V) (map[K]V, error) {
	r := map[K]V{}
	for {
		k, v, ok, err := next()
		if err != nil || !ok {
			return r, err
		}
		ov, ok := r[k]
		r[k] = resolv(ok, k, ov, v)
	}
}

// Clone makes a copy of a map, copies the values as is
func Clone[M ~map[K]V, K comparable, V any](elements M) M {
	return DeepClone(elements, as.Is[V])
}

// DeepClone copies map values using a copier function to a new map and returns it
func DeepClone[M ~map[K]V, K comparable, V any](elements M, valCopier func(V) V) M {
	return ConvertValues(elements, valCopier)
}

// ConvertValues creates a map with converted values
func ConvertValues[V, Vto any, K comparable, M ~map[K]V](elements M, by func(V) Vto) map[K]Vto {
	converted := make(map[K]Vto, len(elements))
	for key, val := range elements {
		converted[key] = by(val)
	}
	return converted
}

// Keys makes a slice of map keys
func Keys[K comparable, V any, M ~map[K]V](elements M) []K {
	keys := make([]K, 0, len(elements))
	for key := range elements {
		keys = append(keys, key)
	}
	return keys
}

// Values makes a slice of map values
func Values[V any, K comparable, M ~map[K]V](elements M) []V {
	return ValuesConverted(elements, as.Is[V])
}

// ValuesConverted makes a slice of converted map values
func ValuesConverted[M ~map[K]V, K comparable, V, Vto any](elements M, by func(V) Vto) []Vto {
	values := make([]Vto, 0, len(elements))
	for _, val := range elements {
		values = append(values, by(val))
	}
	return values
}

// Track applies a tracker for every key/value pairs from a map. To stop traking just return the ErrBreak
func Track[M ~map[K]V, K comparable, V any](elements M, tracker func(K, V) error) error {
	for key, val := range elements {
		if err := tracker(key, val); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// TrackEach applies a tracker for every key/value pairs from a map
func TrackEach[M ~map[K]V, K comparable, V any](elements M, tracker func(K, V)) {
	for key, val := range elements {
		tracker(key, val)
	}
}

// For applies a walker for every key/value pairs from a map. Key/value pair is boxed to the KV. To stop walking just return the ErrBreak
func For[M ~map[K]V, K comparable, V any](elements M, walker func(c.KV[K, V]) error) error {
	for key, val := range elements {
		if err := walker(c.NewKV(key, val)); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEach applies a walker for every key/value pairs from a map. Key/value pair is boxed to the KV
func ForEach[M ~map[K]V, K comparable, V any](elements M, walker func(c.KV[K, V])) {
	for key, val := range elements {
		walker(c.NewKV(key, val))
	}
}

// TrackOrdered applies a tracker for every key/value pairs from a map in order. To stop traking just return the ErrBreak
func TrackOrdered[M ~map[K]V, K comparable, V any](order []K, uniques M, tracker func(K, V) error) error {
	for _, key := range order {
		if err := tracker(key, uniques[key]); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// TrackEachOrdered applies a tracker for every key/value pairs from a map in order
func TrackEachOrdered[M ~map[K]V, K comparable, V any](elements []K, uniques M, tracker func(K, V)) {
	for _, key := range elements {
		tracker(key, uniques[key])
	}
}

// ForOrdered applies a walker for every key/value pairs from a map in order. Key/value pair is boxed to the KV. To stop walking just return the ErrBreak
func ForOrdered[M ~map[K]V, K comparable, V any](elements []K, uniques M, walker func(c.KV[K, V]) error) error {
	for _, key := range elements {
		if err := walker(c.NewKV(key, uniques[key])); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEachOrdered applies a walker for every key/value pairs from a map in order. Key/value pair is boxed to the KV
func ForEachOrdered[M ~map[K]V, K comparable, V any](elements []K, uniques M, walker func(c.KV[K, V])) {
	for _, key := range elements {
		walker(c.NewKV(key, uniques[key]))
	}
}

// ForKeys applies a walker for every key from a map. To stop walking just return the ErrBreak
func ForKeys[M ~map[K]V, K comparable, V any](elements M, walker func(K) error) error {
	for key := range elements {
		if err := walker(key); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEachKey applies a walker for every key from a map
func ForEachKey[M ~map[K]V, K comparable, V any](elements M, walker func(K)) {
	for key := range elements {
		walker(key)
	}
}

// ForValues applies a walker for every value from a map. To stop walking just return the ErrBreak
func ForValues[M ~map[K]V, K comparable, V any](elements M, walker func(V) error) error {
	for _, val := range elements {
		if err := walker(val); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEachValue applies a walker for every value from a map
func ForEachValue[M ~map[K]V, K comparable, V any](elements M, walker func(V)) {
	for _, val := range elements {
		walker(val)
	}
}

// ForOrderedValues applies a walker for every value from a map in order. To stop walking just return the ErrBreak
func ForOrderedValues[M ~map[K]V, K comparable, V any](elements []K, uniques M, walker func(V) error) error {
	for _, key := range elements {
		val := uniques[key]
		if err := walker(val); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEachOrderedValues applies a walker for every value from a map in order
func ForEachOrderedValues[M ~map[K]V, K comparable, V any](elements []K, uniques M, walker func(V)) {
	for _, key := range elements {
		val := uniques[key]
		walker(val)
	}
}

// ToStringOrdered converts elements to the string representation according to the order
func ToStringOrdered[M ~map[K]V, K comparable, V any](order []K, elements M) string {
	return ToStringOrderedf(order, elements, "%+v:%+v", " ")
}

// ToStringOrderedf converts elements to a string representation using a key/value pair format and a delimeter. In order
func ToStringOrderedf[M ~map[K]V, K comparable, V any](order []K, elements M, kvFormat, delim string) string {
	str := strings.Builder{}
	str.WriteString("[")
	for i, K := range order {
		if i > 0 {
			_, _ = str.WriteString(delim)
		}
		str.WriteString(fmt.Sprintf(kvFormat, K, elements[K]))
	}
	str.WriteString("]")
	return str.String()
}

// ToString converts elements to the string representation
func ToString[M ~map[K]V, K comparable, V any](elements M) string {
	return ToStringf(elements, "%+V:%+V", " ")
}

// ToStringf converts elements to a string representation using a key/value pair format and a delimeter
func ToStringf[M ~map[K]V, K comparable, V any](elements M, kvFormat, delim string) string {
	str := strings.Builder{}
	str.WriteString("[")
	i := 0
	for K, V := range elements {
		if i > 0 {
			_, _ = str.WriteString(delim)
		}
		str.WriteString(fmt.Sprintf(kvFormat, K, V))
		i++
	}
	str.WriteString("]")
	return str.String()
}
