package strings

// Empty checks whether the specified string is empty
func Empty(s string) bool {
	return len(s) == 0
}

// NotEmpty checks whether the specified string is not empty
func NotEmpty(s string) bool {
	return !Empty(s)
}
