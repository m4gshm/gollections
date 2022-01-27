package examples

import "github.com/m4gshm/gollections/immutable/vector"
import "github.com/m4gshm/gollections/immutable/set"
import "github.com/m4gshm/gollections/immutable/oset"
import "github.com/m4gshm/gollections/immutable/map_"
import "github.com/m4gshm/gollections/immutable/omap"
import "github.com/m4gshm/gollections/K"

func init() {
	vector.Of(1, 2, 3); vector.New([]int{1, 2, 3})
	
	set.Of(1, 2, 3); set.New([]int{1, 2, 3})
	oset.Of(1, 2, 3); oset.New([]int{1, 2, 3})

	map_.Of(K.V(1, "1"), K.V(2, "2"), K.V(3, "3")); map_.New(map[int]string{1: "2", 2: "2", 3: "3"})
	omap.Of(K.V(1, "1"), K.V(2, "2"), K.V(3, "3")); omap.New(map[int]string{1: "2", 2: "2", 3: "3"})
}
