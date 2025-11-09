package check

import (
	"testing"
)

func Benchmark_Nil(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = Nil[*int](nil)
	}
	b.StopTimer()
}

func Benchmark_Nil_StaticFunc(b *testing.B) {
	Nil := func(i *string) bool { return i == nil }

	b.ResetTimer()
	for b.Loop() {
		_ = Nil(nil)
	}
	b.StopTimer()
}

func Benchmark_Nil_StaticFuncInterface(b *testing.B) {
	Nil := func(i any) bool { return i == nil }

	for b.Loop() {
		_ = Nil(nil)
	}
	b.StopTimer()
}
