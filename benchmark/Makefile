.PHONY: test
test:
	go test  -gcflags -m -bench=^\(Benchmark_Loop_ImmutableVector_FirstNext\)./...

.PHONY: prof
prof:
	# GOGC=off go test -cpuprofile cpu.prof -memprofile mem.prof -bench=^\(Benchmark_MapFlattStructure_Slice_PlainOld\|Benchmark_MapFlattStructure_SliceWithoutNilCheck\) .
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench=^Benchmark_Loop_ImmutableVector_FirstNext$$ .
	pprof -http=localhost:8080 ./cpu.prof

.PHONY: trace
trace:
	go test -trace trace.out -bench=^\(Benchmark_Loop_ImmutableVector_FirstNext\)
	go tool trace -http=localhost:8080 trace.out