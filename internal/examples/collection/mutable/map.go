package mutable

import (
	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/map_"
	"github.com/m4gshm/gollections/k"
)

var (
	_ *mutable.Map[int, string]   = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
	_ *mutable.Map[int, string] = mutable.NewMapOf(map[int]string{1: "2", 2: "2", 3: "3"})
)
