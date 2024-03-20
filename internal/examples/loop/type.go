package loop

type (
	Loop[T any]           func() (element T, ok bool)
	KVLoop[K, V any]      func() (key K, value V, ok bool)
	BreakLoop[T any]      func() (element T, ok bool, err error)
	BreakKVLoop[K, V any] func() (key K, value V, ok bool, err error)
)
