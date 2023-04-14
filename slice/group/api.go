package group

import (
	"github.com/m4gshm/gollections/slice"
)

// Of - group.Of synonym of the slice.Group
func Of[T any, K comparable, TS ~[]T](elements TS, keyProducer func(T) K) map[K]TS {
	return slice.Group(elements, keyProducer)
}

func InMultiple[T any, K comparable, TS ~[]T](elements TS, keysProducer func(T) []K) map[K]TS {
	return slice.GroupInMultiple(elements, keysProducer)
}
