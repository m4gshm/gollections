// Package use provides conditional expression builders
package use

// If builds use.If(condition, tru).Else(fals) expression builder
func If[T any](condition bool, tru T) When[T] {
	return When[T]{condition, tru}
}
