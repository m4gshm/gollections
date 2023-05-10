package chain

func Of[T any](funcs ...func(T) T) func(T) T {
	return func(t T) T {
		for _, f := range funcs {
			t = f(t)
		}
		return t
	}
}