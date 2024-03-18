package mutable

import (
	//
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/collection/mutable/ordered/set"
)

var (
	_ *ordered.Set[int] = set.Of(1, 2, 3)
	_ *ordered.Set[int] = &ordered.Set[int]{}
)
