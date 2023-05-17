// Package join provides string builders
package join

import "github.com/m4gshm/gollections/op/delay/string_"

// NonEmpty returns concatenated string builder
func NonEmpty(joiner string) func(first, second string) string {
	return string_.JoinNonEmpty(joiner)
}
