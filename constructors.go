package examples

import (
	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/immutable/map_"
	"github.com/m4gshm/gollections/immutable/omap"
	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/immutable/vector"
)

func _() {
	vector.Of(1, 2, 3)
	vector.New([]int{1, 2, 3})

	set.Of(1, 2, 3)
	set.New([]int{1, 2, 3})
	oset.Of(1, 2, 3)
	oset.New([]int{1, 2, 3})

	map_.Of(K.V(1, "1"), K.V(2, "2"), K.V(3, "3"))
	map_.New(map[int]string{1: "2", 2: "2", 3: "3"})
	omap.Of(K.V(1, "1"), K.V(2, "2"), K.V(3, "3"))
	omap.New(map[int]string{1: "2", 2: "2", 3: "3"})
}
