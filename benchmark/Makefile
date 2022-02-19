
.PHONY: all
all: prof


.PHONY: prof
prof:
	go test -cpuprofile cpu.prof -bench=^Benchmark_GetNext_Iterator_Immutable_Vector_Impl .
	pprof -svg .\benchmark.test.exe .\cpu.prof