package client

// negamaxAlpha returns a score and a turn
func negamaxAlpha(state state, strategy string, alpha float64, player uint8, maxRec uint8) ([]Move, float64) {
	bestTurn := []Move{}
	maxScore := -1000000.0

	if maxRec == 0 {
		potentialState := PotentialState{s: state, prob: 1}
		return bestTurn, scoreState(potentialState)
	}

	for _, moveList := range getMoves(state, strategy, player) {
		//playturn is a better name for evaluateOutcome, it is
		potentialState := applyMoves(state, player, moveList)

		score := 0.0
		for _, potentialState := range potentialState {
			_, tmpScore := negamaxAlpha(potentialState.s, strategy, maxScore, 1-player, maxRec-1)
			score += tmpScore * potentialState.prob
		}
		score = -score
		if score > maxScore {
			maxScore = score
			bestTurn = moveList
		}
		//todo save previously modified cells by the applyMoves function
		//we use a deep copy for now
		//undoMoves(state, player, moveList) dans le cas
		//ou pas de bataille aleatoire pour eviter les copies

		if maxScore > -alpha {
			break
		}

	}
	return bestTurn, maxScore

func get_possible_moves_from_cell(coord Coordinates){
	//TODO returns possible moves from a cell
	//we don't consider the splits for now
}

// getMoves returns a list of moves
func getMoves(state state, strategy string, player uint8) [][]Move {
	max := 10
	movesMat := make([][]Move, max)

	for i := range movesMat {
		movesMat[i] = []Move{}

func turn_list (state state, strategy string, playerRace uint8) [][]Move {
	//renvoie une liste de de liste de move
	var output [][]Moves
	
	// cartesian product of all moves
	for i, row := range state{
		for j, cell := range row{
			if cell.cout > 0 && cell.race == player
				var output2 [][]Moves
				for _, move in get_possible_moves_from_cell(Coordinates{X:j,Y:i}) {
					for _, move_list := range output {
						output2 = append(output2, append(move_list,move))
				}
				output = output2
			}
		} 
		}
	}



func undoMoves(state state, player uint8, moveList []Move) {
	//TODO
}
