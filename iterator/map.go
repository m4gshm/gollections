package iterator

import (
	"runtime"
	"sync"
)

func newMap[K comparable, V any](values map[K]V) *Map[K, V] {
	m := &Map[K, V]{
		values: values,
		flow:   make(chan *KV[K, V], 1),
	}
	m.start()
	return m
}

type Map[K comparable, V any] struct {
	values  map[K]V
	current *KV[K, V]
	flow    chan *KV[K, V]
}

type KV[K comparable, V any] struct {
	key   K
	value V
}

func (k *KV[K, V]) Key() K {
	return k.key
}

func (k *KV[K, V]) Value() V {
	return k.value
}

var _ Iterator[*KV[interface{}, interface{}]] = (*Map[interface{}, interface{}])(nil)

func (s *Map[K, V]) Next() bool {
	var (
		ok   bool
		stop bool
	)
	for !stop {
		select {
		case s.current, ok = <-s.flow:
			stop = true
		default:
			runtime.Gosched()
		}
	}
	return ok
}

func (s *Map[K, V]) Get() *KV[K, V] {
	return s.current
}

func (s *Map[K, V]) start() {
	waiter := sync.WaitGroup{}
	waiter.Add(1)
	go func() {
		waiter.Done()
		for k, v := range s.values {
			s.flow <- &KV[K, V]{key: k, value: v}
			runtime.Gosched()
		}
		close(s.flow)
	}()
	waiter.Wait()
}
