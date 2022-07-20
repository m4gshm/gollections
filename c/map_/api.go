package map_

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
)

//Map instantiates key/value iterator that converts elements with a converter and returns them.
func Map[K comparable, V any, IT c.KVIterator[K, V], kto comparable, vto any](elements IT, by c.BiConverter[K, V, kto, vto]) c.MapPipe[kto, vto, map[kto]vto] {
	return it.NewKVPipe(it.MapKV(elements, by), collect.Map[kto, vto])
}

//Filter instantiates key/value iterator that checks elements by filters and returns successful ones.
func Filter[K comparable, V any, IT c.KVIterator[K, V]](elements IT, filter c.BiPredicate[K, V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.FilterKV(elements, filter), collect.Map[K, V])
}

//Reduce reduces keys/value pairs to an one pair.
func Reduce[K comparable, V any, IT c.KVIterator[K, V]](elements IT, by op.Quaternary[K, V]) (K, V) {
	return it.ReduceKV(elements, by)
}
