package main

import "fmt"

type Game struct {
	Board Board
	State GameState
}

// GameState:
// 0 - white player's turn
// 1 - black player's turn
// 2 - white player won
// 3 - black player won
// 4 - draw
type GameState int

func (g *Game) Init() {
	g.Board.Init()
	g.State = 0
}

func (g *Game) Print() {
	g.Board.Print()
}

func (g *Game) Play() {
	// If the game is over, do nothing.
	if g.State == 2 || g.State == 3 || g.State == 4 {
		return
	}
	// Perform game logic
	// TODO
}

func main() {
	var board1 = InitFEN("4k2r/8/8/8/8/8/8/R3K2R w KQk - 2 4")
	board1.PrintWithBorder()

	whiteKing := board1.GetPieceAtLocation(Location{X: 0, Y: 4})
	fmt.Println(whiteKing)
	kingMoves := whiteKing.GetMoves(&board1)
	for _, move := range kingMoves {
		fmt.Println(move.ToString())
	}

}
