package map_

import (
	"fmt"
	"strings"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/kvit"
)

// ErrBreak is For, Track breaker
var ErrBreak = it.ErrBreak

// OfLoop builds a map by iterating key\value pairs of a source.
// The hasNext specifies a predicate that tests existing of a next pair in the source.
// The getNext extracts the pair.
func OfLoop[S any, K comparable, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) (map[K]V, error) {
	return OfLoopResolv(source, hasNext, getNext, kvit.FirstVal[K, V])
}

// OfLoopResolv builds a map by iterating key\value pairs of a source.
// The hasNext specifies a predicate that tests existing of a next pair in the source.
// The getNext extracts the pair.
// The resolv selects value for duplicated keys.
func OfLoopResolv[S any, K comparable, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error), resolv func(K, V, V) V) (map[K]V, error) {
	r := map[K]V{}
	for hasNext(source) {
		if k, v, err := getNext(source); err != nil {
			return r, err
		} else if ov, ok := r[k]; ok {
			r[k] = resolv(k, ov, v)
		} else {
			r[k] = v
		}
	}
	return r, nil
}

// Generate builds a map by an generator function.
// The next returns an key\value pair, or false if the generation is over, or an error.
func Generate[K comparable, V any](next func() (K, V, bool, error)) (map[K]V, error) {
	return GenerateResolv(next, kvit.FirstVal[K, V])
}

// GenerateResolv builds a map by an generator function.
// The next returns an key\value pair, or false if the generation is over, or an error.
// The resolv selects value for duplicated keys.
func GenerateResolv[K comparable, V any](next func() (K, V, bool, error), resolv func(K, V, V) V) (map[K]V, error) {
	r := map[K]V{}
	for {
		k, v, ok, err := next()
		if err != nil || !ok {
			return r, err
		} else if ov, ok := r[k]; ok {
			r[k] = resolv(k, ov, v)
		} else {
			r[k] = v
		}
	}
}

// Copy makes a map copy.
func Copy[M ~map[K]V, K comparable, V any](elements M) M {
	copied := make(M, len(elements))
	for key, val := range elements {
		copied[key] = val
	}
	return copied
}

// Keys makes a slice of map keys.
func Keys[M ~map[K]V, K comparable, V any](elements M) []K {
	keys := make([]K, 0, len(elements))
	for key := range elements {
		keys = append(keys, key)
	}
	return keys
}

// Track applies a tracker for every key/value pairs from a map. To stop traking just return the ErrBreak.
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

// TrackEach applies a tracker for every key/value pairs from a map.
func TrackEach[M ~map[K]V, K comparable, V any](elements M, tracker func(K, V)) {
	for key, val := range elements {
		tracker(key, val)
	}
}

// For applies a walker for every key/value pairs from a map. Key/value pair is boxed to the KV. To stop walking just return the ErrBreak.
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

// ForEach applies a walker for every key/value pairs from a map. Key/value pair is boxed to the KV.
func ForEach[M ~map[K]V, K comparable, V any](elements M, walker func(c.KV[K, V])) {
	for key, val := range elements {
		walker(c.NewKV(key, val))
	}
}

// TrackOrdered applies a tracker for every key/value pairs from a map in order. To stop traking just return the ErrBreak.
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

// TrackEachOrdered applies a tracker for every key/value pairs from a map in order.
func TrackEachOrdered[M ~map[K]V, K comparable, V any](elements []K, uniques M, tracker func(K, V)) {
	for _, key := range elements {
		tracker(key, uniques[key])
	}
}

// ForOrdered applies a walker for every key/value pairs from a map in order. Key/value pair is boxed to the KV. To stop walking just return the ErrBreak.
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

// ForEachOrdered applies a walker for every key/value pairs from a map in order. Key/value pair is boxed to the KV.
func ForEachOrdered[M ~map[K]V, K comparable, V any](elements []K, uniques M, walker func(c.KV[K, V])) {
	for _, key := range elements {
		walker(c.NewKV(key, uniques[key]))
	}
}

// ForKeys applies a walker for every key from a map. To stop walking just return the ErrBreak.
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

// ForEachKey applies a walker for every key from a map.
func ForEachKey[M ~map[K]V, K comparable, V any](elements M, walker func(K)) {
	for key := range elements {
		walker(key)
	}
}

// ForValues applies a walker for every value from a map. To stop walking just return the ErrBreak.
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

// ForEachValue applies a walker for every value from a map.
func ForEachValue[M ~map[K]V, K comparable, V any](elements M, walker func(V)) {
	for _, val := range elements {
		walker(val)
	}
}

// ForOrderedValues applies a walker for every value from a map in order. To stop walking just return the ErrBreak.
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

// ForEachOrderedValues applies a walker for every value from a map in order.
func ForEachOrderedValues[M ~map[K]V, K comparable, V any](elements []K, uniques M, walker func(V)) {
	for _, key := range elements {
		val := uniques[key]
		walker(val)
	}
}

// ToStringOrdered converts elements to the string representation according to the order.
func ToStringOrdered[M ~map[K]V, K comparable, V any](order []K, elements M) string {
	return ToStringOrderedf(order, elements, "%+v:%+v", " ")
}

// ToStringOrderedf converts elements to a string representation using a key/value pair format and a delimeter. In order.
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

// ToString converts elements to the string representation.
func ToString[M ~map[K]V, K comparable, V any](elements M) string {
	return ToStringf(elements, "%+V:%+V", " ")
}

// ToStringf converts elements to a string representation using a key/value pair format and a delimeter.
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
