// Package always provides constant predicate implementations
package always

// True always return true
func True[T any](_ T) bool {
	return true
}

// False always return false
func False[T any](_ T) bool {
	return false
}
