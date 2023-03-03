package main

import "fmt"

func main() {

	remotePlayer := &RemotePlayer{}

	game := Game{}
	game.Init(remotePlayer, &RandomComputerPlayer{2})

	remotePlayer.Init(1, &game.Board)

	game.Print()
	for {
		fmt.Println("")
		fmt.Println(game.Board.FullmoveNumber)
		game.Print()
		isEnd := game.Play()
		if isEnd {
			break
		}
	}
	game.Print()

}

func test() {
	b := InitFEN("rnbqkbnr/ppppppp1/8/7p/2B1P3/8/PPPP1PPP/RNBQK1NR w KQkq - 0 1")
	b.PrintWithBorder()

	_, location := GridToLocation("c4")
	bishop := b.GetPieceAtLocation(location)
	fmt.Println(bishop)
	// Get the valid moves
	moves := []Move{}
	for _, move := range b.GetPlayerLegalMoves(1) {
		fmt.Println(move.ToString())
		if move.From.Equals(location) {
			moves = append(moves, move)
		}
	}

	validSquare := []string{}
	for _, move := range moves {
		validSquare = append(validSquare, LocationToGrid(move.To))
	}

	fmt.Println(validSquare)
}
