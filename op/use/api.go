// Package use provides conditional expression builders
package use

// If builds use.If(condition, tru).Else(fals) expression builder
func If[T any](condition bool, tru T) When[T] {
	return When[T]{condition, tru}
}

// If builds use.One(tru).If(condition).Else(fals) expression builder
func One[T any](one T) ThisOne[T] {
	return This(one)
}

// If builds use.This(tru).If(condition).Else(fals) expression builder
func This[T any](one T) ThisOne[T] {
	return ThisOne[T]{one}
}

// When is if...else expression builder
type When[T any] struct {
	Condition bool
	Then      T
}

// Else returns the tru or the fals according to the condition
func (u When[T]) Else(fals T) T {
	if u.Condition {
		return u.Then
	}
	return fals
}

// ElseGet returns the tru or executes the fals according to the condition
func (u When[T]) ElseGet(fals func() T) T {
	if u.Condition {
		return u.Then
	}
	return fals()
}

// ThisOne is if...else expression builder
type ThisOne[T any] struct {
	Value T
}

func (u ThisOne[T]) If(condition bool) When[T] {
	return If(condition, u.Value)
}
