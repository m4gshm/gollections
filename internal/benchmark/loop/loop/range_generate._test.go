package loop

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

var max = 10000

var resultStr = ""

func HighLoad(v int) {
	resultStr = strconv.Itoa(v)
}

var resultInt = 0

func LowLoad(v int) {
	resultInt = v * v
}

func MidLoad(v int) {
	resultInt = v * v * v * v * v * v * v * v * v
}

type benchCase struct {
	name string
	load func(int)
}

var cases = []benchCase{{"high", HighLoad}, {"mid", MidLoad}, {"low", LowLoad}}

func Benchmark_LoopRange_Iterating(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := loop.Range(0, max)
				for v, ok := next(); ok; v, ok = next() {
					casee.load(v)
				}
			}
		})
	}
}

func Benchmark_LoopRange_Iterating2(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				next := loop.Range(0, max)
				v, ok := next()
				for ok {
					casee.load(v)
					v, ok = next()
				}
			}
		})
	}
}

func Benchmark_SliceRange_Iterating(b *testing.B) {
	integers := slice.Range(0, max)
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range integers {
					casee.load(v)
				}
			}
		})
	}
}
