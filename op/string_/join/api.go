// Package join provides concatenate string utils
package join

import "github.com/m4gshm/gollections/op/string_"

// NoEmpty returns concatenated string
func NoEmpty(first, joiner, second string) string {
	return string_.JoinNonEmpty(first, joiner, second)
}
