// Package not provides negalive predicates like 'not equals to'
package not

import "github.com/m4gshm/gollections/op/check"

// Nil - not.Nil checks whether the reference is not nil.
func Nil[T any](reference *T) bool {
	return check.NotNil(reference)
}

// Empty - not.Empty checks that the slice contains elements.
func Empty[TS ~[]T, T any](elements TS) bool {
	return check.NotEmpty(elements)
}

// EmptyStr - not.EmptyStr checks whether the specified string is not empty.
func EmptyStr(s string) bool {
	return check.NotEmptyStr(s)
}
