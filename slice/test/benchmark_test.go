package test

import (
	"testing"

	"github.com/m4gshm/gollections/notsafe"
)

type someType struct {
	field1 string
	field2 int64
}

func Benchmark_GetTypeSize(b *testing.B) {
	var size uintptr
	for b.Loop() {
		size = notsafe.GetTypeSize[someType]()
	}
	_ = size
}
