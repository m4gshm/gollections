package use

// When is if...else expression builder
type When[T any] struct {
	condition bool
	then      T
}

// Else returns result according to the condition
func (w When[T]) Else(fals T) T {
	if w.condition {
		return w.then
	}
	return fals
}

// ElseErr returns the success result or error according to the condition
func (w When[T]) ElseErr(err error) (T, error) {
	if w.condition {
		return w.then, nil
	}
	var fals T
	return fals, err
}

// ElseGetErr returns the success result or a result of the fals function according to the condition
func (w When[T]) ElseGetErr(fals func() (T, error)) (T, error) {
	if w.condition {
		return w.then, nil
	}
	return fals()
}

// ElseGet returns the tru or a return of the fals function according to the condition
func (w When[T]) ElseGet(fals func() T) T {
	if w.condition {
		return w.then
	}
	return fals()
}

// If creates new condition branch in the expression.
// Looks like value := use.If(condition1, variant1).If(condition2, variant2).Else(defaultVariant) .
func (w When[T]) If(condition bool, tru T) When[T] {
	if w.condition {
		return w
	}
	return If(condition, tru)
}

// IfGet creates new condition branch for a getter function.
// Looks like value := use.If(condition1, variant1).IfGet(condition2, getterFunction1).Else(defaultVariant) .
func (w When[T]) IfGet(condition bool, tru func() T) When[T] {
	if w.condition {
		return w
	}
	var otherTru T
	if condition {
		otherTru = tru()
	}
	return If(condition, otherTru)
}

// IfGetErr creates new condition branch for an error return getter function
func (w When[T]) IfGetErr(condition bool, tru func() (T, error)) WhenErr[T] {
	if w.condition {
		return ifErrEvaluated(w.condition, w.then, nil)
	}
	var (
		otherTru T
		otherErr error
	)
	if condition {
		otherTru, otherErr = tru()
	}
	return ifErrEvaluated(condition, otherTru, otherErr)
}

// Other creates new condition branch for a getter function.
// The condition function is called only if the current condition is false.
func (w When[T]) Other(condition func() bool, tru func() T) When[T] {
	if w.condition {
		return w
	}
	var (
		otherCondition = condition()
		otherTru       T
	)
	if otherCondition {
		otherTru = tru()
	}
	return If(otherCondition, otherTru)
}

// OtherErr creates new condition branch for an error return getter function.
// The condition function is called only if the current condition is false.
func (w When[T]) OtherErr(condition func() bool, tru func() (T, error)) WhenErr[T] {
	if w.condition {
		return ifErrEvaluated(w.condition, w.then, nil)
	}
	var (
		otherCondition = condition()
		otherTru       T
		otherErr       error
	)
	if otherCondition {
		otherTru, otherErr = tru()
	}
	return ifErrEvaluated(otherCondition, otherTru, otherErr)
}

func (w When[T]) EvalIf(condition func() bool, tru T) When[T] {
	if w.condition {
		return w
	}
	return If(condition(), tru)
}
