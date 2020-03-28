BINARY = langorou
GOARCH = amd64

# Enable go modules
GOCMD = GO111MODULE=on go

# Build the project
all: build

profile=/tmp/profile.out
benchname=.
benchtime=5s
pkg=./...
maps=./maps

.PHONY: build
build:
	mkdir -p bin
	${GOCMD} build -o ./bin/${BINARY} ./cmd/player/main.go

.PHONY: auto
auto:
	${GOCMD} run cmd/auto/main.go -rand

.PHONY: tournoi
tournoi:
	${GOCMD} run cmd/tournoi/main.go -mapFolder ${maps}

.PHONY: replay
replay:
	${GOCMD} run cmd/replay/main.go -replay ${replayPath}

.PHONY: test
test:
	${GOCMD} vet ${pkg}
	${GOCMD} test -v ${pkg}

.PHONY: benchmark
benchmark:
	${GOCMD} test -run=^$$ -bench=${benchname} -v -cpuprofile ${profile}.cpu -memprofile ${profile}.mem -benchmem -benchtime=${benchtime} -timeout=90s ${pkg}
	@echo "Benchmark profile saved at: ${profile}"

.PHONY: fmt
fmt:
	${GOCMD} fmt $$(go list ./... | grep -v /vendor/)

.PHONY: tidy
tidy:
	${GOCMD} mod tidy
