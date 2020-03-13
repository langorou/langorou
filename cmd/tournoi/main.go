package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
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
	mapSizeMax      = 40
	nHumanGroupsMin = 2
	nHumanGroupsMax = 30
	nMonsterMin     = 4
	nMonsterMax     = 40
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

func (m *mapParams) String() string {
	return fmt.Sprintf("(%2dx%2d) - %2dhums - %2dmonst", m.rows, m.columns, m.humans, m.monsters)
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

type aiPlayer client.IA

type matchResult struct {
	mapName           string
	winner, looser    aiPlayer
	isTie             bool
	winEff, looserEff int
	endTurn           int
}

func (mr *matchResult) String() string {
	return fmt.Sprintf(
		"%-15s VS %-15s | %3d - %3d | %3d turns | %s",
		mr.winner.Name(), mr.looser.Name(), mr.winEff, mr.looserEff, mr.endTurn, mr.mapName,
	)
}

type tournamentResult []matchResult

func (tr tournamentResult) printMatchResults() {
	for _, mr := range tr {
		log.Print(mr.String())
	}
}

func (tr tournamentResult) printLeaderboard() {
	// gagnant 3 points
	// perdant 0 points
	// égalité 1 point chacun
	leaderboard := make(map[string]int)

	for _, mr := range tr {
		if _, ok := leaderboard[mr.looser.Name()]; !ok {
			leaderboard[mr.looser.Name()] = 0 // We might never add points to the looser
		}

		if mr.isTie {
			leaderboard[mr.winner.Name()]++
			leaderboard[mr.looser.Name()]++
		} else {
			leaderboard[mr.winner.Name()] += 3
		}
	}

	for name, score := range leaderboard {
		// we'll not be sorted (random map)
		log.Printf("%15s - %3d points", name, score)
	}
}

func playMap(mapPath string, isRand bool, randMapParams mapParams, p1, p2 aiPlayer, matchResultCh chan matchResult, wg *sync.WaitGroup) {

	defer wg.Done()

	portUsed := make(chan int, 1)
	gameOutcomeCh := make(chan server.GameOutcome, 1)

	go server.StartServer(
		mapPath,
		isRand,
		randMapParams.rows,
		randMapParams.columns,
		randMapParams.humans,
		randMapParams.monsters,
		time.Duration(timeoutS)*time.Second,
		true,
		portUsed,
		true,
		gameOutcomeCh,
	)

	port := <-portUsed

	addr := fmt.Sprintf("localhost:%d", port)

	log.Printf("Launching a game for p1 and p2 on %s", addr)

	player1, err := client.NewTCPClient(addr, p1.Name(), p1)
	failIf(err, "")
	failIf(player1.Init(), "fail to init player 1")

	player2, err := client.NewTCPClient(addr, p2.Name(), p2)
	failIf(err, "")
	failIf(player2.Init(), "fail to init player 2")

	// TODO get match result
	go player1.Play()
	player2.Play()

	outcome := <-gameOutcomeCh

	matchRes := matchResult{
		endTurn: outcome.Turn,
	}

	if isRand {
		matchRes.mapName = randMapParams.String()
	} else {
		matchRes.mapName = mapPath
	}

	switch {
	case outcome.P1Eff > outcome.P2Eff:
		matchRes.winner = p1
		matchRes.looser = p2
		matchRes.winEff = outcome.P1Eff
		matchRes.looserEff = outcome.P2Eff
	case outcome.P1Eff < outcome.P2Eff:
		matchRes.winner = p2
		matchRes.looser = p1
		matchRes.winEff = outcome.P2Eff
		matchRes.looserEff = outcome.P1Eff

	case outcome.P1Eff == outcome.P2Eff:
		matchRes.winner = p2
		matchRes.looser = p1
		matchRes.winEff = outcome.P2Eff
		matchRes.looserEff = outcome.P1Eff
		matchRes.isTie = true
	default:
		log.Fatalf("invalid state for %s vs %s: %d", p1.Name(), p2.Name(), outcome)
	}

	matchResultCh <- matchRes

}

func runTournamentOnMap(mapPath string, isRand bool, competitors []aiPlayer, matchResultCh chan matchResult) {

	var wg sync.WaitGroup
	var randMapParams = newRandomMap()

	for i, p1 := range competitors {
		for j, p2 := range competitors {
			if i != j {
				log.Printf("launching %s vs %s", p1.Name(), p2.Name())
				wg.Add(1)
				go playMap(mapPath, isRand, randMapParams, p1, p2, matchResultCh, &wg)
			}
		}
	}

	wg.Wait()

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	competitors := []aiPlayer{
		client.NewMinMaxIA(2),
		client.NewMinMaxIA(5),
		client.NewDumbIA(),
		client.NewMinMaxIA(7),
	}

	matchResultCh := make(chan matchResult)
	var leaderboard tournamentResult

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for res := range matchResultCh {
			leaderboard = append(leaderboard, res)
		}
		wg.Done()
	}(&wg)

	if mapFolder != "" {
		log.Printf("Using the maps provided for the tournament")

		mapPaths := getMaps(mapFolder)

		for _, mp := range mapPaths {
			// could use go on this, but generate two many games at the same time
			runTournamentOnMap(mp, false, competitors, matchResultCh)
		}
	} else {
		for i := 0; i < nRandMaps; i++ {
			runTournamentOnMap("", true, competitors, matchResultCh)
		}
	}
	close(matchResultCh)
	wg.Wait()

	log.Printf("Games summary")
	log.Printf("--------")
	leaderboard.printMatchResults()
	log.Println()
	log.Printf("Final Scores")
	log.Printf("--------")
	leaderboard.printLeaderboard()

	os.Exit(0)
}
