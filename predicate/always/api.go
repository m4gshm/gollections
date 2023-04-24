package always


func True[T any](t T) bool {
	return true
}

func False[T any](t T) bool {
	return false
}
