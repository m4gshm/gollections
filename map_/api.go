// Package map_ provides map processing helper functions
package map_

import (
	"fmt"
	"strings"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/map_/resolv"
)

// ErrBreak is For, Track breaker
var ErrBreak = c.ErrBreak

// Of instantiates a ap from the specified key/value pairs
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
	return OfLoopResolv(source, hasNext, getNext, resolv.First[K, V])
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
// The next returns a key\value pair, or false if the generation is over, or an error.
func Generate[K comparable, V any](next func() (K, V, bool, error)) (map[K]V, error) {
	return GenerateResolv(next, resolv.First[K, V])
}

// GenerateResolv builds a map by an generator function.
// The next returns a key\value pair, or false if the generation is over, or an error.
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
func ConvertValues[M ~map[K]V, V, Vto any, K comparable](elements M, converter func(V) Vto) map[K]Vto {
	converted := make(map[K]Vto, len(elements))
	for key, val := range elements {
		converted[key] = converter(val)
	}
	return converted
}

// ConvertKeys creates a map with converted keys
func ConvertKeys[M ~map[K]V, K, Kto comparable, V any](elements M, converter func(K) Kto) map[Kto]V {
	converted := make(map[Kto]V, len(elements))
	for key, val := range elements {
		converted[converter(key)] = val
	}
	return converted
}

// Convert creates a map with converted keys and values
func Convert[M ~map[K]V, K, Kto comparable, V, Vto any](elements M, converter func(K, V) (Kto, Vto)) map[Kto]Vto {
	converted := make(map[Kto]Vto, len(elements))
	for key, val := range elements {
		kto, vto := converter(key, val)
		converted[kto] = vto
	}
	return converted
}

// Conv creates a map with converted keys and values
func Conv[M ~map[K]V, K, Kto comparable, V, Vto any](elements M, converter func(K, V) (Kto, Vto, error)) (map[Kto]Vto, error) {
	converted := make(map[Kto]Vto, len(elements))
	for key, val := range elements {
		kto, vto, err := converter(key, val)
		if err != nil {
			return converted, err
		}
		converted[kto] = vto
	}
	return converted, nil
}

// Filter creates a map containing only the filtered elements
func Filter[M ~map[K]V, K comparable, V any](elements M, filter func(K, V) bool) map[K]V {
	filtered := map[K]V{}
	for key, val := range elements {
		if filter(key, val) {
			filtered[key] = val
		}
	}
	return filtered
}

// FilterKeys creates a map containing only the filtered elements
func FilterKeys[M ~map[K]V, K comparable, V any](elements M, filter func(K) bool) map[K]V {
	filtered := map[K]V{}
	for key, val := range elements {
		if filter(key) {
			filtered[key] = val
		}
	}
	return filtered
}

// FilterValues creates a map containing only the filtered elements
func FilterValues[M ~map[K]V, K comparable, V any](elements M, filter func(V) bool) map[K]V {
	filtered := map[K]V{}
	for key, val := range elements {
		if filter(val) {
			filtered[key] = val
		}
	}
	return filtered
}

// Keys returns keys of the 'elements' map as a slice
func Keys[M ~map[K]V, K comparable, V any](elements M) []K {
	if elements == nil {
		return nil
	}
	return AppendKeys(elements, make([]K, 0, len(elements)))
}

// AppendKeys collects keys of the specified 'elements' map into the specified 'out' slice
func AppendKeys[M ~map[K]V, K comparable, V any](elements M, out []K) []K {
	for key := range elements {
		out = append(out, key)
	}
	return out
}

// Values returns values of the 'elements' map as a slice
func Values[M ~map[K]V, K comparable, V any](elements M) []V {
	if elements == nil {
		return nil
	}
	return AppendValues(elements, make([]V, 0, len(elements)))
}

// AppendValues collects values of the specified 'elements' map into the specified 'out' slice
func AppendValues[M ~map[K]V, K comparable, V any](elements M, out []V) []V {
	for _, val := range elements {
		out = append(out, val)
	}
	return out
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
func Reduce[M ~map[K]V, K comparable, V any](elements M, merge func(K, K, V, V) (K, V)) (rk K, rv V) {
	first := true
	for k, v := range elements {
		if first {
			rk, rv = k, v
			first = false
		} else {
			rk, rv = merge(rk, k, rv, v)
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

// ToSlice collects key\value elements to a slice by applying the specified converter to evety element
func ToSlice[M ~map[K]V, K comparable, V any, T any](elements M, converter func(key K, val V) T) []T {
	out := make([]T, 0, len(elements))
	for key, val := range elements {
		out = append(out, converter(key, val))
	}
	return out
}

// ToSlicee collects key\value elements to a slice by applying the specified erroreable converter to evety element
func ToSlicee[M ~map[K]V, K comparable, V any, T any](elements M, converter func(key K, val V) (T, error)) ([]T, error) {
	out := make([]T, 0, len(elements))
	for key, val := range elements {
		t, err := converter(key, val)
		if err != nil {
			return out, err
		}
		out = append(out, t)
	}
	return out, nil
}

// Empty checks the val map is empty
func Empty[M ~map[K]V, K comparable, V any](val M) bool {
	return len(val) == 0
}

// NotEmpty checks the val map is not empty
func NotEmpty[M ~map[K]V, K comparable, V any](val M) bool {
	return !Empty(val)
}

// Get returns the value by the specified key from the map m or zero if the map doesn't contain that key
func Get[M ~map[K]V, K comparable, V any](m M, key K) V {
	val, _ := GetOk(m, key)
	return val
}

// GetOk returns the value, and true by the specified key from the map m or zero and false if the map doesn't contain that key
func GetOk[M ~map[K]V, K comparable, V any](m M, key K) (val V, ok bool) {
	if m != nil {
		val, ok = m[key]
	}
	return val, ok
}

// Getter creates a function that can be used for retrieving a value from the map m by a key
func Getter[M ~map[K]V, K comparable, V any](m M) func(key K) V {
	return func(key K) V { return Get(m, key) }
}

// GetterOk creates a function that can be used for retrieving a value from the map m by a key
func GetterOk[M ~map[K]V, K comparable, V any](m M) func(key K) (V, bool) {
	return func(key K) (V, bool) { return GetOk(m, key) }
}

// Contains checks is the map contains a key
func Contains[M ~map[K]V, K comparable, V any](m M, key K) (ok bool) {
	if m != nil {
		_, ok = m[key]
	}
	return ok
}

// KeyChecker creates a function that can be used to check if the map contains a key
func KeyChecker[M ~map[K]V, K comparable, V any](m M) func(key K) (ok bool) {
	return func(key K) (ok bool) {
		return Contains(m, key)
	}
}
