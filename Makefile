.PHONY: build
build:
	go install ./...
	go test ./...
	go vet ./...

.PHONY: bench
bench:
	go test -run=^\$$ -bench=. -benchmem -cpuprofile profile.cpu -memprofile profile.mem -benchtime=10s

.PHONY: clean
clean:
	rm profile.* *.test -f
