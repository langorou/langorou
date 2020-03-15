package tournament

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

type AIPlayer client.IA

type matchResult int

const (
	tie matchResult = iota
	player1Won
	player2Won
)

type MatchSummary struct {
	MapName                string
	Player1                AIPlayer
	Player2                AIPlayer
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
	return fmt.Sprintf("%s_VS_%sON%s", mr.Player1.Name(), mr.Player2.Name(), mr.MapName)

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

func playMap(
	mapPath string,
	isRand bool,
	randMapParams mapParams,
	timeoutS int,
	p1,
	p2 AIPlayer,
	matchSummaryCh chan MatchSummary,
	wg *sync.WaitGroup,
) error {

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
	if err != nil {
		return err
	}
	if err = player1.Init(); err != nil {
		return fmt.Errorf("fail to init player 1: %s", err)
	}

	player2, err := client.NewTCPClient(addr, p2.Name(), p2)
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
	competitors []AIPlayer,
	matchSummaryCh chan MatchSummary,
) {

	var wg sync.WaitGroup
	var randMapParams = newRandomMap(limits)

	for i, p1 := range competitors {
		for j, p2 := range competitors {
			if i != j {
				log.Printf("launching %s vs %s", p1.Name(), p2.Name())
				wg.Add(1)
				go playMap(
					mapPath,
					isRand,
					randMapParams,
					timeoutS,
					p1,
					p2,
					matchSummaryCh,
					&wg,
				)
			}
		}
	}

	wg.Wait()

}
