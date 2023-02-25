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

}

func main() {
	var board Board
	board.Init()
	fmt.Println("Board initialized.")
	board.Print()

	fmt.Println()
	fmt.Println("Test what the pawn can do:")
	fmt.Println(board.Locations[1][0])
	moves := board.Locations[1][0].GetMoves(&board)
	fmt.Println(moves)
}
