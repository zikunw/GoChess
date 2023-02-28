package main

import (
	"fmt"
)

func main() {

	game := Game{}
	game.Init(&MinimaxComputerPlayer{1, 4}, &MinimaxComputerPlayer{2, 3})

	game.Print()
	for {
		fmt.Println("")
		fmt.Println(game.Board.FullmoveNumber)
		game.Print()
		isEnd := game.Play()
		if isEnd {
			break
		}
	}
	game.Print()

}
