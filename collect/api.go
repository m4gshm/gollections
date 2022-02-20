//Package collect provides util functions and types to implement the Collect method of a c.Transformable implementation
package collect

import "github.com/m4gshm/gollections/c"

//Collector is Converter of Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result.
type Collector[t any, out any] c.Converter[c.Iterator[t], out]

//CollectorKV is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result.
type CollectorKV[k, v any, out any] func(c.KVIterator[k, v]) out

//Map collects the map of key/value pairs obtained by passing over a key/value iterator.
func Map[K comparable, V any](it c.KVIterator[K, V]) map[K]V {
	e := map[K]V{}
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		e[k] = v
	}
	return e
}

//Groups collects sets of values grouped by keys obtained by passing a key/value iterator.
func Groups[K comparable, V any](it c.KVIterator[K, V]) map[K][]V {
	e := map[K][]V{}
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		e[k] = append(e[k], v)
	}
	return e
}
