package use

// When is if...else expression builder
type When[T any] struct {
	condition bool
	then      T
}

// Else evaluates expression and returns result
func (w When[T]) Else(fals T) T {
	if w.condition {
		return w.then
	}
	return fals
}

// ElseZero evaluates expression and returns zero value as default
func (w When[T]) ElseZero() (out T) {
	if w.condition {
		return w.then
	}
	return out
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

// If creates new condition branch in the expression
func (w When[T]) If(condition bool, tru T) When[T] {
	if w.condition {
		return w
	}
	return newWhen(condition, tru)
}

// IfGet creates new condition branch for a getter function
func (w When[T]) IfGet(condition bool, tru func() T) When[T] {
	if w.condition {
		return w
	}
	return newWhen(condition, evaluate(condition, tru))
}

// IfGetErr creates new condition branch for an error return getter function
func (w When[T]) IfGetErr(condition bool, tru func() (T, error)) WhenErr[T] {
	if w.condition {
		return newWhenErr(w.condition, w.then, nil)
	}
	return If_(condition, tru)
}

// Other creates new condition branch for a getter function.
// The condition function is called only if the current condition is false.
func (w When[T]) Other(condition func() bool, tru func() T) When[T] {
	if w.condition {
		return w
	}

	c := condition()
	return newWhen(c, evaluate(c, tru))
}

// OtherErr creates new condition branch for an error return getter function.
// The condition function is called only if the current condition is false.
func (w When[T]) OtherErr(condition func() bool, tru func() (T, error)) WhenErr[T] {
	if w.condition {
		return newWhenErr(w.condition, w.then, nil)
	}
	return If_(condition(), tru)
}

// Eval evaluates the expression and returns ok==false if there is no satisfied condition
func (w When[T]) Eval() (out T, ok bool) {
	if w.condition {
		return w.then, true
	}
	return out, false
}

func newWhen[T any](condition bool, then T) When[T] {
	return When[T]{condition, then}
}

func evaluate[T any](condition bool, tru func() T) (out T) {
	if condition {
		out = tru()
	}
	return out
}
