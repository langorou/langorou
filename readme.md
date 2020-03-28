# Langorou

## Playing

Run `make auto` to launch a game, you can view it on [http://localhost:8080](http://localhost:8080)

## Tournament

You can launch a tournament on predefined maps with are located in [`maps/`](maps/).

Or you can launch a tournament on random maps with `go run cmd/tournoi/main.go`. You can configure the random maps generation with some flags (more details in the [`cmd/tournoi/main.go`](cmd/tournoi/main.go))

## Replays

After a tournament (or a game if you saved it), you can replay the matches with `cmd/replay/main.go -replay "<path_to_replay>"`, or the more convenient `make replay replayPath="<path_to_replay>"` and analyse it at [http://localhost:8080](http://localhost:8080).

The initial position 0 isn't display, it starts after the first move.
