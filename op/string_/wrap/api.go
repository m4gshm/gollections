// Package wrap provides string wrap utils
package wrap

import "github.com/m4gshm/gollections/op/string_"

// NoEmpty - wrap.NoEmpty wraps the target string only if it is non-empty
func NoEmpty(pref, target, post string) string {
	return string_.WrapNoEmpty(pref, target, post)
}
