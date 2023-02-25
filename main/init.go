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
	fmt.Print("White Move:")
	m := Move{
		Type:             'M',
		Piece:            'P',
		IsDisambiguation: false,
		From:             Location{4, 1},
		To:               Location{4, 3},
	}
	fmt.Println(m.ToString())
	board.MakeMove(m)
	board.Print()

	fmt.Println()
	fmt.Print("Black Move:")
	m = Move{
		Type:             'M',
		Piece:            'P',
		IsDisambiguation: false,
		From:             Location{4, 6},
		To:               Location{4, 4},
	}
	fmt.Println(m.ToString())
	board.MakeMove(m)
	board.Print()
}
