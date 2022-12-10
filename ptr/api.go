package ptr

// Of is value-to-pointer conversion helper
func Of[T any](t T) *T {
	return &t
}
