package mutable

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/collection/immutable/ordered/set"
)

var (
	_ ordered.Set[int]    = set.Of(1, 2, 3)
	_ collection.Set[int] = ordered.NewSet([]int{1, 2, 3})
)
