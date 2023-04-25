// Package ptr provides value, pointer convert helpers
package ptr

// Of is value-to-pointer conversion helper
func Of[T any](t T) *T {
	return &t
}
