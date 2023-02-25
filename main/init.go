package main

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
	var board1 = InitFEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2 ")
	board1.Print()

}
