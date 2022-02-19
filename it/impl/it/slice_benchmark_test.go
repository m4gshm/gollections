package it

import "testing"

func Benchmark_IsValidIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := IsValidIndex(5, 0)
		r = IsValidIndex(5, 5)
		r = IsValidIndex(5, -1)
		_ = r
	}
}


func Benchmark_IsValidIndex2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := IsValidIndex2(5, 0)
		r = IsValidIndex2(5, 5)
		r = IsValidIndex2(5, -1)
		_ = r
	}
}

func Benchmark_CanIterateByRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := CanIterateByRange(NoStarted, 5, 4)
		r = CanIterateByRange(NoStarted, 5, 6)	
		r = CanIterateByRange(NoStarted, 5, NoStarted)
		_ = r
	}
}
