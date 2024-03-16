package over

func Filtered[TS ~[]T, T any](elements TS, filter func(T) bool) func(func(int, T) bool) {
	return func(consumer func(int, T) bool) {
		for i, e := range elements {
			if filter(e) {
				if !consumer(i, e) {
					return
				}
			}
		}
	}
}

func Converted[FS ~[]From, From, To any](elements FS, converter func(From) To) func(func(int, To) bool) {
	return func(consumer func(int, To) bool) {
		for i, e := range elements {
				if !consumer(i, converter(e)) {
					return
				}
		}
	}
}
