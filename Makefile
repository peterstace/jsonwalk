.PHONY: build
build:
	go install ./...
	go test ./...
	go vet ./...

.PHONY: bench
bench:
	go test -run=^\$$ -bench=. -benchmem -cpuprofile profile.cpu -memprofile profile.mem -benchtime=10s

.PHONY: pprof-cpu
pprof-cpu:
	go tool pprof -http 0.0.0.0:8080 profile.cpu

.PHONY: pprof-mem
pprof-mem:
	go tool pprof -http 0.0.0.0:8080 profile.mem

.PHONY: clean
clean:
	rm profile.* *.test -f
