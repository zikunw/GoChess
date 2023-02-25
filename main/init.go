package main

func main() {
	game := Game{}
	game.Init(&HumanPlayer{1}, &HumanPlayer{2})
	for {
		game.Print()
		game.Play()
	}
}
