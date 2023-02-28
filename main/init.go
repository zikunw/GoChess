package main

import (
	"fmt"
)

func main() {

	//game := Game{}
	//game.Init(&HumanPlayer{1}, &MinimaxComputerPlayer{2, 3})
	//
	//game.Print()
	//for {
	//	fmt.Println("")
	//	fmt.Println(game.Board.FullmoveNumber)
	//	game.Print()
	//	isEnd := game.Play()
	//	if isEnd {
	//		break
	//	}
	//}
	//game.Print()

	b := Board{}
	b.Init()
	//b.Print()
	fmt.Println(b.Serialize())

	b2 := Board{}
	b2.Deserialize(b.Serialize())
	b2.Print()
	fmt.Print(b2.Serialize())
}
