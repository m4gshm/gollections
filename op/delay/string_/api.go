package string_

import (
	"strings"
)

// Of returns string builder
func Of(parts ...string) func() string {
	return func() string { return strings.Join(parts, "") }
}

// Wrap returns wrapped string builder
func Wrap(pref, post string) func(s string) string {
	return func(s string) string { return pref + s + post }
}
