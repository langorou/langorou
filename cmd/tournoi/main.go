package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/langorou/langorou/pkg/utils"

	_ "net/http/pprof"

	"github.com/langorou/langorou/pkg/client"
	"github.com/langorou/langorou/pkg/tournament"
)

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %v", msg, err)
	}
}

const nRandMaps = 1 // number of random maps to generate

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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	competitors := []tournament.AIPlayer{
		// client.NewMinMaxIA(200 * time.Millisecond),
		// client.NewMinMaxIA(500 * time.Millisecond),
		client.NewDumbIA(),
		client.NewMinMaxIA(1500 * time.Millisecond),
	}

	matchSummaryCh := make(chan tournament.MatchSummary)
	var leaderboard tournament.TournamentResult

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
			if strings.HasSuffix(mp, ".xml") {
				log.Printf("Launching tournament on map %s", mp)
				tournament.RunTournamentOnMap(mp, false, tournament.RandMapLimits{}, timeoutS, competitors, matchSummaryCh)
			}
		}
	} else {
		limits := tournament.RandMapLimits{
			MapSizeMin:      10,
			MapSizeMax:      16,
			NHumanGroupsMin: 2,
			NHumanGroupsMax: 30,
			NMonsterMin:     4,
			NMonsterMax:     40,
		}

		for i := 0; i < nRandMaps; i++ {
			tournament.RunTournamentOnMap("", true, limits, timeoutS, competitors, matchSummaryCh)
		}
	}
	close(matchSummaryCh)
	wg.Wait()

	log.Printf("\nGames summary\n--------\n%s\n", leaderboard.MatchResults())
	log.Printf("\nFinal Scores\n--------\n%s", leaderboard.Leaderboard())

	failIf(utils.CreateDirIfNotExist("./out"), "")
	failIf(leaderboard.Save("./out/"), "saving")

	os.Exit(0)
}
