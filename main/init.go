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
	var board1 = InitFEN("4k3/6B1/5PR1/8/2R5/8/6N1/4K3 w - - 0 1")
	board1.PrintWithBorder()

	var piece1 = board1.GetPieceAtLocation(Location{X: 3, Y: 2})
	fmt.Println(piece1)
	var moves1 = piece1.GetMoves(&board1)
	for _, move := range moves1 {
		fmt.Println(move.ToString())
	}

	fmt.Println()
	var piece2 = board1.GetPieceAtLocation(Location{X: 5, Y: 6})
	fmt.Println(piece2)
	var moves2 = piece2.GetMoves(&board1)
	for _, move := range moves2 {
		fmt.Println(move.ToString())
	}

}
