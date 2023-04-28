// Package as provides as.Is alias
package as

import (
	"github.com/m4gshm/gollections/convert"
	// "github.com/m4gshm/gollections/slice/iter"
)

// Is an alias of the conv.AsIs
func Is[T any](value T) T { return convert.AsIs(value) }

// Is an alias of the conv.AsSlice
func Slice[T any](value T) []T { return convert.AsSlice(value) }

// func KV[T any, K comparable, V any](element T, keyExtractor func(T) K, valExtractor func(T) V) *iter.KeyValuer[T, K, V] {
// 	kv := iter.NewKeyValuer([]T{element}, keyExtractor, valExtractor)
// 	return kv
// }

func ErrTail[I, O any](f func(I) O) func(I) (O, error) {
	return func(in I) (O, error) { return f(in), nil }
}
