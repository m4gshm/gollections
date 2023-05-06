package get

import "github.com/m4gshm/gollections/op"

// If builds get.If(tr, condition).Else(fals) expression builder
func If[T any](tru func() T, condition bool) If_[T] {
	return If_[T]{tru, condition}
}

func IfCalc[T any](tru func() T, condition func() bool) IfCalc_[T] {
	return IfCalc_[T]{tru, condition}
}

// If builds use.One(tr).If(condition).Else(fals) expression builder
func One[T any](one func() T) One_[T] {
	return One_[T]{one}
}

// If_ is if...else expression builder
type If_[T any] struct {
	Then      func() T
	Condition bool
}

// OrElse returns the tru or the fals according to the condition
func (u If_[T]) OrElse(fals func() T) T {
	return op.IfDoElse(u.Condition, u.Then, fals)
}

// Else returns the tru or the fals according to the condition
func (u If_[T]) Else(fals T) T {
	return op.IfDoElse(u.Condition, u.Then, func() T { return fals })
}

// ElseGet returns the tru or executes the fals according to the condition
func (u If_[T]) ElseGet(fals func() T) T {
	return op.IfDoElse(u.Condition, u.Then, fals)
}

// If_ is if...else expression builder
type IfCalc_[T any] struct {
	Then      func() T
	Condition func() bool
}

// OrElse returns the tru or the fals according to the condition
func (u IfCalc_[T]) OrElse(fals func() T) T {
	return op.IfDoElse(u.Condition(), u.Then, fals)
}

// Else returns the tru or the fals according to the condition
func (u IfCalc_[T]) Else(fals T) T {
	if u.Condition() {
		return u.Then()
	}
	return fals
}

// ElseGet returns the tru or executes the fals according to the condition
func (u IfCalc_[T]) ElseGet(fals func() T) T {
	return op.IfDoElse(u.Condition(), u.Then, fals)
}

// One_ is if...else expression builder
type One_[T any] struct {
	Getter func() T
}

func (u One_[T]) If(condition bool) If_[T] {
	return If(u.Getter, condition)
}

func (u One_[T]) If_(condition func() bool) IfCalc_[T] {
	return IfCalc(u.Getter, condition)
}
