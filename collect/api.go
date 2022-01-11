package collect

import "github.com/m4gshm/gollections/typ"

type Collector[T any, OUT any] typ.Converter[typ.Iterator[T], OUT]

func Map[k comparable, v any](it typ.Iterator[*typ.KV[k, v]]) map[k]v {
	e := map[k]v{}
	for it.HasNext() {
		n, err := it.Next()
		if err != nil {
			panic(err)
		}
		e[n.Key()] = n.Value()
	}
	return e
}

func Groups[k comparable, v any](it typ.Iterator[*typ.KV[k, v]]) map[k][]v {
	e := map[k][]v{}
	for it.HasNext() {
		n, err := it.Next()
		if err != nil {
			panic(err)
		}
		key := n.Key()
		group := e[key]
		if group == nil {
			group = make([]v, 0)
		}
		e[key] = append(group, n.Value())
	}
	return e
}
