package mutable

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/collection/immutable/oset"
)

var (
	_ ordered.Set[int]    = oset.Of(1, 2, 3)
	_ collection.Set[int] = ordered.NewSet([]int{1, 2, 3})
)
