// Package string_ provides string utils
package string_

// WrapNoEmpty returns wrapped the target string if it is non-empty
func WrapNoEmpty(pref, target, post string) string {
	if len(target) == 0 {
		return ""
	}
	return pref + target + post
}
