package main

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
	board := InitFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	board.PrintWithBorder()
}
