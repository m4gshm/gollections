// Package wrap provides wrap string builders
package wrap

import "github.com/m4gshm/gollections/op/delay/string_"

// By returns wrapped string builder
func By(pref, post string) func(s string) string {
	return string_.Wrap(pref, post)
}

// NonEmpty returns wrapped string builder
func NonEmpty(pref, post string) func(target string) string {
	return string_.WrapNonEmpty(pref, post)
}
