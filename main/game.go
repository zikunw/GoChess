package main

import "fmt"

// The game logic

type Game struct {
	Board Board
	State GameState

	WhitePlayer PlayerController
	BlackPlayer PlayerController
}

// GameState:
// 0 - white player's turn
// 1 - black player's turn
// 2 - white player won
// 3 - black player won
// 4 - draw
type GameState int

func (g *Game) Init(whitePlayer PlayerController, blackPlayer PlayerController) {
	g.Board.Init()
	g.State = 0
	g.WhitePlayer = whitePlayer
	g.BlackPlayer = blackPlayer
}

// Init the game with a FEN string
// For debugging purposes
func (g *Game) InitWithFEN(fen string, whitePlayer PlayerController, blackPlayer PlayerController) {
	g.Board = InitFEN(fen)
	g.State = 0
	g.WhitePlayer = whitePlayer
	g.BlackPlayer = blackPlayer
}

func (g *Game) Print() {
	g.Board.PrintWithBorder()
}

// Initiate the game logic
// Returns true if the game is over
func (g *Game) Play() bool {
	// If the game is over, do nothing.
	if g.State == 2 || g.State == 3 || g.State == 4 {
		return true
	}

	// If it's white player's turn, get the move from the white player
	var move = Move{}
	if g.State == 0 {
		fmt.Println("White player's turn")
		move = g.WhitePlayer.GetMove(&g.Board)
		g.Board.MakeMove(move)
		g.State = 1
	} else {
		fmt.Println("Black player's turn")
		move = g.BlackPlayer.GetMove(&g.Board)
		g.Board.MakeMove(move)
		g.State = 0
	}

	// Check if the game is over
	if g.Board.CheckPlayerInCheckmate(1) {
		g.State = 3
		fmt.Println("Black player won")
		return true
	}
	if g.Board.CheckPlayerInCheckmate(2) {
		g.State = 2
		fmt.Println("White player won")
		return true
	}

	// Check if the game is a draw
	if g.Board.CheckPlayerInStalemate(1) || g.Board.CheckPlayerInStalemate(2) {
		g.State = 4
		fmt.Println("Draw by stalemate")
		return true
	}
	// Check half move rule
	if g.Board.HalfmoveClock >= 100 {
		g.State = 4
		fmt.Println("Draw by half move rule")
		return true
	}

	return false
}
