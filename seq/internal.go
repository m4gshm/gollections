package seq

func empty[T any](_ func(T) bool)        {}
func empty2[K, V any](_ func(K, V) bool) {}
