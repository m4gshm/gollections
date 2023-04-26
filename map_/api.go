// Package map_ provides map processing helper functions
package map_

import (
	"fmt"
	"strings"

	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/map_/resolv"
)

// ErrBreak is For, Track breaker
var ErrBreak = c.ErrBreak

func New[K comparable, V any]() map[K]V {
	return map[K]V{}
}

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
	return OfLoopResolv(source, hasNext, getNext, resolv.FirstVal[K, V])
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
	return GenerateResolv(next, resolv.FirstVal[K, V])
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

// Clone makes a copy of the 'elements' map. The values are copied as is.
func Clone[M ~map[K]V, K comparable, V any](elements M) M {
	return DeepClone(elements, as.Is[V])
}

// DeepClone makes a copy of the 'elements' map. The values are copied by the 'copier' function
func DeepClone[M ~map[K]V, K comparable, V any](elements M, copier func(V) V) M {
	return ConvertValues(elements, copier)
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
	values := make([]V, 0, len(elements))
	for _, val := range elements {
		values = append(values, val)
	}
	return values
}

// ValuesConverted makes a slice of converted map values
func ValuesConverted[M ~map[K]V, K comparable, V, Vto any](elements M, by func(V) Vto) []Vto {
	values := make([]Vto, 0, len(elements))
	for _, val := range elements {
		values = append(values, by(val))
	}
	return values
}

// Track applies the 'tracker' function for every key/value pairs from the 'elements' map. Return the c.ErrBreak to stop
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

// TrackEach applies the 'tracker' function for every key/value pairs from the 'elements' map
func TrackEach[M ~map[K]V, K comparable, V any](elements M, tracker func(K, V)) {
	for key, val := range elements {
		tracker(key, val)
	}
}

// For applies the 'walker' function for key/value pairs from the elements. Return the c.ErrBreak to stop.
func For[M ~map[K]V, K comparable, V any](elements M, walker func(c.KV[K, V]) error) error {
	for key, val := range elements {
		if err := walker(kv.New(key, val)); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEach applies the 'walker' function for every key/value pair from the elements map
func ForEach[M ~map[K]V, K comparable, V any](elements M, walker func(c.KV[K, V])) {
	for key, val := range elements {
		walker(kv.New(key, val))
	}
}

// TrackOrdered applies the 'tracker' function for key/value pairs from the 'elements' map in order of the 'order' slice. Return the c.ErrBreak to stop
func TrackOrdered[M ~map[K]V, K comparable, V any](order []K, elements M, tracker func(K, V) error) error {
	for _, key := range order {
		if err := tracker(key, elements[key]); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// TrackEachOrdered applies the 'tracker' function for evey key/value pair from the 'elements' map in order of the 'order' slice
func TrackEachOrdered[M ~map[K]V, K comparable, V any](order []K, uniques M, tracker func(K, V)) {
	for _, key := range order {
		tracker(key, uniques[key])
	}
}

// ForOrdered applies the 'walker' function for every key/value pair from the 'elements' map in order of the 'order' slice. Return the c.ErrBreak to stop.
func ForOrdered[M ~map[K]V, K comparable, V any](order []K, elements M, walker func(c.KV[K, V]) error) error {
	for _, key := range order {
		if err := walker(kv.New(key, elements[key])); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEachOrdered applies the 'walker' function for every key/value pair from the 'elements' map in order of the 'order' slice.
func ForEachOrdered[M ~map[K]V, K comparable, V any](order []K, elements M, walker func(c.KV[K, V])) {
	for _, key := range order {
		walker(kv.New(key, elements[key]))
	}
}

// ForKeys applies the 'walker' function for keys from the 'elements' map . Return the c.ErrBreak to stop.
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

// ForEachKey applies the 'walker' function for every key from from the 'elements' map
func ForEachKey[M ~map[K]V, K comparable, V any](elements M, walker func(K)) {
	for key := range elements {
		walker(key)
	}
}

// ForValues applies the 'walker' function for values from the 'elements' map . Return the c.ErrBreak to stop..
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

// ForEachValue applies the 'walker' function for every value from from the 'elements' map
func ForEachValue[M ~map[K]V, K comparable, V any](elements M, walker func(V)) {
	for _, val := range elements {
		walker(val)
	}
}

// ForOrderedValues applies the 'walker' function for values from the 'elements' map in order of the 'order' slice. Return the c.ErrBreak to stop..
func ForOrderedValues[M ~map[K]V, K comparable, V any](order []K, elements M, walker func(V) error) error {
	for _, key := range order {
		val := elements[key]
		if err := walker(val); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEachOrderedValues applies the 'walker' function for each value from the 'elements' map in order of the 'order' slice
func ForEachOrderedValues[M ~map[K]V, K comparable, V any](order []K, elements M, walker func(V)) {
	for _, key := range order {
		val := elements[key]
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

// Reduce reduces the key/value pairs by the 'next' function into an one pair using the 'merge' function
func Reduce[M ~map[K]V, K comparable, V any](elements M, merge func(K, V, K, V) (K, V)) (rk K, rv V) {
	first := true
	for k, v := range elements {
		if first {
			rk, rv = k, v
			first = false
		} else {
			rk, rv = merge(rk, rv, k, v)
		}
	}
	return rk, rv
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func HasAny[M ~map[K]V, K comparable, V any](elements M, predicate func(K, V) bool) bool {
	for k, v := range elements {
		if predicate(k, v) {
			return true
		}
	}
	return false
}
