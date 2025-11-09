// Package convert provides converting helpers
package convert

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/k"
)

// AsIs helper for Map, Flatt.
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

// ToSlice convert an one element to a slice.
func ToSlice[T any](value T) []T { return []T{value} }

// AsSlice convert an one element to a slice.
func AsSlice[T any](value T) []T { return ToSlice(value) }

// KeyValue transforms one element to one key/value pair.
func KeyValue[T, K, V any](element T, keyExtractor func(T) K, valExtractor func(T) V) c.KV[K, V] {
	return k.V(keyExtractor(element), valExtractor(element))
}

// KeysValues transforms one element to multiple key/value pairs slices.
func KeysValues[T, K, V any](element T, keysExtractor func(T) []K, valsExtractor func(T) []V) (out []c.KV[K, V]) {
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

// ExtraVals transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func ExtraVals[T, V any](element T, valsExtractor func(T) []V) (out []c.KV[T, V]) {
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

// ExtraKeys transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func ExtraKeys[T, K any](element T, keysExtractor func(T) []K) (out []c.KV[K, T]) {
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

// ToPtr converts a value to the value pointer
func ToPtr[T any](value T) *T {
	return &value
}

// ToVal returns a value referenced by the pointer or the zero value if the pointer is nil
func ToVal[T any](pointer *T) (t T) {
	if pointer != nil {
		t = *pointer
	}
	return t
}

// ToValNotNil returns a value referenced by the pointer or ok==false if the pointer is nil
func ToValNotNil[T any](pointer *T) (t T, ok bool) {
	if pointer != nil {
		return *pointer, true
	}
	return t, false
}

// ToType converts I to T. If successful, returns the converted value and true.
func ToType[T, I any](i I) (T, bool) {
	var a any = i
	t, ok := a.(T)
	return t, ok
}

// NilSafe filters not nil elements, converts that ones, filters not nils after converting and returns them
func NilSafe[From, To any](converter func(*From) *To) func(f *From) (*To, bool) {
	return func(f *From) (*To, bool) {
		if f != nil {
			if t := converter(f); t != nil {
				return t, true
			}
		}
		return nil, false
	}
}
