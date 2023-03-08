package main

import "fmt"

func main() {

	// remotePlayer := &RemotePlayer{}

	// game := Game{}
	// game.Init(remotePlayer, &RandomComputerPlayer{2})

	// remotePlayer.Init(1, &game.Board)

	// game.Print()
	// for {
	// 	fmt.Println("")
	// 	fmt.Println(game.Board.FullmoveNumber)
	// 	game.Print()
	// 	isEnd := game.Play()
	// 	if isEnd {
	// 		break
	// 	}
	// }
	// game.Print()

	//test()

	createServer()

}

func test() {
	b := InitFEN("1nb1kb2/rp2p1P1/1qp2p1P/3p4/4PQ2/8/PbPP1P2/1NB1KBNR w KQkq - 0 1")
	b.PrintWithBorder()

	_, location := GridToLocation("g7")
	bishop := b.GetPieceAtLocation(location)
	fmt.Println(bishop)
	// Get the valid moves
	moves := []Move{}
	for _, move := range b.GetPlayerLegalMoves(1) {
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
