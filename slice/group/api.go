package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
)

// Of - group.Of synonym of the slice.Group
func Of[T any, K comparable, TS ~[]T](elements TS, keyProducer c.Converter[T, K]) map[K]TS {
	return slice.Group(elements, keyProducer)
}

func InMultiple[T any, K comparable, TS ~[]T](elements TS, keysProducer c.Converter[T, []K]) map[K]TS {
	return slice.GroupInMultiple(elements, keysProducer)
}
