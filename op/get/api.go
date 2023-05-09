// Package get provides conditional expression builders
package get

// If builds get.If(condition, tru).Else(fals) expression builder
func If[T any](condition bool, tru func() T) When[T] {
	return When[T]{condition, tru}
}

// One builds get.One(tru).If(condition).Else(fals) expression builder
func One[T any](one func() T) ThisOne[T] {
	return This(one)
}

// This builds get.This(tru).If(condition).Else(fals) expression builder
func This[T any](one func() T) ThisOne[T] {
	return ThisOne[T]{one}
}

// When is if...else expression builder
type When[T any] struct {
	Condition bool
	Then      func() T
}

// Else returns the tru or the fals according to the condition
func (u When[T]) Else(fals T) T {
	if u.Condition {
		return u.Then()
	}
	return fals
}

// ElseGet returns the tru or executes the fals according to the condition
func (u When[T]) ElseGet(fals func() T) T {
	if u.Condition {
		return u.Then()
	}
	return fals()
}

// ThisOne is if...else expression builder
type ThisOne[T any] struct {
	Value func() T
}

// If is condition part of get.One(condition).If(tru).Else(fals) expression builder
func (u ThisOne[T]) If(condition bool) When[T] {
	return If(condition, u.Value)
}
