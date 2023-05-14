package use

import "github.com/m4gshm/gollections/expr/get"

// When is if...else expression builder
type When[T any] struct {
	Condition bool
	Then      T
}

// Else returns result according to the condition
func (w When[T]) Else(fals T) T {
	if w.Condition {
		return w.Then
	}
	return fals
}

// ElseErr returns the success result or error according to the condition
func (w When[T]) ElseErr(err error) (T, error) {
	if w.Condition {
		return w.Then, nil
	}
	var fals T
	return fals, err
}

// ElseGetErr returns the success result or a result of the fals function according to the condition
func (w When[T]) ElseGetErr(fals func() (T, error)) (T, error) {
	if w.Condition {
		return w.Then, nil
	}
	return fals()
}

// ElseGet returns the tru or a return of the fals function according to the condition
func (w When[T]) ElseGet(fals func() T) T {
	if w.Condition {
		return w.Then
	}
	return fals()
}

// If creates new condition branch in the expression.
// Looks like value := use.If(condition1, variant1).If(condition2, variant2).Else(defaultVariant) .
func (w When[T]) If(condition bool, tru T) When[T] {
	if w.Condition {
		return w
	}
	return If(condition, tru)
}

// IfGet creates new condition branch for a getter function.
// Looks like value := use.If(condition1, variant1).IfGet(condition2, getterFunction1).Else(defaultVariant) .
func (w When[T]) IfGet(condition bool, tru func() T) get.When[T] {
	if w.Condition {
		return get.If(true, func() T { return w.Then })
	}
	return get.If(condition, tru)
}

// IfGetErr creates new condition branch for an error return getter function
func (w When[T]) IfGetErr(condition bool, tru func() (T, error)) get.WhenErr[T] {
	if w.Condition {
		return get.If_(true, func() (T, error) { return w.Then, nil })
	}
	return get.If_(condition, tru)
}

// Other creates new condition branch for a getter function.
// Other creates new condition branch for a getter function.
func (w When[T]) Other(condition func() bool, tru func() T) get.When[T] {
	if w.Condition {
		return get.If(true, func() T { return w.Then })
	}
	return get.If(condition(), tru)
}

// OtherErr creates new condition branch for an error return getter function.
// The condition function is called only if the current condition is false.
func (w When[T]) OtherErr(condition func() bool, tru func() (T, error)) get.WhenErr[T] {
	if w.Condition {
		return get.If_(true, func() (T, error) { return w.Then, nil })
	}
	return get.If_(condition(), tru)
}
