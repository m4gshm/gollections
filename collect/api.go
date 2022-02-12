//Package collect provides util functions and types to implement the Collect method of a c.Transformable implementation
package collect

import "github.com/m4gshm/gollections/c"

//Collector is Converter of Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result.
type Collector[t any, out any] c.Converter[c.Iterator[t], out]

//CollectorKV is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result.
type CollectorKV[k, v any, out any] func(c.KVIterator[k, v]) out

//Map collects the map of key/value pairs obtained by passing over a key/value iterator.
func Map[k comparable, v any](it c.KVIterator[k, v]) map[k]v {
	e := map[k]v{}
	for it.HasNext() {
		key, val := it.Get()
		e[key] = val
	}
	return e
}

//Groups collects sets of values grouped by keys obtained by passing a key/value iterator.
func Groups[k comparable, v any](it c.KVIterator[k, v]) map[k][]v {
	e := map[k][]v{}
	for it.HasNext() {
		key, val := it.Get()
		group := e[key]
		if group == nil {
			group = make([]v, 0)
		}
		e[key] = append(group, val)
	}
	return e
}
