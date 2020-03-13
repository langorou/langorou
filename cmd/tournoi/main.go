package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/langorou/langorou/pkg/utils"

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
	nRandMaps       = 1 // number of random maps to generate
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
	return fmt.Sprintf("%dx%d_h%d_m%d", m.rows, m.columns, m.humans, m.monsters)
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

type matchResult int

const (
	tie matchResult = iota
	player1Won
	player2Won
)

type matchSummary struct {
	MapName                string
	Player1                aiPlayer
	Player2                aiPlayer
	Winner                 matchResult
	Player1Eff, Player2Eff int
	EndTurn                int
	History                []server.Packed
}

func (mr *matchSummary) String() string {
	// First player is always werewolves, second Vampire
	switch mr.Winner {
	case tie:
		return fmt.Sprintf(
			"%-15s (P1) VS (P2) %-15s | %3d - %3d | %3d turns | %s",
			mr.Player1.Name(), mr.Player2.Name(), mr.Player1Eff, mr.Player2Eff, mr.EndTurn, mr.MapName,
		)
	case player1Won:
		return fmt.Sprintf(
			"%-15s (P1) VS (P2) %-15s | %3d - %3d | %3d turns | %s",
			mr.Player1.Name(), mr.Player2.Name(), mr.Player1Eff, mr.Player2Eff, mr.EndTurn, mr.MapName,
		)
	case player2Won:
		return fmt.Sprintf(
			"%-15s (P2) VS (P1) %-15s | %3d - %3d | %3d turns | %s",
			mr.Player2.Name(), mr.Player1.Name(), mr.Player2Eff, mr.Player1Eff, mr.EndTurn, mr.MapName,
		)
	}

	return fmt.Sprintf("error: invalid winner code %d", mr.Winner)
}

func (mr *matchSummary) shortName() string {
	// First player is always werewolves, second Vampire
	return fmt.Sprintf("%s_VS_%sON%s", mr.Player1.Name(), mr.Player2.Name(), mr.MapName)

}

func (mr *matchSummary) saveJSON(path string) {
	f, err := os.Create(path)
	failIf(err, "")
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(mr)
	failIf(err, "")
}

type tournamentResult []matchSummary

func (tr tournamentResult) matchResults() string {
	var output string
	for _, mr := range tr {
		output += mr.String()
		output += "\n"
	}
	return output
}

func (tr tournamentResult) leaderboard() string {
	// gagnant 3 points
	// perdant 0 points
	// égalité 1 point chacun
	leaderboard := make(map[string]int)

	for _, mr := range tr {
		switch mr.Winner {
		case tie:
			leaderboard[mr.Player1.Name()]++
			leaderboard[mr.Player2.Name()]++
		case player1Won:
			leaderboard[mr.Player1.Name()] += 3
			if _, ok := leaderboard[mr.Player2.Name()]; !ok {
				leaderboard[mr.Player2.Name()] = 0 // We might never add points to the looser
			}
		case player2Won:
			leaderboard[mr.Player2.Name()] += 3
			if _, ok := leaderboard[mr.Player1.Name()]; !ok {
				leaderboard[mr.Player1.Name()] = 0 // We might never add points to the looser
			}
		}

	}

	var output string
	for name, score := range leaderboard {
		// we'll not be sorted (random map)
		output += fmt.Sprintf("%15s - %3d points\n", name, score)
	}

	return output
}

func (tr tournamentResult) save(path string) {
	t := strconv.FormatInt(time.Now().Unix(), 10)

	f, err := os.Create(filepath.Join(path, fmt.Sprintf("%s_tournament.txt", t)))
	failIf(err, "")

	content := fmt.Sprintf("Leaderboard\n%s\n----------\n%s", tr.leaderboard(), tr.matchResults())

	_, err = f.WriteString(content)
	failIf(err, "")
	f.Sync()
	f.Close()

	dirPath := filepath.Join(path, fmt.Sprintf("%s_matches", t))
	failIf(utils.CreateDirIfNotExist(dirPath), "")

	for _, mr := range tr {
		filename := fmt.Sprintf("%s.json", mr.shortName())
		mr.saveJSON(filepath.Join(dirPath, filename))
	}
}

func playMap(mapPath string, isRand bool, randMapParams mapParams, p1, p2 aiPlayer, matchSummaryCh chan matchSummary, wg *sync.WaitGroup) {

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

	go player1.Play()
	player2.Play()

	outcome := <-gameOutcomeCh

	matchRes := matchSummary{
		EndTurn:    outcome.Turn,
		History:    outcome.History,
		Player1:    p1,
		Player2:    p2,
		Player1Eff: outcome.P1Eff,
		Player2Eff: outcome.P2Eff,
	}

	if isRand {
		matchRes.MapName = randMapParams.String()
	} else {
		matchRes.MapName = mapPath
	}

	switch {
	case outcome.P1Eff > outcome.P2Eff:
		matchRes.Winner = player1Won
	case outcome.P1Eff < outcome.P2Eff:
		matchRes.Winner = player2Won
	case outcome.P1Eff == outcome.P2Eff:
		matchRes.Winner = tie
	}

	matchSummaryCh <- matchRes

}

func runTournamentOnMap(mapPath string, isRand bool, competitors []aiPlayer, matchSummaryCh chan matchSummary) {

	var wg sync.WaitGroup
	var randMapParams = newRandomMap()

	for i, p1 := range competitors {
		for j, p2 := range competitors {
			if i != j {
				log.Printf("launching %s vs %s", p1.Name(), p2.Name())
				wg.Add(1)
				go playMap(mapPath, isRand, randMapParams, p1, p2, matchSummaryCh, &wg)
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
		// client.NewMinMaxIA(5),
		client.NewDumbIA(),
		// client.NewMinMaxIA(7),
	}

	matchSummaryCh := make(chan matchSummary)
	var leaderboard tournamentResult

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for res := range matchSummaryCh {
			leaderboard = append(leaderboard, res)
		}
		wg.Done()
	}(&wg)

	if mapFolder != "" {
		log.Printf("Using the maps provided for the tournament")

		mapPaths := getMaps(mapFolder)

		for _, mp := range mapPaths {
			// could use go on this, but generate two many games at the same time
			runTournamentOnMap(mp, false, competitors, matchSummaryCh)
		}
	} else {
		for i := 0; i < nRandMaps; i++ {
			runTournamentOnMap("", true, competitors, matchSummaryCh)
		}
	}
	close(matchSummaryCh)
	wg.Wait()

	log.Printf("\nGames summary\n--------\n%s\n", leaderboard.matchResults())
	log.Printf("\nFinal Scores\n--------\n%s", leaderboard.leaderboard())

	failIf(utils.CreateDirIfNotExist("./out"), "")
	leaderboard.save("./out/")

	os.Exit(0)
}
