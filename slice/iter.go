package slice

// IsValidIndex checks if index is out of range
func IsValidIndex(size, index int) bool {
	return index > -1 && index < size
}

// Get safely returns an element of the 'elements' slice by the 'current' index or return zero value of T if the index is more than size-1 or less 0
func Get[TS ~[]T, T any](elements TS, current int) T {
	v, _ := Gett(elements, current)
	return v
}

// Gett safely returns an element of the 'elements' slice by the 'current' index or return zero value of T if the index is more than size-1 or less 0
// ok == true if success
func Gett[TS ~[]T, T any](elements TS, current int) (element T, ok bool) {
	if current >= 0 && current < len(elements) {
		element, ok = (elements)[current], true
	}
	return element, ok
}
