package main

import (
	"flag"
	"log"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/langorou/langorou/pkg/client"
	"github.com/langorou/twilight/server"
)

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %v", msg, err)
	}
}

var mapPath string
var useRand bool
var rows int
var columns int
var humans int
var monster int
var timeoutS int

func init() {
	flag.StringVar(&mapPath, "map", "", "path to the map to load (or save if randomly generating)")
	flag.BoolVar(&useRand, "rand", false, "use a randomly generated map")
	flag.IntVar(&rows, "rows", 10, "total number of rows")
	flag.IntVar(&columns, "columns", 10, "total number of columns")
	flag.IntVar(&humans, "humans", 16, "quantity of humans group")
	flag.IntVar(&monster, "monster", 8, "quantity of monster in the start case")
	flag.IntVar(&timeoutS, "timeout", 8, "timeout in seconds for each move")
}

func main() {
	flag.Parse()

	log.Print("starting server...")

	// For profiling
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go server.StartServer(mapPath, useRand, rows, columns, humans, monster, time.Duration(timeoutS)*time.Second, false, nil, false)

	addr := "localhost:5555"
	player1, err := client.NewTCPClient(addr, "langone", client.NewMinMaxIA(7))
	failIf(err, "")
	failIf(player1.Init(), "fail to init player 1")

	player2, err := client.NewTCPClient(addr, "langtwo", client.NewMinMaxIA(5))
	failIf(err, "")
	failIf(player2.Init(), "fail to init player 2")

	go player1.Play()
	player2.Play()

	time.Sleep(5 * time.Minute)
	os.Exit(0)
}
