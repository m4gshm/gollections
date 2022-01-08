package check

import (
	"testing"
)

func Benchmark_Empty(b *testing.B) {
	v := []int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Empty(v)
	}
	b.StopTimer()
}

func Benchmark_Nil(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Nil[*int](nil)
	}
	b.StopTimer()
}

func Benchmark_Nil_StaticFunc(b *testing.B) {
	Nil := func(i *string) bool { return i == nil }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Nil(nil)
	}
	b.StopTimer()
}

func Benchmark_Nil_StaticFuncInterface(b *testing.B) {

	Nil := func(i interface{}) bool { return i == nil }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Nil(nil)
	}
	b.StopTimer()
}
