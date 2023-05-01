package mutable

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/collection/immutable/vector"
)

var (
	_ immutable.Vector[int]  = vector.Of(1, 2, 3)
	_ collection.Vector[int] = immutable.NewVector([]int{1, 2, 3})
)
