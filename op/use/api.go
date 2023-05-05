// Package use provides conditional expression builders
package use

import "github.com/m4gshm/gollections/op"

// If builds useIf(tr, condition).Else(fals) expression builder
func If[T any](tru T, condition bool) UseIf[T] {
	return UseIf[T]{tru, condition}
}

// UseIf is if...else expression builder
type UseIf[T any] struct {
	Tru       T
	Condition bool
}

// Else returns the tru or the fals according to the condition
func (u UseIf[T]) Else(fals T) T {
	return op.IfElse(u.Condition, u.Tru, fals)
}
