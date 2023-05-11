package wrap

import "github.com/m4gshm/gollections/op/string_"

func NoEmpty(pref, target, post string) string {
	return string_.WrapNoEmpty(pref, target, post)
}
