package convert

func Key[V, K, KOUT any](by func(K) KOUT) (out func(key K, val V) (KOUT, V)) {
	if by != nil {
		out = func(key K, val V) (KOUT, V) { return by(key), val }
	}
	return
}

func Value[K, V, VOUT any](by func(V) VOUT) (out func(key K, val V) (K, VOUT)) {
	if by != nil {
		out = func(key K, val V) (K, VOUT) { return key, by(val) }
	}
	return
}
