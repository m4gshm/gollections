package loop

import "testing"

func Benchmark_int_Range_Iterating(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for b.Loop() {
				for v := range maxValOfRange {
					casee.load(v)
				}
			}
		})
	}
}
