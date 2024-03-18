package mutable

import (
	//
	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/set"
)

var (
	_ *mutable.Set[int] = set.Of(1, 2, 3)
	_ *mutable.Set[int] = &mutable.Set[int]{}
)
