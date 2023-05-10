package replace

func By[T any](r T) func(t T) T {
	return func(t T) T { return r }
}
