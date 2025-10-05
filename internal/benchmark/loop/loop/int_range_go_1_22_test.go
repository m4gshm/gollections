package loop

import "testing"

func Benchmark_int_Range_Iterating(b *testing.B) {
	for _, casee := range cases {
		b.Run(casee.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for v := range max {
					casee.load(v)
				}
			}
		})
	}
}
