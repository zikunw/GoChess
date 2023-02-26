package main

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

func (g *Game) Print() {
	g.Board.PrintWithBorder()
}

// Initiate the game logic
func (g *Game) Play() {
	// If the game is over, do nothing.
	if g.State == 2 || g.State == 3 || g.State == 4 {
		return
	}

	// If it's white player's turn, get the move from the white player
	if g.State == 0 {
		move := g.WhitePlayer.GetMove(&g.Board)
		g.Board.MakeMove(move)
		g.State = 1
	} else {
		move := g.BlackPlayer.GetMove(&g.Board)
		g.Board.MakeMove(move)
		g.State = 0
	}

	// Check if the game is over
	// TODO: Currently the game wont stop
}
