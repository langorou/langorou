BINARY = langorou
GOARCH = amd64

# Enable go modules
GOCMD = GO111MODULE=on go

# Build the project
all: build

profile=/tmp/profile.out
bench_name=.
bench_time=5s
pkg=./...

.PHONY: build
build:
	mkdir -p bin
	${GOCMD} build -o ./bin/${BINARY} ./cmd/player/main.go

.PHONY: auto
auto:
	${GOCMD} run cmd/auto/main.go -rand

.PHONY: tournoi
tournoi:
	${GOCMD} run cmd/tournoi/main.go

.PHONY: replay
replay:
	${GOCMD} run cmd/replay/main.go -replay

.PHONY: test
test:
	${GOCMD} test -v ${pkg}

.PHONY: benchmark
benchmark:
	${GOCMD} test -run=^$$ -bench=${bench_name} -v -cpuprofile ${profile} -benchmem -benchtime=${bench_time} -timeout=30s ${pkg}
	@echo "Benchmark profile saved at: ${profile}"

.PHONY: fmt
fmt:
	${GOCMD} fmt $$(go list ./... | grep -v /vendor/)

.PHONY: tidy
tidy:
	${GOCMD} mod tidy
