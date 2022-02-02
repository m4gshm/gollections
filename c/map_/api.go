package map_

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
)

//Map creates the key/value iterator that converts elements with a converter and returns them.
func Map[k comparable, v any, IT c.KVIterator[k, v], kto comparable, vto any](elements IT, by c.BiConverter[k, v, kto, vto]) c.MapPipe[kto, vto, map[kto]vto] {
	return it.NewKVPipe(it.MapKV(elements, by), collect.Map[kto, vto])
}

//Filter creates the key/value iterator that checks elements by filters and returns successful ones.
func Filter[k comparable, v any, IT c.KVIterator[k, v]](elements IT, filter c.BiPredicate[k, v]) c.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(elements, filter), collect.Map[k, v])
}

//Reduce reduces keys/value pairs to an one pair.
func Reduce[k comparable, v any, IT c.KVIterator[k, v]](elements IT, by op.Quaternary[k, v]) (k, v) {
	return it.ReduceKV(elements, by)
}
