package mutable

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/collection/immutable/map_"
	"github.com/m4gshm/gollections/k"
)

var (
	_ immutable.Map[int, string]  = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
	_ collection.Map[int, string] = immutable.NewMapOf(map[int]string{1: "2", 2: "2", 3: "3"})
)
