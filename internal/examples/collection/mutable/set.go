package mutable

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/set"
)

var (
	_ *mutable.Set[int]                          = set.Of(1, 2, 3)
	_ collection.Set[int, *mutable.SetIter[int]] = &mutable.Set[int]{}
)
