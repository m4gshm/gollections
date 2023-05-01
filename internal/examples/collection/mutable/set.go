package mutable

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/collection/immutable/set"
)

var (
	_ immutable.Set[int]  = set.Of(1, 2, 3)
	_ collection.Set[int] = immutable.NewSet([]int{1, 2, 3})
)
