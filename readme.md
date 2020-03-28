# Langorou

This is an IA for the game "Vampires VS Loup Garous" of the CentraleSupélec IA class.

The rules of the game are described [here](./doc/Projet.pdf).

A thorough documentation is available [here](./doc/documentation.md)

## Requirements

You will need to have `make` and a [`go` toolchain](https://golang.org/) installed (version >= 1.11 to support [go modules](https://blog.golang.org/using-go-modules)).

## Build and run

To build the project simply run `make`, this will create a binary at `build/langorou`, that you can then run:

`langorou -name <player_name> <host> <port>`
- the `-name` parameter is optional.
- `host` and `port` are the locations of the game server.

## Playing

Run `make auto` to launch a game, you can view it on [http://localhost:8080](http://localhost:8080)

## Tournament

You can launch a tournament on predefined maps with are located in [`maps/`](maps/).

Or you can launch a tournament on random maps with `go run cmd/tournoi/main.go`. You can configure the random maps generation with some flags (more details in the [`cmd/tournoi/main.go`](cmd/tournoi/main.go))

## Testing

To run the tests you can run: `make test`, by default this will run all the tests of this project.
To limit the tests you want to run you can do `make test pkg=./pkg/client` to run only the tests of the `pkg/client` package.

## Benchmarking

You can run benchmarks by running `make benchmark` a lot of parameters are available:
- `benchname` to specify which benchmark to run (by default all benchmarks are run)
- `profile` (by default `/tmp/profile.out` to specify where to save a profile of the benchmark, this can then be analyzed by [pprof](https://github.com/google/pprof) using: `go tool pprof -http localhost:6080 /tmp/profile.out.cpu` and going at http://localhost:6080 on your browser, it produces two profiles:
- - a cpu profile at `<profile>.cpu`
- - a memory profile at `<profile>.mem`
- `benchtime` to limit the benchmark time (format is for instance: `2s` or `500ms`)
- `pkg` to only run benchmarks of a specific package

## Replays

After a tournament (or a game if you saved it), you can replay the matches with `cmd/replay/main.go -replay "<path_to_replay>"`, or the more convenient `make replay replayPath="<path_to_replay>"` and analyse it at [http://localhost:8080](http://localhost:8080).

The initial position 0 isn't display, it starts after the first move.
