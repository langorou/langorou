BINARY = langorou
GOARCH = amd64

# Enable go modules
GOCMD = GO111MODULE=on go

# Build the project
all: build

profile=/tmp/profile.out
bench_name=.

.PHONY: build
build:
	mkdir -p bin
	${GOCMD} build -o ./bin/${BINARY} ./cmd/player/main.go

.PHONY: auto
auto:
	${GOCMD} run cmd/auto/main.go -rand

.PHONY: test
test:
	${GOCMD} test -v ./...

.PHONY: benchmark
benchmark:
	${GOCMD} test -bench=${bench_name} -v -cpuprofile ${profile} ./pkg/client/
	@echo "Benchmark profile saved at: ${profile}"

.PHONY: fmt
fmt:
	${GOCMD} fmt $$(go list ./... | grep -v /vendor/)

.PHONY: tidy
tidy:
	${GOCMD} mod tidy
