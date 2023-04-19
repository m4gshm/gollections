package from

func To[F, I, T any](f func(F) I, t func(I) T) func(F) T {
	return func(v F) T { return t(f(v)) }
}
