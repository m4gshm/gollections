package mutable

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/vector"
)

var (
	_ *mutable.Vector[int]   = vector.Of(1, 2, 3)
	_ collection.Vector[int] = &mutable.Vector[int]{}
)
