package client

// negamax_alpha returns a score and a turn
func negamax_alpha (state state, strategy string, alpha float32, player uint8, max_rec uint8) ([]Move, float32){
	var best_turn []Move  
	var max_score int = -1000000
	if max_rec == 0 {
		potstate = PotentialState {s: state, prob:1}
		score = scoreState(potstate)
		return (best_turn, score)
	}
	for move_list in turn_list(state, strategy, player){
		//playturn is a better name for evaluateOutcome, it is
		potstate_list = play_turn(state, player, move_list)
		
		score = 0
		for potstate in potstate_list {
			_, temp_score = negamax_alpha(potstate.s, strategy, max_score, 1- player, max_rec - 1 )
			score += temp_score * potstate.prob
		}
		score = -score 
		if score> max_score{
			max_score = score
			best_turn = move_list
		}
		//todo save previously modified cells by the play_turn function
		//we use a deep copy for now 
		//reverse_play_turn(state, player, move_list) dans le cas 
		//ou pas de bataille aleatoire pour eviter les copies
		
		if max_score > -alpha{
			break
		}
		
	}
	return best_turn, max_score
	
}




func turn_list (state state, strategy string, player uint8) [][]Move {
	//renvoie une liste de de liste de move 
	for cell in 	
	var nb_move int = 10
	var turn_list [][]Move = make([][]Move, 0)
	for i:=0; i<nb_move; i++ {
		turn_list[i]  = make([]Move, 0) 
		for j:= 0; j<nb_move; j++ {
			turn_list[i][j] = Move{ Start: Coordinates{X:i,Y:j}, N:1, End:Coordinates{X:i,Y:j+1}}
		}
	}
}

func play_turn(state state, player uint8, move_list []Move) []PotentialState {
	evaluateMoveOutcome ()
	for move in move_list
}

func reverse_play_turn(state state, player uint8, move_list []Move){
	//TODO
}

