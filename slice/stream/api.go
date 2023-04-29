package stream

import (
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/slice/iter"
	"github.com/m4gshm/gollections/stream"
)

func Convert[FS ~[]From, From, To any](elements FS, by func(From) To) stream.Iter[To] {
	conv := iter.Convert(elements, by)
	return stream.New(conv.Next)
}

func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, by func(From) To) stream.Iter[To] {
	f := iter.FilterAndConvert(elements, filter, by)
	return stream.New(f.Next)
}

func Flatt[FS ~[]From, From, To any](elements FS, by func(From) []To) stream.Iter[To] {
	f := iter.Flatt(elements, by)
	return stream.New(f.Next)
}

func FilterAndFlatt[FS ~[]From, From, To any](elements FS, filter func(From) bool, flatt func(From) []To) stream.Iter[To] {
	f := iter.FilterAndFlatt(elements, filter, flatt)
	return stream.New(f.Next)
}

func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) stream.Iter[T] {
	f := iter.Filter(elements, filter)
	return stream.New(f.Next)
}

func NotNil[T any, TRS ~[]*T](elements TRS) stream.Iter[*T] {
	return Filter(elements, check.NotNil[T])
}
