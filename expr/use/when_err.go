package use

// WhenErr if..else expression builder
type WhenErr[T any] struct {
	condition bool
	then      T
	err       error
}

// Else evaluates expression and returns result
func (w WhenErr[T]) Else(fals T) (T, error) {
	if w.condition {
		return w.then, w.err
	}
	return fals, nil
}

// ElseZero evaluates expression and returns zero value as default
func (w WhenErr[T]) ElseZero() (out T, err error) {
	if w.condition {
		return w.then, w.err
	}
	return out, nil
}

// ElseErr returns the success result or error according to the condition
func (w WhenErr[T]) ElseErr(err error) (T, error) {
	if w.condition {
		return w.then, w.err
	}
	var fals T
	return fals, err
}

// ElseGetErr returns the success result or a result of the fals function according to the condition
func (w WhenErr[T]) ElseGetErr(fals func() (T, error)) (T, error) {
	if w.condition {
		return w.then, w.err
	}
	return fals()
}

// ElseGet returns the tru or a return of the fals function according to the condition
func (w WhenErr[T]) ElseGet(fals func() T) (T, error) {
	if w.condition {
		return w.then, w.err
	}
	return fals(), nil
}

// If creates new condition branch in the expression
func (w WhenErr[T]) If(condition bool, tru T) WhenErr[T] {
	if w.condition {
		return w
	}
	return newWhenErr(condition, tru, nil)
}

// IfGet creates new condition branch for a getter function
func (w WhenErr[T]) IfGet(condition bool, tru func() T) WhenErr[T] {
	if w.condition {
		return w
	}
	return newWhenErr(condition, evaluate(condition, tru), nil)
}

// IfGetErr creates new condition branch for an error return getter function
func (w WhenErr[T]) IfGetErr(condition bool, tru func() (T, error)) WhenErr[T] {
	if w.condition {
		return w
	}
	return If_(condition, tru)
}

// Other creates new condition branch for a getter function.
// The condition function is called only if the current condition is false.
func (w WhenErr[T]) Other(condition func() bool, tru func() T) WhenErr[T] {
	if w.condition {
		return w
	}

	c := condition()
	return newWhenErr(c, evaluate(c, tru), nil)
}

// OtherErr creates new condition branch for an error return getter function.
// The condition function is called only if the current condition is false.
func (w WhenErr[T]) OtherErr(condition func() bool, tru func() (T, error)) WhenErr[T] {
	if w.condition {
		return w
	}
	return If_(condition(), tru)
}

// Eval evaluates the expression and returns ok==false if there is no satisfied condition
func (w WhenErr[T]) Eval() (out T, ok bool, err error) {
	if w.condition {
		return w.then, true, w.err
	}
	return out, false, nil
}

func newWhenErr[T any](condition bool, then T, err error) WhenErr[T] {
	return WhenErr[T]{condition, then, err}
}

func evaluateErr[T any](condition bool, tru func() (T, error)) (out T, err error) {
	if condition {
		out, err = tru()
	}
	return out, err
}
