package main

import (
	"flag"
	"io/ioutil"
	"log"
	"github.com/langorou/twilight/server"
)

var replayPath string

func init() {
	flag.StringVar(&replayPath, "replay", "", "path to the replay file")
}

func main() {
	flag.Parse()

	if replayPath == "" {
		log.Fatal("please specify a replay file path with -replay")
	}

	replayBytes, err := ioutil.ReadFile(replayPath)
	if err != nil {
		log.Fatalf("failed to read replay file: %s", err)
	}

	var replay 

	server.StartWebAppFromHistory(repl)

}
