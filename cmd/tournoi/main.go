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

	competitors := generatePlayers()

	matchSummaryCh := make(chan tournament.MatchSummary)
	var leaderboard tournament.Result

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

func generatePlayers() []tournament.Participant {
	dur := 1 * time.Second

	players := []tournament.Participant{
		{Dumb: true},
		{Timeout: dur, Params: client.NewDefaultHeuristicParameters()},
		{Timeout: dur, Params: client.HeuristicParameters{
			Counts:           1,
			Battles:          0.5,
			NeutralBattles:   0.5,
			CumScore:         client.DefaultCumScore,
			WinScore:         client.DefaultWinScore,
			LoseOverWinRatio: 1,
			WinThreshold:     1,
			MaxGroups:        3,
			Groups:           0,
		}},
		{Timeout: dur, Params: client.HeuristicParameters{
			Counts:           1,
			Battles:          0.2,
			NeutralBattles:   0.4,
			CumScore:         client.DefaultCumScore,
			WinScore:         client.DefaultWinScore,
			LoseOverWinRatio: 1,
			WinThreshold:     1,
			MaxGroups:        3,
			Groups:           0,
		}},
		{Timeout: dur, Params: client.HeuristicParameters{
			Counts:           1,
			Battles:          0.2,
			NeutralBattles:   0.2,
			CumScore:         client.DefaultCumScore,
			WinScore:         client.DefaultWinScore,
			LoseOverWinRatio: 1,
			WinThreshold:     1,
			MaxGroups:        3,
			Groups:           0,
		}},
		{Timeout: dur, Params: client.HeuristicParameters{
			Counts:           1,
			Battles:          0.05,
			NeutralBattles:   0.05,
			CumScore:         client.DefaultCumScore,
			WinScore:         client.DefaultWinScore,
			LoseOverWinRatio: 1,
			WinThreshold:     1,
			MaxGroups:        3,
			Groups:           0,
		}},
		{Timeout: dur, Params: client.HeuristicParameters{
			Counts:           1,
			Battles:          0.02,
			NeutralBattles:   0.03,
			CumScore:         client.DefaultCumScore,
			WinScore:         client.DefaultWinScore,
			LoseOverWinRatio: 1,
			WinThreshold:     0.8,
			MaxGroups:        3,
			Groups:           0,
		}},
		{Timeout: dur, Params: client.HeuristicParameters{
			Counts:           1,
			Battles:          0.,
			NeutralBattles:   0.,
			CumScore:         0,
			WinScore:         client.DefaultWinScore,
			LoseOverWinRatio: 1,
			WinThreshold:     1,
			MaxGroups:        3,
			Groups:           -0.001,
		}},
		{Timeout: dur, Params: client.HeuristicParameters{
			// Not risk averse at all
			Counts:           1,
			Battles:          client.DefaultBattles,
			NeutralBattles:   client.DefaultNeutralBattles,
			CumScore:         client.DefaultCumScore,
			WinScore:         client.DefaultWinScore,
			LoseOverWinRatio: 0.8,
			WinThreshold:     0.8,
			MaxGroups:        2,
			Groups:           client.DefaultGroups,
		}},
		{Timeout: dur, Params: client.HeuristicParameters{
			// risk averse
			Counts:           1,
			Battles:          client.DefaultBattles,
			NeutralBattles:   client.DefaultNeutralBattles,
			CumScore:         client.DefaultCumScore,
			WinScore:         client.DefaultWinScore,
			LoseOverWinRatio: 1.2,
			WinThreshold:     1,
			MaxGroups:        2,
			Groups:           client.DefaultGroups,
		}},
	}

	return players
}
