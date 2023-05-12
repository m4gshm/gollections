package get

// When is if...else expression builder
type When[T any] struct {
	Condition bool
	Then      func() T
}

// Else returns the tru or the fals according to the condition
func (w When[T]) Else(fals T) T {
	if w.Condition {
		return w.Then()
	}
	return fals
}

func (w When[T]) ElseErr(err error) (T, error) {
	if w.Condition {
		return w.Then(), nil
	}
	var fals T
	return fals, err
}

func (w When[T]) ElseGetErr(fals func() (T, error)) (T, error) {
	if w.Condition {
		return w.Then(), nil
	}
	return fals()
}

// ElseGet returns the tru or executes the fals according to the condition
func (w When[T]) ElseGet(fals func() T) T {
	if w.Condition {
		return w.Then()
	}
	return fals()
}

func (w When[T]) If(condition bool, tru T) When[T] {
	if w.Condition {
		return w
	}
	return If(condition, func() T { return tru })
}

func (w When[T]) IfGet(condition bool, tru func() T) When[T] {
	if w.Condition {
		return w
	}
	return If(condition, tru)
}

func (w When[T]) IfGetErr(condition bool, tru func() (T, error)) WhenErr[T] {
	if w.Condition {
		return If_(true, func() (T, error) { return w.Then(), nil })
	}
	return If_(condition, tru)
}

func (w When[T]) Other(condition func() bool, tru func() T) When[T] {
	if w.Condition {
		return w
	}
	return If(condition(), tru)
}

func (w When[T]) OtherErr(condition func() bool, tru func() (T, error)) WhenErr[T] {
	if w.Condition {
		return If_(true, func() (T, error) { return w.Then(), nil })
	}
	return If_(condition(), tru)
}
