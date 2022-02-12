GO ?= go
BENCH_COUNT ?= 5
REF_NAME ?= $(shell git symbolic-ref HEAD --short | tr / - 2>/dev/null)

## Run benchmark, iterations count controlled by BENCH_COUNT, default 5.
bench:
	@$(GO) test -bench=. -count=$(BENCH_COUNT) -run=^a  ./... >bench-$(REF_NAME).txt
	@test -e bench-master.txt && benchstat bench-master.txt bench-$(REF_NAME).txt || benchstat bench-$(REF_NAME).txt

