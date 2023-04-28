// Package always provides constant predicate implementations
package always

// True always return true
func True[T any](_ T) (bool, error) {
	return true, nil
}

// False always return false
func False[T any](_ T) (bool, error) {
	return false, nil
}
