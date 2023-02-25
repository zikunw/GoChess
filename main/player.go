package main

// player is the struct that represents a player in the game.
// The game sturct is going ask this struct for the next move.

// Therefore we can have different types of players:
// - human player
// - computer player (and more variations of alg)

type PlayerController interface {
	GetMove() Move
}

type HumanPlayer struct {
	Color int // 1 - white, 2 - black
}

//func (p *HumanPlayer) GetMove(b Board) Move {
//
//	reader := bufio.NewReader(os.Stdin)
//
//	// Prompt user for which piece to move
//	// If the piece is not owned by the player, prompt again
//	for {
//		fmt.Print("Which piece do you want to move? (e.g. a2): ")
//		text, _ := reader.ReadString('\n')
//	}
//}
