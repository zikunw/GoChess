package main

import (
	"fmt"
	"time"
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
	game.Init(&RandomComputerPlayer{1}, &RandomComputerPlayer{2})
	for {
		fmt.Println("")
		fmt.Println(game.Board.FullmoveNumber)
		game.Print()
		isEnd := game.Play()
		if isEnd {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

}
