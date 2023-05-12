// Package string_ provides string builders
package string_

import (
	"strings"

	"github.com/m4gshm/gollections/op/string_"
)

// Of returns string builder
func Of(parts ...string) func() string {
	return func() string { return strings.Join(parts, "") }
}

// Wrap returns wrapped string builder
func Wrap(pref, post string) func(target string) string {
	return func(s string) string { return pref + s + post }
}

// WrapNonEmpty returns wrapped string builder
func WrapNonEmpty(pref, post string) func(target string) string {
	return func(target string) string {
		return string_.WrapNonEmpty(pref, target, post)
	}
}

func JoinNonEmpty(joiner string) func(first, second string) string {
	return func(first, second string) string {
		return string_.JoinNonEmpty(first, joiner, second)
	}
}
