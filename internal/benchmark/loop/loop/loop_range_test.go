package loop

import (
	"testing"

	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
)

var max = 10000

func HighLoad(v int) {
	resultInt = v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v * v
}

var resultInt = 0

func LowLoad(v int) {
	resultInt = v * v * v
}

func MidLoad(v int) {
	resultInt = v * v * v * v * v * v * v * v * v
}

type benchCase struct {
	name string
	load func(int)
}

var cases = []benchCase{ /*{"high", HighLoad}, {"mid", MidLoad},*/ {"low", LowLoad}}

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

func Benchmark_SeqRange_Iterating(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for v := range seq.Range(0, max) {
					casee.load(v)
				}
			}
		})
	}
}
