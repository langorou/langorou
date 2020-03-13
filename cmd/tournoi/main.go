package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	_ "net/http/pprof"

	"github.com/langorou/langorou/pkg/client"
	"github.com/langorou/twilight/server"
)

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %v", msg, err)
	}
}

const (
	nRandMaps       = 5 // number of random maps to generate
	mapSizeMin      = 5
	mapSizeMax      = 50
	nHumanGroupsMin = 2
	nHumanGroupsMax = 30
	nMonsterMin     = 4
	nMonsterMax     = 60
)

var mapPath string
var mapFolder string
var useRand bool
var rows int
var columns int
var humans int
var monster int
var timeoutS int

func getMaps(root string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

func getRandIntInRange(min, max int) int {
	return min + rand.Intn(max+min+1)
}

type mapParams struct {
	rows, columns, humans, monsters int
}

func newRandomMap() mapParams {
	return mapParams{
		rows:     getRandIntInRange(mapSizeMin, mapSizeMax),
		columns:  getRandIntInRange(mapSizeMin, mapSizeMax),
		humans:   getRandIntInRange(nHumanGroupsMin, nHumanGroupsMax),
		monsters: getRandIntInRange(nMonsterMin, nMonsterMax),
	}
}

func init() {
	flag.StringVar(&mapPath, "map", "", "path to the map to load (or save if randomly generating)")
	flag.StringVar(&mapFolder, "mapFolder", "", "path for the folder in which we can test several maps")
	flag.BoolVar(&useRand, "rand", false, "use a randomly generated map")
	flag.IntVar(&rows, "rows", 10, "total number of rows")
	flag.IntVar(&columns, "columns", 10, "total number of columns")
	flag.IntVar(&humans, "humans", 16, "quantity of humans group")
	flag.IntVar(&monster, "monster", 8, "quantity of monster in the start case")
	flag.IntVar(&timeoutS, "timeout", 8, "timeout in seconds for each move")
}

type aiPlayer struct {
	ai client.IA
}

type matchResult struct {
	winner, looser    aiPlayer
	winEff, looserEff int
	endTurn           int
}

type tournamentResult []matchResult

func playMap(mapPath string, isRand bool, randMapParams mapParams, p1, p2 aiPlayer, matchResultCh chan matchResult) {

	portUsed := make(chan int, 1)

	go server.StartServer(mapPath, isRand, randMapParams.rows, randMapParams.columns, randMapParams.humans, randMapParams.monsters, time.Duration(timeoutS)*time.Second, true, portUsed, true)
	port := <-portUsed

	addr := fmt.Sprintf("localhost:%d", port)

	log.Printf("Launching a game for p1 and p2 on %s", addr)

	player1, err := client.NewTCPClient(addr, "langone", p1.ai)
	failIf(err, "")
	failIf(player1.Init(), "fail to init player 1")

	player2, err := client.NewTCPClient(addr, "langtwo", p2.ai)
	failIf(err, "")
	failIf(player2.Init(), "fail to init player 2")

	// TODO get match result
	go player1.Play()
	player2.Play()
	matchResultCh <- matchResult{}
}

func runTournamentOnMap(mapPath string, isRand bool, competitors []aiPlayer, matchResultCh chan matchResult) {

	var randMapParams = newRandomMap()

	for i, p1 := range competitors {
		for j, p2 := range competitors {
			if i != j {
				go playMap(mapPath, isRand, randMapParams, p1, p2, matchResultCh)
			}
		}

	}

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	competitors := []aiPlayer{
		aiPlayer{ai: client.NewMinMaxIA(7)},
		aiPlayer{ai: client.NewMinMaxIA(7)},
	}

	matchResultCh := make(chan matchResult)
	var leaderboard tournamentResult

	if mapFolder != "" {
		log.Printf("Using the maps provided for the tournament")

		mapPaths := getMaps(mapFolder)

		for _, mp := range mapPaths {
			runTournamentOnMap(mp, false, competitors, matchResultCh)
		}
	} else {
		for i := 0; i < nRandMaps; i++ {
			runTournamentOnMap("", true, competitors, matchResultCh)
		}
	}

	for res := range matchResultCh {
		leaderboard = append(leaderboard, res)
	}

	os.Exit(0)
}
