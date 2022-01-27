//package collect provides util functions and types to implement the Collect method of a typ.Transformable implementation
package collect

import "github.com/m4gshm/gollections/typ"

type Collector[T any, OUT any] typ.Converter[typ.Iterator[T], OUT]
type CollectorKV[k, v any, OUT any] func(typ.KVIterator[k, v]) OUT

func Map[k comparable, v any](it typ.KVIterator[k, v]) map[k]v {
	e := map[k]v{}
	for it.HasNext() {
		key, val, err := it.Get()
		if err != nil {
			panic(err)
		}
		e[key] = val
	}
	return e
}

func Groups[k comparable, v any](it typ.KVIterator[k, v]) map[k][]v {
	e := map[k][]v{}
	for it.HasNext() {
		key, val, err := it.Get()
		if err != nil {
			panic(err)
		}
		group := e[key]
		if group == nil {
			group = make([]v, 0)
		}
		e[key] = append(group, val)
	}
	return e
}
