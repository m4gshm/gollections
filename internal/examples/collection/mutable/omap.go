package mutable

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/collection/mutable/ordered/map_"
	"github.com/m4gshm/gollections/k"
)

var (
	_ *ordered.Map[int, string]   = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
	_ collection.Map[int, string] = ordered.NewMap(
		/*order  */ []int{3, 1, 2},
		/*uniques*/ map[int]string{1: "2", 2: "2", 3: "3"},
	)
)
