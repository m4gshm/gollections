// Package convert provides converting helpers
package convert

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/k"
)

// AsIs helper for Map, Flatt
func AsIs[T any](value T) T { return value }

// And apply two converters in order.
func And[I, O, N any](first func(I) O, second func(O) N) func(I) N {
	return func(i I) N { return second(first(i)) }
}

// Or applies first Converter, applies second Converter if the first returns zero.
func Or[I, O comparable](first func(I) O, second func(I) O) func(I) O {
	return func(i I) O {
		c := first(i)
		var no O
		if no == c {
			return second(i)
		}
		return c
	}
}

// ToSlice convert an one element to a slice
func ToSlice[T any](value T) []T { return []T{value} }

// AsSlice convert an one element to a slice
func AsSlice[T any](value T) []T { return ToSlice(value) }

// ToKV transforms one element to one key/value pair
func ToKV[T, K, V any](element T, keyExtractor func(T) K, valExtractor func(T) V) c.KV[K, V] {
	return k.V(keyExtractor(element), valExtractor(element))
}

// ToKVs transforms one element to multiple key/value pairs slices
func ToKVs[T, K, V any](element T, keysExtractor func(T) []K, valsExtractor func(T) []V) (out []c.KV[K, V]) {
	var (
		keys   = keysExtractor(element)
		values = valsExtractor(element)
		lv     = len(values)
	)
	if len(keys) == 0 {
		var key K
		for _, value := range values {
			out = append(out, k.V(key, value))
		}
	} else {
		for _, key := range keys {
			if lv == 0 {
				var value V
				out = append(out, k.V(key, value))
			} else {
				for _, value := range values {
					out = append(out, k.V(key, value))
				}
			}
		}
	}
	return out
}

// FlatValues transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlatValues[T, V any](element T, valsExtractor func(T) []V) (out []c.KV[T, V]) {
	var (
		key    = element
		values = valsExtractor(element)
		lv     = len(values)
	)
	if lv == 0 {
		var value V
		out = append(out, k.V(key, value))
	} else {
		for _, value := range values {
			out = append(out, k.V(key, value))
		}
	}
	return out
}

// FlatKeys transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlatKeys[T, K any](element T, keysExtractor func(T) []K) (out []c.KV[K, T]) {
	var (
		keys  = keysExtractor(element)
		value = element
	)
	if len(keys) == 0 {
		var key K
		out = append(out, k.V(key, value))
	} else {
		for _, key := range keys {
			out = append(out, k.V(key, value))
		}
	}
	return out
}

// ToPointer converts a value to the value pointer
func ToPointer[T any](value T) *T {
	return &value
}

// ToValue returns a value referenced by the pointer or the zero value if the pointer is nil
func ToValue[T any](pointer *T) (t T) {
	if pointer != nil {
		t = *pointer
	}
	return t
}

// GetValue returns a value referenced by the pointer or ok==false if the pointer is nil
func GetValue[T any](pointer *T) (t T, ok bool) {
	if pointer != nil {
		return *pointer, true
	}
	return t, false
}
