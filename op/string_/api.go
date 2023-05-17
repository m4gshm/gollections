// Package string_ provides string utils
package string_

// WrapNonEmpty returns wrapped the target string if it is non-empty
func WrapNonEmpty(pref, target, post string) string {
	if len(target) == 0 {
		return ""
	}
	return pref + target + post
}

// JoinNonEmpty returns concatenated string
func JoinNonEmpty(first, joiner, second string) string {
	if len(first) == 0 || len(second) == 0 {
		return first + second
	}
	return first + joiner + second
}

// Empty checks whether the specified string is empty
func Empty(s string) bool {
	return len(s) == 0
}

// NotEmpty checks whether the specified string is not empty
func NotEmpty(s string) bool {
	return !Empty(s)
}
