package main

import (
	"flag"
	"github.com/langorou/langorou/pkg/client"
	"github.com/langorou/twilight/server"
	"log"
	"os"
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

func init() {
	flag.StringVar(&mapPath, "map", "", "path to the map to load (or save if randomly generating)")
	flag.BoolVar(&useRand, "rand", false, "use a randomly generated map")
	flag.IntVar(&rows, "rows", 10, "total number of rows")
	flag.IntVar(&columns, "columns", 10, "total number of columns")
	flag.IntVar(&humans, "humans", 16, "quantity of humans group")
	flag.IntVar(&monster, "monster", 8, "quantity of monster in the start case")
}


func main() {
	flag.Parse()
	log.Print("starting server...")
	go server.StartServer(mapPath, useRand, rows, columns, humans, monster)


	addr := "localhost:5555"
	player1, err := client.NewTCPClient(addr, "langone", client.NewDumbIA())
	failIf(err, "")
	player2, err := client.NewTCPClient(addr, "langtwo", client.NewDumbIA())
	failIf(err, "")

	go player1.Start()
	player2.Start()

	os.Exit(0)
}
