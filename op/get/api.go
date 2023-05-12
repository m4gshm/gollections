// Package get provides conditional expression builders
package get

// If builds get.If(condition, tru).Else(fals) expression builder
func If[T any](condition bool, tru func() T) When[T] {
	return When[T]{condition, tru}
}

func If_[T any](condition bool, tru func() (T, error)) WhenErr[T] {
	return WhenErr[T]{condition, tru}
}
