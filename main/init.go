package main

import (
	"fmt"
)

func main() {
	// game := Game{}
	// game.Init(&HumanPlayer{1}, &HumanPlayer{2})
	// for {
	// 	fmt.Println("")
	// 	game.Print()
	// 	game.Play()
	// }

	// Play a game with random player
	game := Game{}
	//game.Init(&RandomComputerPlayer{1}, &MinimaxComputerPlayer{2, 3})
	//game.Init(&MinimaxComputerPlayer{1, 3}, &MinimaxComputerPlayer{2, 3})
	game.Init(&MinimaxComputerPlayer{1, 3}, &RandomComputerPlayer{2})

	game.Print()
	for {
		fmt.Println("")
		fmt.Println(game.Board.FullmoveNumber)
		game.Print()
		isEnd := game.Play()
		if isEnd {
			break
		}
		//time.Sleep(time.Millisecond * 500)
	}
	game.Print()

}
