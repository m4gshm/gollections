package test

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/map_"
)

var m = map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9}

func Benchmark_Keys(b *testing.B) {
	var k []int
	for b.Loop() {
		k = map_.Keys(m)
	}
	_ = k
}

func Benchmark_KeysConvert(b *testing.B) {
	var k []string
	for b.Loop() {
		k = map_.KeysConvert(m, strconv.Itoa)
	}
	_ = k
}

func Benchmark_KeysConvertAsIs(b *testing.B) {
	var k []int
	for b.Loop() {
		k = map_.KeysConvert(m, as.Is)
	}
	_ = k
}

func Benchmark_Values(b *testing.B) {
	var v []int
	for b.Loop() {
		v = map_.Values(m)
	}
	_ = v
}

func Benchmark_ValuesConvert(b *testing.B) {
	var k []string
	for b.Loop() {
		k = map_.ValuesConvert(m, strconv.Itoa)
	}
	_ = k
}

func Benchmark_ValuesConvertAsIs(b *testing.B) {
	var k []int
	for b.Loop() {
		k = map_.ValuesConvert(m, as.Is)
	}
	_ = k
}
