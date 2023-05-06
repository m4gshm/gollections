// Package use provides conditional expression builders
package use

import "github.com/m4gshm/gollections/op"

// If builds useIf(tr, condition).Else(fals) expression builder
func If[T any](tru T, condition bool) If_[T] {
	return If_[T]{tru, condition}
}

// If builds use.One(tr).If(condition).Else(fals) expression builder
func One[T any](one T) One_[T] {
	return One_[T]{one}
}

// If_ is if...else expression builder
type If_[T any] struct {
	Then      T
	Condition bool
}

// OrElse returns the tru or the fals according to the condition
func (u If_[T]) OrElse(fals T) T {
	return op.IfElse(u.Condition, u.Then, fals)
}

// Else returns the tru or the fals according to the condition
func (u If_[T]) Else(fals T) T {
	return op.IfElse(u.Condition, u.Then, fals)
}

// ElseGet returns the tru or executes the fals according to the condition
func (u If_[T]) ElseGet(fals func() T) T {
	return op.IfDoElse(u.Condition, func() T { return u.Then }, fals)
}

// One_ is if...else expression builder
type One_[T any] struct {
	Value T
}

func (u One_[T]) If(condition bool) If_[T] {
	return If(u.Value, condition)
}
