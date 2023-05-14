package use

// WhenErr if..else expression builder
type WhenErr[T any] struct {
	Condition bool
	Then      T
	Err       error
}

// Else returns result according to the condition
func (w WhenErr[T]) Else(fals T) (T, error) {
	if w.Condition {
		return w.Then, w.Err
	}
	return fals, nil
}

// ElseErr returns the success result or error according to the condition
func (w WhenErr[T]) ElseErr(err error) (T, error) {
	if w.Condition {
		return w.Then, w.Err
	}
	var fals T
	return fals, err
}

// ElseGetErr returns the success result or a result of the fals function according to the condition
func (w WhenErr[T]) ElseGetErr(fals func() (T, error)) (T, error) {
	if w.Condition {
		return w.Then, w.Err
	}
	return fals()
}

// ElseGet returns the tru or a return of the fals function according to the condition
func (w WhenErr[T]) ElseGet(fals func() T) (T, error) {
	if w.Condition {
		return w.Then, w.Err
	}
	return fals(), nil
}

// If creates new condition branch in the expression
func (w WhenErr[T]) If(condition bool, tru T) WhenErr[T] {
	if w.Condition {
		return w
	}
	return ifErrEvaluated(condition, tru, nil)
}

// IfGet creates new condition branch for a getter function
func (w WhenErr[T]) IfGet(condition bool, tru func() T) WhenErr[T] {
	if w.Condition {
		return w
	}

	var other T
	if condition {
		other = tru()
	}
	return ifErrEvaluated(condition, other, nil)
}

// IfGetErr creates new condition branch for an error return getter function
func (w WhenErr[T]) IfGetErr(condition bool, tru func() (T, error)) WhenErr[T] {
	if w.Condition {
		return w
	}
	return If_(condition, tru)
}

// Other creates new condition branch for a getter function.
// The condition function is called only if the current condition is false.
func (w WhenErr[T]) Other(condition func() bool, tru func() T) WhenErr[T] {
	if w.Condition {
		return w
	}

	var (
		otherCondition = condition()
		other          T
	)
	if otherCondition {
		other = tru()
	}
	return ifErrEvaluated(otherCondition, other, nil)
}

// OtherErr creates new condition branch for an error return getter function.
// The condition function is called only if the current condition is false.
func (w WhenErr[T]) OtherErr(condition func() bool, tru func() (T, error)) WhenErr[T] {
	if w.Condition {
		return w
	}
	return If_(condition(), tru)
}

func ifErrEvaluated[T any](condition bool, tru T, err error) WhenErr[T] {
	var (
		then    T
		thenErr error
	)
	if condition {
		then, thenErr = tru, err

	}
	return WhenErr[T]{condition, then, thenErr}
}
