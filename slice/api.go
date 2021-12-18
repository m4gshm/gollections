package slice

func Of[T any](values ...T) []T {
	return values
}

func Convert[From, To any](values []From, by Converter[From, To]) ([]To, error) {
	out := make([]To, len(values))
	for i, v := range values {
		o, err := by(v)
		if err != nil {
			return nil, err
		}
		out[i] = o
	}
	return out, nil
}

type Converter[From, To any] func(v From) (To, error)

func And[I, O, N any](f Converter[I, O], s Converter[O, N]) Converter[I, N] {
	return func(i I) (N, error) {
		c, err := f(i)
		if err != nil {
			var n N
			return n, err
		}
		return s(c)
	}
}

func Or[I, O any](f Converter[I, O], s Converter[I, O]) Converter[I, O] {
	return func(i I) (O, error) {
		if c, err := f(i); err != nil {
			return s(i)
		} else {
			return c, nil
		}
	}
}
