package tournament

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/langorou/langorou/pkg/client"
	"github.com/langorou/langorou/pkg/utils"
	"github.com/langorou/twilight/server"
)

type mapParams struct {
	rows, columns, humans, monsters int
}

func (m *mapParams) String() string {
	return fmt.Sprintf("%dx%d_h%d_m%d", m.rows, m.columns, m.humans, m.monsters)
}

func newRandomMap(limits RandMapLimits) mapParams {
	return mapParams{
		rows:     utils.GetRandIntInRange(limits.MapSizeMin, limits.MapSizeMax+1),
		columns:  utils.GetRandIntInRange(limits.MapSizeMin, limits.MapSizeMax+1),
		humans:   utils.GetRandIntInRange(limits.NHumanGroupsMin, limits.NHumanGroupsMax+1),
		monsters: utils.GetRandIntInRange(limits.NMonsterMin, limits.NMonsterMax+1),
	}
}

type Participant struct {
	Dumb    bool
	Timeout time.Duration
	Params  client.HeuristicParameters
}

func (p Participant) createPlayer() client.IA {
	if p.Dumb {
		return client.NewDumbIA()
	}

	return client.NewMinMaxIAP(p.Timeout, p.Params)
}

func (p Participant) Name() string {
	if p.Dumb {
		return "dumb IA"
	}

	return fmt.Sprintf("min_max_%d_%s", p.Timeout, p.Params.ShortString())
}

type matchResult int

const (
	tie matchResult = iota
	player1Won
	player2Won
)

type MatchSummary struct {
	MapName                string
	Player1                Participant
	Player2                Participant
	Winner                 matchResult
	Player1Eff, Player2Eff int
	EndTurn                int
	History                []server.Packed
}

func (mr *MatchSummary) String() string {
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

func (mr *MatchSummary) shortName() string {
	// First player is always werewolves, second Vampire
	return fmt.Sprintf("%s_VS_%sON%s", mr.Player1.Name(), mr.Player2.Name(), filepath.Base(mr.MapName))

}

func (mr *MatchSummary) saveJSON(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(mr)
	if err != nil {
		return err
	}

	return nil
}

type TournamentResult []MatchSummary

func (tr TournamentResult) MatchResults() string {
	var output string
	for _, mr := range tr {
		output += mr.String()
		output += "\n"
	}
	return output
}

func (tr TournamentResult) Leaderboard() string {
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

func (tr TournamentResult) Save(path string) error {
	t := strconv.FormatInt(time.Now().Unix(), 10)

	f, err := os.Create(filepath.Join(path, fmt.Sprintf("%s_tournament.txt", t)))
	if err != nil {
		return err
	}

	content := fmt.Sprintf("Leaderboard\n%s\n----------\n%s", tr.Leaderboard(), tr.MatchResults())

	if _, err = f.WriteString(content); err != nil {
		return err
	}

	f.Sync()
	f.Close()

	dirPath := filepath.Join(path, fmt.Sprintf("%s_matches", t))
	if err = utils.CreateDirIfNotExist(dirPath); err != nil {
		return err
	}

	for _, mr := range tr {
		filename := fmt.Sprintf("%s.json", mr.shortName())
		if err = mr.saveJSON(filepath.Join(dirPath, filename)); err != nil {
			return err
		}
	}

	return nil
}

type job interface {
	execute() error
}

type playMap struct {
	mapPath        string
	isRand         bool
	randMapParams  mapParams
	timeoutS       int
	p1             Participant
	p2             Participant
	matchSummaryCh chan MatchSummary
	wg             *sync.WaitGroup
}

func (pm playMap) execute() error {

	defer pm.wg.Done()

	portUsed := make(chan int, 1)
	gameOutcomeCh := make(chan server.GameOutcome, 1)

	go server.StartServer(
		pm.mapPath,
		pm.isRand,
		pm.randMapParams.rows,
		pm.randMapParams.columns,
		pm.randMapParams.humans,
		pm.randMapParams.monsters,
		time.Duration(pm.timeoutS)*time.Second,
		true,
		portUsed,
		true,
		gameOutcomeCh,
	)

	port := <-portUsed

	addr := fmt.Sprintf("localhost:%d", port)

	log.Printf("Launching %s vs %s on %s", pm.p1.Name(), pm.p2.Name(), addr)

	player1, err := client.NewTCPClient(addr, pm.p1.Name(), pm.p1.createPlayer())
	if err != nil {
		return err
	}
	if err = player1.Init(); err != nil {
		return fmt.Errorf("fail to init player 1: %s", err)
	}

	player2, err := client.NewTCPClient(addr, pm.p2.Name(), pm.p2.createPlayer())
	if err != nil {
		return err
	}
	if err = player2.Init(); err != nil {
		return fmt.Errorf("fail to init player 2: %s", err)
	}

	go player1.Play()
	player2.Play()

	outcome := <-gameOutcomeCh

	matchRes := MatchSummary{
		EndTurn:    outcome.Turn,
		History:    outcome.History,
		Player1:    pm.p1,
		Player2:    pm.p2,
		Player1Eff: outcome.P1Eff,
		Player2Eff: outcome.P2Eff,
	}

	if pm.isRand {
		matchRes.MapName = pm.randMapParams.String()
	} else {
		matchRes.MapName = pm.mapPath
	}

	switch {
	case outcome.P1Eff > outcome.P2Eff:
		matchRes.Winner = player1Won
	case outcome.P1Eff < outcome.P2Eff:
		matchRes.Winner = player2Won
	case outcome.P1Eff == outcome.P2Eff:
		matchRes.Winner = tie
	}

	pm.matchSummaryCh <- matchRes

	return nil
}

type RandMapLimits struct {
	MapSizeMin      int
	MapSizeMax      int
	NHumanGroupsMin int
	NHumanGroupsMax int
	NMonsterMin     int
	NMonsterMax     int
}

func RunTournamentOnMap(
	mapPath string,
	isRand bool,
	limits RandMapLimits,
	timeoutS int,
	competitors []Participant,
	matchSummaryCh chan MatchSummary,
) {

	var wg sync.WaitGroup
	var randMapParams = newRandomMap(limits)

	maxConcurrentPlay := int(math.Max(1, float64(runtime.NumCPU()-1)))

	concurrentPlays := make(chan job)
	log.Printf("Launching %d games at the same time.", maxConcurrentPlay)

	for i := 0; i < maxConcurrentPlay; i++ {
		go func(id int) {
			for j := range concurrentPlays {
				log.Printf("Starting a game with worker %d", id)
				if err := j.execute(); err != nil {
					log.Print(err)
				}
				log.Printf("Finished a game with worker %d", id)
			}
		}(i)
	}

	for i, p1 := range competitors {
		for j, p2 := range competitors {
			if i != j {
				wg.Add(1)
				concurrentPlays <- playMap{
					mapPath,
					isRand,
					randMapParams,
					timeoutS,
					p1,
					p2,
					matchSummaryCh,
					&wg,
				}
			}
		}
	}
	close(concurrentPlays)

	wg.Wait()

}
