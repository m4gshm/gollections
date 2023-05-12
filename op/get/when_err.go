package get

type WhenErr[T any] struct {
	Condition bool
	Then      func() (T, error)
}

// Else returns the tru or the fals according to the condition
func (w WhenErr[T]) Else(fals T) (T, error) {
	if w.Condition {
		return w.Then()
	}
	return fals, nil
}

func (w WhenErr[T]) ElseErr(err error) (T, error) {
	if w.Condition {
		return w.Then()
	}
	var fals T
	return fals, err
}

func (w WhenErr[T]) ElseGetErr(fals func() (T, error)) (T, error) {
	if w.Condition {
		return w.Then()
	}
	return fals()
}

// ElseGet returns the tru or executes the fals according to the condition
func (w WhenErr[T]) ElseGet(fals func() T) (T, error) {
	if w.Condition {
		return w.Then()
	}
	return fals(), nil
}

func (w WhenErr[T]) If(condition bool, tru T) WhenErr[T] {
	if w.Condition {
		return w
	}
	return If_(condition, func() (T, error) { return tru, nil })
}

func (w WhenErr[T]) IfGet(condition bool, tru func() T) WhenErr[T] {
	if w.Condition {
		return w
	}
	return If_(condition, func() (T, error) { return tru(), nil })
}

func (w WhenErr[T]) IfGetErr(condition bool, tru func() (T, error)) WhenErr[T] {
	if w.Condition {
		return w
	}
	return If_(condition, tru)
}

func (w WhenErr[T]) Other(condition func() bool, tru func() T) WhenErr[T] {
	if w.Condition {
		return w
	}
	return If_(condition(), func() (T, error) { return tru(), nil })
}

func (w WhenErr[T]) OtherErr(condition func() bool, tru func() (T, error)) WhenErr[T] {
	if w.Condition {
		return w
	}
	return If_(condition(), tru)
}
