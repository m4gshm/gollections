package map_

import (
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func Map[k comparable, v any, IT typ.KVIterator[k, v], kto comparable, vto any](elements IT, by typ.BiConverter[k, v, kto, vto]) typ.MapPipe[kto, vto, map[kto]vto] {
	return it.NewKVPipe(it.MapKV(elements, by), collect.Map[kto, vto])
}

func Filter[k comparable, v any, IT typ.KVIterator[k, v]](elements IT, filter typ.BiPredicate[k, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(elements, filter), collect.Map[k, v])
}

func Reduce[k comparable, v any, IT typ.KVIterator[k, v]](elements IT, by op.Quaternary[k, v]) (k, v) {
	return it.ReduceKV(elements, by)
}
