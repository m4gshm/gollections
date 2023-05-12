// Package use provides conditional expression builders
package use

import "github.com/m4gshm/gollections/op/get"

// If builds use.If(condition, tru).Else(fals) expression builder
func If[T any](condition bool, tru T) When[T] {
	return When[T]{condition, tru}
}

// One builds use.One(tru).If(condition).Else(fals) expression builder
func One[T any](one T) ThisOne[T] {
	return This(one)
}

// This builds use.This(tru).If(condition).Else(fals) expression builder
func This[T any](one T) ThisOne[T] {
	return ThisOne[T]{one}
}

// When is if...else expression builder
type When[T any] struct {
	Condition bool
	Then      T
}

// Else returns the tru or the fals according to the condition
func (w When[T]) Else(fals T) T {
	if w.Condition {
		return w.Then
	}
	return fals
}

func (w When[T]) ElseErr(err error) (T, error) {
	if w.Condition {
		return w.Then, nil
	}
	var fals T
	return fals, err
}

func (w When[T]) ElseGetErr(err func() error) (T, error) {
	if w.Condition {
		return w.Then, nil
	}
	var fals T
	return fals, err()
}

func (w When[T]) ElseOptErr(fals func() (T, error)) (T, error) {
	if w.Condition {
		return w.Then, nil
	}
	return fals()
}

// ElseGet returns the tru or executes the fals according to the condition
func (w When[T]) ElseGet(fals func() T) T {
	if w.Condition {
		return w.Then
	}
	return fals()
}

// ThisOne is if...else expression builder
type ThisOne[T any] struct {
	Value T
}

// If is condition part of use.One(condition).If(tru).Else(fals) expression builder
func (t ThisOne[T]) If(condition bool) When[T] {
	return If(condition, t.Value)
}

func (w When[T]) ElseIf(condition bool, tru T) When[T] {
	if w.Condition {
		return w
	}
	return If(condition, tru)
}

func (w When[T]) ElseIfGet(condition bool, tru func() T) get.When[T] {
	if w.Condition {
		return get.If(true, func() T { return w.Then })
	}
	return get.If(condition, tru)
}

func (w When[T]) ElsIfGet(condition func() bool, tru func() T) get.When[T] {
	if w.Condition {
		return get.If(true, func() T { return w.Then })
	}
	return get.If(condition(), tru)
}