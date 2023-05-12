// Package wrap provides string wrap utils
package wrap

import "github.com/m4gshm/gollections/op/string_"

// NonEmpty - wrap.NonEmpty wraps the target string only if it is non-empty
func NonEmpty(pref, target, post string) string {
	return string_.WrapNonEmpty(pref, target, post)
}
