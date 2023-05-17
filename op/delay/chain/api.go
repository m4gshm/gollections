package chain

// Of returns a function that calls the specified functions in sequence with the argument returned by the previous one
func Of[T any](funcs ...func(T) T) func(T) T {
	return func(t T) T {
		for _, f := range funcs {
			t = f(t)
		}
		return t
	}
}
